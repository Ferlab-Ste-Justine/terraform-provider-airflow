package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

var _ datasource.DataSource = (*AirflowFabUserDataSource)(nil)

type AirflowFabUserDataSource struct {
	client *airflow.Client
}

func NewAirflowFabUserDataSource() datasource.DataSource {
	return &AirflowFabUserDataSource{}
}

func (d *AirflowFabUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fab_user"
}

func (d *AirflowFabUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches an Airflow FAB user.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{Required: true},
			"email":    schema.StringAttribute{Computed: true},
			"first_name": schema.StringAttribute{Computed: true},
			"last_name": schema.StringAttribute{Computed: true},
			"roles": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"active": schema.BoolAttribute{Computed: true},
			"created_on": schema.StringAttribute{Computed: true},
			"changed_on": schema.StringAttribute{Computed: true},
			"last_login": schema.StringAttribute{Computed: true},
			"login_count": schema.Int64Attribute{Computed: true},
			"fail_login_count": schema.Int64Attribute{Computed: true},
		},
	}
}

func (d *AirflowFabUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*airflow.Client)
}

func (d *AirflowFabUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var config AirflowFabUserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	state, err := NewUserDataSourceModelFromApi(ctx, d.client, config.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
