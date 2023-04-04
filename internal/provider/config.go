package provider

import (
	"errors"
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey"
	stdlog "log"
	"strings"

	"github.com/control-monkey-customer-z/terraform-provider-cm/version"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/credentials"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/featureflag"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/log"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/session"
	"github.com/control-monkey/controlmonkey-sdk-go/service/stack"
	"github.com/control-monkey/controlmonkey-sdk-go/service/variable"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
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
	stack    stack.Service
	variable variable.Service
}

// Client configures and returns a fully initialized ControlMonkey client.
func (c *Config) Client() (*Client, diag.Diagnostics) {
	stdlog.Println("[INFO] Configuring a new ControlMonkey client")

	// Create a new session.
	sess, err := c.getSession()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// Create a new client.
	client := &Client{
		stack:    stack.New(sess),
		variable: variable.New(sess),
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
		{Product: "Terraform Plugin SDK", Version: meta.SDKVersionString()},
		{Product: "Terraform Provider ControlMonkey", Version: "v2-" + version.String()},
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
