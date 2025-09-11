package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/control-monkey/controlmonkey-sdk-go/services/commons"
	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_config"
	"github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmVariable = "cm_variable"

	namespaceVariable              = "var_namespace"
	namespaceVariableScope         = "namespace"
	namespaceVariableKey           = "namespaceVar"
	namespaceVariableType          = "tfVar"
	namespaceVariableValue         = "TfValue"
	namespaceVariableDisplayName   = "Display Name"
	namespaceVariableIsSensitive   = "false"
	namespaceVariableIsOverridable = "true"
	namespaceVariableIsRequired    = "false"

	namespaceVariableValueAfterUpdate        = "TfValue2"
	namespaceVariableNumericValue            = "5"
	namespaceVariableNumericValueAfterUpdate = "10"

	orgVariable              = "var_org"
	orgVariableScope         = "organization"
	orgVariableKey           = "orgKey"
	orgVariableValue         = "TfValue"
	orgVariableType          = "envVar"
	orgVariableIsSensitive   = "false"
	orgVariableIsOverridable = "false"

	orgVariableValueAfterUpdate         = "TfValue2"
	orgVariableIsOverridableAfterUpdate = "true"

	blueprintVariable              = "var_blueprint"
	blueprintVariableScope         = "blueprint"
	blueprintVariableKey           = "blueprintVar"
	blueprintVariableType          = "tfVar"
	blueprintVariableValue         = "BlueprintValue"
	blueprintVariableIsSensitive   = "false"
	blueprintVariableIsOverridable = "true"
	blueprintVariableManagedBy     = "stack"
	blueprintVariableDescription   = "Blueprint variable description"

	blueprintVariableValueAfterUpdate     = "BlueprintValue2"
	blueprintVariableManagedByAfterUpdate = "inCode"
)

func TestAccVariableResourceNamespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace" {
  name = "variable test"
}

resource "%s" "%s" {
	scope          = "%s"
	scope_id       = cm_namespace.namespace.id
	key            = "%s"
	type           = "%s"
	value          = "%s"
	display_name   = "%s"
	is_sensitive   = %s
	is_overridable = %s
	is_required    = %s
}
					`, tfCmVariable, namespaceVariable,
					namespaceVariableScope, namespaceVariableKey, namespaceVariableType, namespaceVariableValue, namespaceVariableDisplayName,
					namespaceVariableIsSensitive, namespaceVariableIsOverridable, namespaceVariableIsRequired),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "scope", namespaceVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(namespaceVariable), "scope_id", "cm_namespace.namespace", "id"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "key", namespaceVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "type", namespaceVariableType),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value", namespaceVariableValue),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "display_name", namespaceVariableDisplayName),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_sensitive", namespaceVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_overridable", namespaceVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_required", namespaceVariableIsRequired),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(namespaceVariable), "id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace" {
  name = "variable test"
}

resource "%s" "%s" {
	scope          = "%s"
	scope_id       = cm_namespace.namespace.id
	key            = "%s"
	type           = "%s"
	value          = "%s"
	is_sensitive   = %s
	is_overridable = %s
	is_required = %s
}
						`, tfCmVariable, namespaceVariable,
					namespaceVariableScope, namespaceVariableKey, namespaceVariableType,
					namespaceVariableValueAfterUpdate, namespaceVariableIsSensitive, namespaceVariableIsOverridable, namespaceVariableIsRequired),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "scope", namespaceVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(namespaceVariable), "scope_id", "cm_namespace.namespace", "id"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "key", namespaceVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "type", namespaceVariableType),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value", namespaceVariableValueAfterUpdate),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_sensitive", namespaceVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_overridable", namespaceVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_required", namespaceVariableIsRequired),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(namespaceVariable), "id"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "display_name"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace" {
  name = "variable test"
}

