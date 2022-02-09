package provider

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type adminUploadMappingDataSourceType struct{}

func (t adminUploadMappingDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Upload mapping data source.",

		Attributes: map[string]tfsdk.Attribute{
			"folder": {
				MarkdownDescription: "The name of the folder.",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"template": {
				MarkdownDescription: "The URL to be mapped to the folder.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t adminUploadMappingDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return adminUploadMappingDataSource{
		provider: provider,
	}, diags
}

type adminUploadMappingDataSourceData struct {
	Folder   types.String `tfsdk:"folder"`
	ID       types.String `tfsdk:"id"`
	Template types.String `tfsdk:"template"`
}

type adminUploadMappingDataSource struct {
	provider provider
}

func (d adminUploadMappingDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data adminUploadMappingDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Folder

	params := admin.GetUploadMappingParams{
		Folder: data.Folder.Value,
	}

	result, err := d.provider.client.Admin.GetUploadMapping(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read upload mapping, got error: %s", err),
		)

		return
	}

	data.Template = types.String{Value: result.Template}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
