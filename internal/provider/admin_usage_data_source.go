package provider

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type adminUsageDataSourceType struct{}

func (t adminUsageDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Usage data source.",

		Attributes: map[string]tfsdk.Attribute{
			"bandwidth": {
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"limit": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"usage": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"used_percent": {
							Computed: true,
							Type:     types.Float64Type,
						},
					},
				),
				Computed: true,
			},
			"derived_resources": {
				Computed: true,
				Type:     types.Int64Type,
			},
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"last_updated": {
				Computed: true,
				Type:     types.StringType,
			},
			"media_limits": {
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"asset_max_total_px": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"image_max_px": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"image_max_size_bytes": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"raw_max_size_bytes": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"video_max_size_bytes": {
							Computed: true,
							Type:     types.Int64Type,
						},
					},
				),
				Computed: true,
			},
			"objects": {
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"limit": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"usage": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"used_percent": {
							Computed: true,
							Type:     types.Float64Type,
						},
					},
				),
				Computed: true,
			},
			"plan": {
				Computed: true,
				Type:     types.StringType,
			},
			"requests": {
				Computed: true,
				Type:     types.Int64Type,
			},
			"resources": {
				Computed: true,
				Type:     types.Int64Type,
			},
			"storage": {
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"limit": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"usage": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"used_percent": {
							Computed: true,
							Type:     types.Float64Type,
						},
					},
				),
				Computed: true,
			},
			"transformations": {
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"limit": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"usage": {
							Computed: true,
							Type:     types.Int64Type,
						},
						"used_percent": {
							Computed: true,
							Type:     types.Float64Type,
						},
					},
				),
				Computed: true,
			},
		},
	}, nil
}

func (t adminUsageDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return adminUsageDataSource{
		provider: provider,
	}, diags
}

type adminUsageUsageData struct {
	Limit       types.Int64   `tfsdk:"limit"`
	Usage       types.Int64   `tfsdk:"usage"`
	UsedPercent types.Float64 `tfsdk:"used_percent"`
}

type adminUsageMediaLimitsData struct {
	AssetMaxTotalPx   types.Int64 `tfsdk:"asset_max_total_px"`
	ImageMaxPx        types.Int64 `tfsdk:"image_max_px"`
	ImageMaxSizeBytes types.Int64 `tfsdk:"image_max_size_bytes"`
	RawMaxSizeBytes   types.Int64 `tfsdk:"raw_max_size_bytes"`
	VideoMaxSizeBytes types.Int64 `tfsdk:"video_max_size_bytes"`
}

type adminUsageDataSourceData struct {
	Bandwidth        adminUsageUsageData       `tfsdk:"bandwidth"`
	DerivedResources types.Int64               `tfsdk:"derived_resources"`
	ID               types.String              `tfsdk:"id"`
	LastUpdated      types.String              `tfsdk:"last_updated"`
	MediaLimits      adminUsageMediaLimitsData `tfsdk:"media_limits"`
	Objects          adminUsageUsageData       `tfsdk:"objects"`
	Plan             types.String              `tfsdk:"plan"`
	Requests         types.Int64               `tfsdk:"requests"`
	Resources        types.Int64               `tfsdk:"resources"`
	Storage          adminUsageUsageData       `tfsdk:"storage"`
	Transformations  adminUsageUsageData       `tfsdk:"transformations"`
}

type adminUsageDataSource struct {
	provider provider
}

func (d adminUsageDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data adminUsageDataSourceData

	params := admin.UsageParams{}

	res, err := d.provider.client.Admin.Usage(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read usage, got error: %s", err),
		)
		return
	}

	if res.Error.Message != "" {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read usage, got error: %s", res.Error.Message),
		)
		return
	}

	data.Bandwidth.Limit = types.Int64{Value: res.Bandwidth.Limit}
	data.Bandwidth.Usage = types.Int64{Value: res.Bandwidth.Usage}
	data.Bandwidth.UsedPercent = types.Float64{Value: res.Bandwidth.UsedPercent}
	data.DerivedResources = types.Int64{Value: int64(res.DerivedResources)}
	data.LastUpdated = types.String{Value: res.LastUpdated}
	data.MediaLimits.AssetMaxTotalPx = types.Int64{Value: int64(res.MediaLimits.AssetMaxTotalPx)}
	data.MediaLimits.ImageMaxPx = types.Int64{Value: int64(res.MediaLimits.ImageMaxPx)}
	data.MediaLimits.ImageMaxSizeBytes = types.Int64{Value: int64(res.MediaLimits.ImageMaxSizeBytes)}
	data.MediaLimits.RawMaxSizeBytes = types.Int64{Value: int64(res.MediaLimits.RawMaxSizeBytes)}
	data.MediaLimits.VideoMaxSizeBytes = types.Int64{Value: int64(res.MediaLimits.VideoMaxSizeBytes)}
	data.Objects.Limit = types.Int64{Value: int64(res.Objects.Limit)}
	data.Objects.Usage = types.Int64{Value: int64(res.Objects.Usage)}
	data.Objects.UsedPercent = types.Float64{Value: res.Objects.UsedPercent}
	data.Plan = types.String{Value: res.Plan}
	data.Requests = types.Int64{Value: res.Requests}
	data.Resources = types.Int64{Value: int64(res.Resources)}
	data.Storage.Limit = types.Int64{Value: res.Storage.Limit}
	data.Storage.Usage = types.Int64{Value: res.Storage.Usage}
	data.Storage.UsedPercent = types.Float64{Value: res.Storage.UsedPercent}
	data.Transformations.Limit = types.Int64{Value: int64(res.Transformations.Limit)}
	data.Transformations.Usage = types.Int64{Value: int64(res.Transformations.Usage)}
	data.Transformations.UsedPercent = types.Float64{Value: res.Transformations.UsedPercent}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
