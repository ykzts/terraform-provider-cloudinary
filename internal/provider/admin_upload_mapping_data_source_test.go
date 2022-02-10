package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAdminUploadMappingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAdminUploadMappingDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudinary_admin_upload_mapping.test", "folder", "example-data"),
				),
			},
		},
	})
}

const testAccAdminUploadMappingDataSourceConfig = `
data "cloudinary_admin_upload_mapping" "test" {
  folder = "example-data"
}
`