resource "%s" "%s" {
	scope          = "%s"
	scope_id       = cm_namespace.namespace.id
	key            = "%s"
	type           = "%s"
	value          = "%s"
	is_sensitive   = %s
	is_overridable = %s
	value_conditions = [
	  {
	    operator = "lt"
	    value    = 50
	  },
	  {
	    operator = "gt"
	    value    = 5
	  },
  	]
}
						`, tfCmVariable, namespaceVariable,
					namespaceVariableScope, namespaceVariableKey, namespaceVariableType,
					namespaceVariableNumericValue, namespaceVariableIsSensitive, namespaceVariableIsOverridable),
				ExpectError: regexp.MustCompile(commons.ErrorCodeValidationError),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace" {
  name = "variable test"
}

resource "%s" "%s" {
	scope          = "%s"
	scope_id       = cm_namespace.namespace.id
	key            = "%s"
	type           = "%s"
	value          = "%s"
	is_sensitive   = %s
	is_overridable = %s
	value_conditions = [
	  {
	    operator = "lt"
	    value    = 50
	  },
	  {
	    operator = "gt"
	    value    = 5
	  },
  	]
}
						`, tfCmVariable, namespaceVariable,
					namespaceVariableScope, namespaceVariableKey, namespaceVariableType,
					namespaceVariableNumericValueAfterUpdate, namespaceVariableIsSensitive, namespaceVariableIsOverridable),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "scope", namespaceVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(namespaceVariable), "scope_id", "cm_namespace.namespace", "id"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "key", namespaceVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "type", namespaceVariableType),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value", namespaceVariableNumericValueAfterUpdate),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_sensitive", namespaceVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_overridable", namespaceVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.#", "2"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(namespaceVariable), "id"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "display_name"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "is_required"),
				),
			},
			{
				ConfigVariables: config.Variables{
					"lte_value": config.StringVariable("50"),
					"gte_value": config.StringVariable("5"),
				},
				Config: providerConfig + fmt.Sprintf(`
resource "cm_namespace" "namespace" {
  name = "variable test"
}

variable "lte_value" {
	type = string
}

variable "gte_value" {
	type = string
}

resource "%s" "%s" {
	scope          = "%s"
	scope_id       = cm_namespace.namespace.id
	key            = "%s"
	type           = "%s"
	value          = "%s"
	is_sensitive   = %s
	is_overridable = %s
	value_conditions = [
	  {
	    operator = "lt"
	    value    = var.lte_value
	  },
	  {
	    operator = "gt"
	    value    = var.gte_value
	  },
  	]
}
						`, tfCmVariable, namespaceVariable,
					namespaceVariableScope, namespaceVariableKey, namespaceVariableType,
					namespaceVariableNumericValueAfterUpdate, namespaceVariableIsSensitive, namespaceVariableIsOverridable),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "scope", namespaceVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(namespaceVariable), "scope_id", "cm_namespace.namespace", "id"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "key", namespaceVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "type", namespaceVariableType),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value", namespaceVariableNumericValueAfterUpdate),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_sensitive", namespaceVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_overridable", namespaceVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.#", "2"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.0.value", "50"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.1.value", "5"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(namespaceVariable), "id"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "display_name"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "is_required"),
				),
			},
			{
				ConfigVariables: config.Variables{
					"lte_value": config.StringVariable("50"),
					"gte_value": config.StringVariable("5"),
				},
				ResourceName:      fmt.Sprintf("%s.%s", tfCmVariable, namespaceVariable),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "scope", namespaceVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(namespaceVariable), "scope_id", "cm_namespace.namespace", "id"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "key", namespaceVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "type", namespaceVariableType),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value", namespaceVariableNumericValueAfterUpdate),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_sensitive", namespaceVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "is_overridable", namespaceVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.#", "2"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.0.value", "50"),
					resource.TestCheckResourceAttr(variableResourceName(namespaceVariable), "value_conditions.1.value", "5"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(namespaceVariable), "id"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "display_name"),
					resource.TestCheckNoResourceAttr(variableResourceName(namespaceVariable), "is_required"),
				),
			},
		},
	})
}

