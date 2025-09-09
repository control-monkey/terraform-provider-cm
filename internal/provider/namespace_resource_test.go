package provider

import (
	"fmt"
	"testing"

	"github.com/control-monkey/terraform-provider-cm/internal/provider/commons/test_helpers"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNamespace = "cm_namespace"

	n1ResourceName = "namespace1"
	n1Name         = "namespace1"
	n1Description  = "first namespace test"

	n1NameAfterUpdate = "namespace2"
)

func TestAccNamespaceResourceNamespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  description = "%s"
}
`, cmNamespace, n1ResourceName, n1Name, n1Description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "name", n1Name),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "description", n1Description),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "id"),
					// No Attributes
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
  name = "Namespace Test 1"
}

resource "cm_team" "team2" {
  name = "Namespace Test 2"
}

resource "%s" "%s" {
  name = "%s"
  deployment_approval_policy = {
  	override_behavior = "allow"
    rules = [
      {
        type = "requireTeamsApproval"
        parameters = jsonencode({
          teams = [cm_team.team1.id, cm_team.team2.id]
        })
      },
      {
        type = "requireTwoApprovals"
      },
    ]
  }
}
`, cmNamespace, n1ResourceName, n1NameAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "name", n1NameAfterUpdate),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "id"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.0.type", "requireTeamsApproval"),
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.0.parameters"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.1.type", "requireTwoApprovals"),
					// No Attributes
					resource.TestCheckNoResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.1.parameters"),
					resource.TestCheckNoResourceAttr(namespaceResourceName(n1ResourceName), "description"),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			// Test capabilities
			{
				Config: providerConfig + fmt.Sprintf(`
resource "cm_team" "team1" {
  name = "Namespace Test 1"
}

resource "cm_team" "team2" {
  name = "Namespace Test 2"
}

resource "%s" "%s" {
  name = "%s"
  deployment_approval_policy = {
  	override_behavior = "allow"
    rules = [
      {
        type = "requireTeamsApproval"
        parameters = jsonencode({
          teams = [cm_team.team1.id, cm_team.team2.id]
        })
      },
      {
        type = "requireTwoApprovals"
      },
    ]
  }

  capabilities = {
    deploy_on_push = {
      status = "enabled"
      is_overridable = true
    }
    plan_on_pr = {
      status = "disabled" 
      is_overridable = false
    }
    drift_detection = {
      status = "enabled"
      is_overridable = true
    }
  }
}
`, cmNamespace, n1ResourceName, n1NameAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "name", n1NameAfterUpdate),
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "id"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.0.type", "requireTeamsApproval"),
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.0.parameters"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.1.type", "requireTwoApprovals"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.deploy_on_push.status", "enabled"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.deploy_on_push.is_overridable", "true"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.plan_on_pr.status", "disabled"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.plan_on_pr.is_overridable", "false"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.drift_detection.status", "enabled"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.drift_detection.is_overridable", "true"),
					resource.TestCheckNoResourceAttr(namespaceResourceName(n1ResourceName), "description"),
				),
			},
			test_helpers.GetValidateNoDriftStep(),
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmNamespace, n1ResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "name", n1NameAfterUpdate),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "id"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.0.type", "requireTeamsApproval"),
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.0.parameters"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.1.type", "requireTwoApprovals"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.deploy_on_push.status", "enabled"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.deploy_on_push.is_overridable", "true"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.plan_on_pr.status", "disabled"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.plan_on_pr.is_overridable", "false"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.drift_detection.status", "enabled"),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "capabilities.drift_detection.is_overridable", "true"),
					// No Attributes
					resource.TestCheckNoResourceAttr(namespaceResourceName(n1ResourceName), "deployment_approval_policy.rules.1.parameters"),
					resource.TestCheckNoResourceAttr(namespaceResourceName(n1ResourceName), "description"),
				),
			},
		},
	})
}

func namespaceResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmNamespace, s)
}
