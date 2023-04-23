package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cmNamespace = "cm_namespace"

	n1ResourceName      = "namespace1"
	n1Name              = "namespace1"
	n1Description       = "first namespace test"
	n1PolicyMaxTtlType  = "days"
	n1PolicyMaxTtlValue = "2"
	n1PolicyDefTtlType  = "hours"
	n1PolicyDefTtlValue = "3"

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
  policy = {
	ttl_config = {
	  max_ttl = {
	    type = "%s"
 	    value = %s
	  }
	  default_ttl = {
	    type = "%s"
 	    value = %s
	  }
	}
  }
}
`, cmNamespace, n1ResourceName, n1Name, n1Description, n1PolicyMaxTtlType, n1PolicyMaxTtlValue, n1PolicyDefTtlType, n1PolicyDefTtlValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "name", n1Name),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "description", n1Description),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.max_ttl.type", n1PolicyMaxTtlType),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.max_ttl.value", n1PolicyMaxTtlValue),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.default_ttl.type", n1PolicyDefTtlType),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.default_ttl.value", n1PolicyDefTtlValue),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "id"),
					// No Attributes
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  policy = {
	ttl_config = {
	  max_ttl = {
	    type = "%s"
 	    value = %s
	  }
	  default_ttl = {
	    type = "%s"
 	    value = %s
	  }
	}
  }
}
`, cmNamespace, n1ResourceName, n1NameAfterUpdate, n1PolicyMaxTtlType, n1PolicyMaxTtlValue, n1PolicyDefTtlType, n1PolicyDefTtlValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "name", n1NameAfterUpdate),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.max_ttl.type", n1PolicyMaxTtlType),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.max_ttl.value", n1PolicyMaxTtlValue),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.default_ttl.type", n1PolicyDefTtlType),
					resource.TestCheckResourceAttr(namespaceResourceName(n1ResourceName), "policy.ttl_config.default_ttl.value", n1PolicyDefTtlValue),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(namespaceResourceName(n1ResourceName), "id"),
					// No Attributes
					resource.TestCheckNoResourceAttr(namespaceResourceName(n1ResourceName), "description"),
				),
			},
			{
				ResourceName:      fmt.Sprintf("%s.%s", cmNamespace, n1ResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func namespaceResourceName(s string) string {
	return fmt.Sprintf("%s.%s", cmNamespace, s)
}
