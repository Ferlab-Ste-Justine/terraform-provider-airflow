package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

var _ provider.Provider = AirflowProvider{}

type AirflowProvider struct {
	version string
}

type AirflowProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Address  types.String `tfsdk:"address"`
	CaCert   types.String `tfsdk:"ca_cert"`
}

func (p AirflowProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "airflow"
}

func (p AirflowProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform Provider for Airflow.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for Airflow API authentication.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for Airflow API authentication.",
				Required:            true,
				Sensitive:           true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The base URL of the Airflow API server.",
				Required:            true,
			},
			"ca_cert": schema.StringAttribute{
				MarkdownDescription: "Path to a CA certificate file to verify the Airflow API server's TLS certificate.",
				Optional:            true,
			},
		},
	}
}

func (p AirflowProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AirflowProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, clientErr := airflow.NewClient(&airflow.ClientConfig{
		ServerAddress: data.Address.ValueString(),
		Auth: airflow.Auth{
			CaCert:   data.CaCert.ValueString(),
			Username: data.Username.ValueString(),
			Password: data.Password.ValueString(),
		},
		ConnectionTimeout: 5 * time.Minute,
		RequestTimeout:    5 * time.Minute,
	})
	if clientErr != nil {
		resp.Diagnostics.AddError(
			"Error initializing Airflow client",
			clientErr.Error(),
		)
		return
	}

	accessTokenErr := client.GetAccessToken()
	if accessTokenErr != nil {
		resp.Diagnostics.AddError(
			"Error fetching access token for Airflow client",
			accessTokenErr.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p AirflowProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAirflowFabUserResource,
		NewAirflowFabRoleResource,
	}
}

func (p AirflowProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAirflowFabUserDataSource,
		NewAirflowFabRoleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return AirflowProvider{
			version: version,
		}
	}
}
