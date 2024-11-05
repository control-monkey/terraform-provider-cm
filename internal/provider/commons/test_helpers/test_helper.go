package test_helpers

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func GetValidateNoDriftStep() resource.TestStep {
	return resource.TestStep{
		RefreshState: true,
		RefreshPlanChecks: resource.RefreshPlanChecks{
			PostRefresh: []plancheck.PlanCheck{
				plancheck.ExpectEmptyPlan(),
			},
		},
	}
}
