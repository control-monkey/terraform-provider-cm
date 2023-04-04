package provider

import (
	"fmt"
	"github.com/control-monkey/controlmonkey-sdk-go/controlmonkey/credentials"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"os"
	"testing"
)

const (
	// providerConfig is a shared configuration to combine with the actual test configuration.
	providerConfig = `
provider "cm" {}
`
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(credentials.EnvCredentialsVarToken); v == "" {
		t.Fatal(fmt.Printf("%s must be set for acceptance tests", credentials.EnvCredentialsVarToken))
	}
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"cm": providerserver.NewProtocol6WithError(New()),
	}
)
