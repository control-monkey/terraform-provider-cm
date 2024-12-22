package provider

import (
	"errors"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/services/blueprint"
	"github.com/control-monkey/controlmonkey-sdk-go/services/control_policy"
	"github.com/control-monkey/controlmonkey-sdk-go/services/control_policy_group"
	"github.com/control-monkey/controlmonkey-sdk-go/services/custom_abac_configuration"
	"github.com/control-monkey/controlmonkey-sdk-go/services/custom_role"
	"github.com/control-monkey/controlmonkey-sdk-go/services/disaster_recovery"
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace_permissions"
	"github.com/control-monkey/controlmonkey-sdk-go/services/notification"
	"github.com/control-monkey/controlmonkey-sdk-go/services/organization"
	"github.com/control-monkey/controlmonkey-sdk-go/services/team"
	stdlog "log"
	"strings"

	"github.com/control-monkey/controlmonkey-sdk-go/services/template"

	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/credentials"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/featureflag"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/log"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/session"
	"github.com/control-monkey/controlmonkey-sdk-go/services/namespace"
	"github.com/control-monkey/controlmonkey-sdk-go/services/stack"
	"github.com/control-monkey/controlmonkey-sdk-go/services/variable"
	"github.com/control-monkey/terraform-provider-cm/version"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var ErrNoValidCredentials = errors.New("\n\nNo valid credentials found " +
	"for ControlMonkey Provider.\nPlease see https://www.terraform.io/docs/" +
	"providers/internal/index.html\nfor more information on providing " +
	"credentials for ControlMonkey Provider.")

type Config struct {
	Token        string
	FeatureFlags string

	terraformVersion string
}

type Client struct {
	blueprint               blueprint.Service
	controlPolicy           control_policy.Service
	controlPolicyGroup      control_policy_group.Service
	customAbacConfiguration custom_abac_configuration.Service
	customRole              custom_role.Service
	disasterRecovery        disaster_recovery.Service
	namespace               namespace.Service
	namespacePermissions    namespace_permissions.Service
	notification            notification.Service
	organization            organization.Service
	stack                   stack.Service
	team                    team.Service
	template                template.Service
	variable                variable.Service
}

// Client configures and returns a fully initialized ControlMonkey client.
func (c *Config) Client() (*Client, diag.Diagnostics) {
	stdlog.Println("[INFO] Configuring a new ControlMonkey client")

	// Create a new session.
	sess, err := c.getSession()
	if err != nil {
		diags := new(diag.Diagnostics)
		diags.AddError("Failed to configure ControlMonkey client", err.Error())
		return nil, *diags
	}

	// Create a new client.
	client := &Client{
		blueprint:               blueprint.New(sess),
		controlPolicy:           control_policy.New(sess),
		controlPolicyGroup:      control_policy_group.New(sess),
		customAbacConfiguration: custom_abac_configuration.New(sess),
		customRole:              custom_role.New(sess),
		disasterRecovery:        disaster_recovery.New(sess),
		namespace:               namespace.New(sess),
		namespacePermissions:    namespace_permissions.New(sess),
		notification:            notification.New(sess),
		organization:            organization.New(sess),
		stack:                   stack.New(sess),
		team:                    team.New(sess),
		template:                template.New(sess),
		variable:                variable.New(sess),
	}

	stdlog.Println("[INFO] ControlMonkey client configured")
	return client, nil
}

func (c *Config) getSession() (*session.Session, error) {
	config := controlmonkey.DefaultConfig()

	// HTTP options.
	{
		config.WithHTTPClient(cleanhttp.DefaultPooledClient())
		config.WithUserAgent(c.getUserAgent())
	}

	// Credentials.
	{
		v, err := c.getCredentials()
		if err != nil {
			return nil, err
		}
		config.WithCredentials(v)
	}

	// Logging.
	{
		config.WithLogger(log.LoggerFunc(func(format string, args ...interface{}) {
			stdlog.Printf(fmt.Sprintf("[DEBUG] [internal-sdk-go] %s", format), args...)
		}))
	}

	return session.New(config), nil
}

func (c *Config) getUserAgent() string {
	agents := []struct {
		Product string
		Version string
		Comment []string
	}{
		{Product: "HashiCorp", Version: "1.0"},
		{Product: "Terraform", Version: c.terraformVersion, Comment: []string{"+https://www.terraform.io"}},
		{Product: "Terraform Provider ControlMonkey", Version: "v" + version.String()},
	}

	var ua string
	for _, agent := range agents {
		pv := fmt.Sprintf("%s/%s", agent.Product, agent.Version)
		if len(agent.Comment) > 0 {
			pv += fmt.Sprintf(" (%s)", strings.Join(agent.Comment, "; "))
		}
		if len(ua) > 0 {
			ua += " "
		}
		ua += pv
	}

	return ua
}

func (c *Config) getCredentials() (*credentials.Credentials, error) {
	var providers []credentials.Provider
	var static *credentials.StaticProvider

	featureflag.Set(c.FeatureFlags)

	if c.Token != "" {
		static = &credentials.StaticProvider{
			Value: credentials.Value{
				Token: c.Token,
			},
		}
	}
	if static != nil {
		providers = append(providers, static)
	}

	providers = append(providers,
		new(credentials.EnvProvider),
		new(credentials.FileProvider))

	creds := credentials.NewChainCredentials(providers...)

	if _, err := creds.Get(); err != nil {
		stdlog.Printf("[ERROR] Failed to instantiate ControlMonkey client: %v", err)
		return nil, ErrNoValidCredentials
	}

	return creds, nil
}