func TestAccVariableResourceOrganization(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
	scope          = "%s"
	key            = "%s"
	type           = "%s"
	value          = "%s"
	is_sensitive   = %s
	is_overridable = %s
}
					`, tfCmVariable, orgVariable, orgVariableScope, orgVariableKey, orgVariableType, orgVariableValue,
					orgVariableIsSensitive, orgVariableIsOverridable),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "scope", orgVariableScope),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "key", orgVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "type", orgVariableType),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "value", orgVariableValue),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "is_sensitive", orgVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "is_overridable", orgVariableIsOverridable),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(orgVariable), "id"),
					resource.TestCheckNoResourceAttr(variableResourceName(orgVariable), "is_required"),
					resource.TestCheckNoResourceAttr(variableResourceName(orgVariable), "description"),
				),
			},
			//Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
						resource "%s" "%s" {
							scope          = "%s"
							key            = "%s"
							type           = "%s"
							value          = "%s"
							is_sensitive   = %s
							is_overridable = %s
						}
						`, tfCmVariable, orgVariable,
					orgVariableScope, orgVariableKey, orgVariableType, orgVariableValueAfterUpdate,
					orgVariableIsSensitive, orgVariableIsOverridableAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "scope", orgVariableScope),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "key", orgVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "type", orgVariableType),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "value", orgVariableValueAfterUpdate),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "is_sensitive", orgVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(orgVariable), "is_overridable", orgVariableIsOverridableAfterUpdate),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(orgVariable), "id"),
				),
			},
		},
	})
}

func TestAccVariableResourceBlueprint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccVariableBlueprintResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
	scope                        = "%s"
	scope_id                     = cm_blueprint.test_blueprint.id
	key                          = "%s"
	type                         = "%s"
	value                        = "%s"
	description                  = "%s"
	is_sensitive                 = %s
	is_overridable               = %s
	blueprint_variable_managed_by = "%s"
}
`, tfCmVariable, blueprintVariable, blueprintVariableScope, blueprintVariableKey,
					blueprintVariableType, blueprintVariableValue, blueprintVariableDescription,
					blueprintVariableIsSensitive, blueprintVariableIsOverridable, blueprintVariableManagedBy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "scope", blueprintVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(blueprintVariable), "scope_id", "cm_blueprint.test_blueprint", "id"),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "key", blueprintVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "type", blueprintVariableType),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "value", blueprintVariableValue),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "description", blueprintVariableDescription),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "is_sensitive", blueprintVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "is_overridable", blueprintVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "blueprint_variable_managed_by", blueprintVariableManagedBy),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(blueprintVariable), "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccVariableBlueprintResourceSetup() + fmt.Sprintf(`
resource "%s" "%s" {
	scope                        = "%s"
	scope_id                     = cm_blueprint.test_blueprint.id
	key                          = "%s"
	type                         = "%s"
	value                        = "%s"
	description                  = "%s"
	is_sensitive                 = %s
	is_overridable               = %s
	blueprint_variable_managed_by = "%s"
}
`, tfCmVariable, blueprintVariable, blueprintVariableScope, blueprintVariableKey,
					blueprintVariableType, blueprintVariableValueAfterUpdate, blueprintVariableDescription,
					blueprintVariableIsSensitive, blueprintVariableIsOverridable, blueprintVariableManagedByAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "scope", blueprintVariableScope),
					resource.TestCheckResourceAttrPair(variableResourceName(blueprintVariable), "scope_id", "cm_blueprint.test_blueprint", "id"),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "key", blueprintVariableKey),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "type", blueprintVariableType),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "value", blueprintVariableValueAfterUpdate),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "description", blueprintVariableDescription),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "is_sensitive", blueprintVariableIsSensitive),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "is_overridable", blueprintVariableIsOverridable),
					resource.TestCheckResourceAttr(variableResourceName(blueprintVariable), "blueprint_variable_managed_by", blueprintVariableManagedByAfterUpdate),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(variableResourceName(blueprintVariable), "id"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", tfCmVariable, blueprintVariable),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVariableBlueprintResourceSetup() string {
	// Test environment variables used by this function
	providerId := test_config.GetProviderId()
	repoName := test_config.GetRepoName()

	return providerConfig + fmt.Sprintf(`
resource "cm_blueprint" "test_blueprint" {
    name = "Variable Test Blueprint"
    description = "Blueprint for testing variables"

    blueprint_vcs_info = {
        provider_id = "%s"
        repo_name = "%s"
        path = "cm/blueprint"
    }

    stack_configuration = {
        name_pattern = "{stack_name}"
        iac_type = "terraform"

        vcs_info_with_patterns = {
            provider_id = "%s"
            repo_name = "%s"
            path_pattern = "{stack_path}"
        }
    }

    substitute_parameters = [
		{
			key = "stack_name"
			description = "any name you want"
		},
		{
			key = "stack_path"
			description = "path"
		}
	]
}
`, providerId, repoName, providerId, repoName)
}

func variableResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmVariable, s)
}
