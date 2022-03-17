package provider

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type uploadMappingResourceType struct{}

func (t uploadMappingResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Upload Mapping resource.",

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
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t uploadMappingResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return uploadMappingResource{
		provider: provider,
	}, diags
}

type uploadMappingResourceData struct {
	Folder   types.String `tfsdk:"folder"`
	ID       types.String `tfsdk:"id"`
	Template types.String `tfsdk:"template"`
}

type uploadMappingResource struct {
	provider provider
}

func (r uploadMappingResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data uploadMappingResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Folder

	params := admin.CreateUploadMappingParams{
		Folder:   data.Folder.Value,
		Template: data.Template.Value,
	}

	res, err := r.provider.client.Admin.CreateUploadMapping(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create upload mapping, got error: %s", err),
		)
		return
	}

	if res.Error.Message != "" {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create upload mapping, got error: %s", res.Error.Message),
		)
		return
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r uploadMappingResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data uploadMappingResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Folder

	params := admin.GetUploadMappingParams{
		Folder: data.Folder.Value,
	}

	res, err := r.provider.client.Admin.GetUploadMapping(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read upload mapping, got error: %s", err),
		)
		return
	}

	if res.Error.Message != "" {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read upload mapping, got error: %s", res.Error.Message),
		)
		return
	}

	data.Template = types.String{Value: res.Template}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r uploadMappingResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data uploadMappingResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Folder

	params := admin.UpdateUploadMappingParams{
		Folder:   data.Folder.Value,
		Template: data.Template.Value,
	}

	res, err := r.provider.client.Admin.UpdateUploadMapping(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update upload mapping, got error: %s", err),
		)
		return
	}

	if res.Error.Message != "" {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update upload mapping, got error: %s", res.Error.Message),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r uploadMappingResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data uploadMappingResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Folder

	params := admin.DeleteUploadMappingParams{
		Folder: data.Folder.Value,
	}

	res, err := r.provider.client.Admin.DeleteUploadMapping(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete upload mapping, got error: %s", err),
		)
		return
	}

	if res.Error.Message != "" {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete upload mapping, got error: %s", res.Error.Message),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r uploadMappingResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("folder"), req, resp)
}
