package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	client *cloudinary.Cloudinary

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	APIKey        types.String `tfsdk:"api_key"`
	APISecret     types.String `tfsdk:"api_secret"`
	CloudName     types.String `tfsdk:"cloud_name"`
	CloudinaryURL types.String `tfsdk:"cloudinary_url"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var cloud string
	if data.CloudName.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as cloud_name",
		)
		return
	}

	if data.CloudName.Null {
		cloud = os.Getenv("CLOUDINARY_CLOUD_NAME")
	} else {
		cloud = data.CloudName.Value
	}

	var key string
	if data.APIKey.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)
		return
	}

	if data.APIKey.Null {
		key = os.Getenv("CLOUDINARY_API_KEY")
	} else {
		key = data.APIKey.Value
	}

	var secret string
	if data.APISecret.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_secret",
		)
		return
	}

	if data.APISecret.Null {
		secret = os.Getenv("CLOUDINARY_API_SECRET")
	} else {
		secret = data.APISecret.Value
	}

	var cldURL string
	if data.CloudinaryURL.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as cloudinary_url",
		)
		return
	}

	if data.CloudinaryURL.Null {
		cldURL = os.Getenv("CLOUDINARY_URL")

		if cldURL == "" && cloud != "" && key != "" && secret != "" {
			cldURL = fmt.Sprintf("cloudinary://%s:%s@%s", key, secret, cloud)
		}
	} else {
		cldURL = data.CloudinaryURL.Value
	}

	cld, err := cloudinary.NewFromURL(cldURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create cloudinary client:\n\n"+err.Error(),
		)

		return
	}

	p.client = cld
	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"cloudinary_admin_upload_mapping": adminUploadMappingResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"cloudinary_admin_upload_mapping": adminUploadMappingDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				Computed:  true,
				Optional:  true,
				Sensitive: true,
				Type:      types.StringType,
			},
			"api_secret": {
				Computed:  true,
				Optional:  true,
				Sensitive: true,
				Type:      types.StringType,
			},
			"cloud_name": {
				Computed: true,
				Optional: true,
				Type:     types.StringType,
			},
			"cloudinary_url": {
				Computed:  true,
				Optional:  true,
				Sensitive: true,
				Type:      types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
