package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	tfCmControlPolicyMapping = "cm_control_policy_mapping"

	mappingResourceName = "mapping"
	controlPolicyId     = "pol-868pl1qopp"
	targetId            = "ns-x82yjdyahc"
	targetType          = "namespace"
	enforcementLevel    = "hardMandatory"

	enforcementLevelAfterUpdate = "softMandatory"
)

func TestAccControlPolicyMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
						resource "%s" "%s" {
						control_policy_id = "%s"
						target_id         = "%s"
						target_type       = "%s"
						enforcement_level = "%s"
					}
					`, tfCmControlPolicyMapping, mappingResourceName,
					controlPolicyId, targetId, targetType, enforcementLevel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "target_id", targetId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "target_type", targetType),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "enforcement_level", enforcementLevel),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "id", getId()),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
						resource "%s" "%s" {
						control_policy_id = "%s"
						target_id         = "%s"
						target_type       = "%s"
						enforcement_level = "%s"
					}
					`, tfCmControlPolicyMapping, mappingResourceName,
					controlPolicyId, targetId, targetType, enforcementLevelAfterUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "target_id", targetId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "target_type", targetType),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "enforcement_level", enforcementLevelAfterUpdate),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "id", getId()),
				),
			},
			{
				ResourceName:  controlPolicyMappingResourceName(mappingResourceName),
				ImportStateId: fmt.Sprintf("%s/%s/%s", controlPolicyId, targetId, targetType),
				ImportState:   true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "control_policy_id", controlPolicyId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "target_id", targetId),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "target_type", targetType),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "enforcement_level", enforcementLevelAfterUpdate),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "enforcement_level", enforcementLevelAfterUpdate),

					resource.TestCheckResourceAttrSet(controlPolicyMappingResourceName(mappingResourceName), "id"),
					resource.TestCheckResourceAttr(controlPolicyMappingResourceName(mappingResourceName), "id", getId()),
				),
			},
		},
	})
}

func getId() string {
	return fmt.Sprintf("%s/%s/%s", controlPolicyId, targetId, targetType)
}

func controlPolicyMappingResourceName(s string) string {
	return fmt.Sprintf("%s.%s", tfCmControlPolicyMapping, s)
}
