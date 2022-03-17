package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUploadMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUploadMappingResourceConfig("https://example.com/images/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"cloudinary_upload_mapping.test",
						"template",
						"https://example.com/images/",
					),
					resource.TestCheckResourceAttr("cloudinary_upload_mapping.test", "folder", "example"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cloudinary_upload_mapping.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"template"},
			},
			// Update and Read testing
			{
				Config: testAccUploadMappingResourceConfig("https://example.org/images/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"cloudinary_upload_mapping.test",
						"template",
						"https://example.org/images/",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUploadMappingResourceConfig(template string) string {
	return fmt.Sprintf(`
resource "cloudinary_upload_mapping" "test" {
  folder   = "example"
	template = %[1]q
}
`, template)
}
