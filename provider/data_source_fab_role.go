package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

var _ datasource.DataSource = (*AirflowFabRoleDataSource)(nil)

type AirflowFabRoleDataSource struct {
	client *airflow.Client
}

func NewAirflowFabRoleDataSource() datasource.DataSource {
	return &AirflowFabRoleDataSource{}
}

func (d *AirflowFabRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fab_role"
}

func (d *AirflowFabRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches an Airflow FAB role.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Required: true},
			"permissions": schema.SetNestedAttribute{
				Computed:    true,
				Description: "A list of permissions (action/resource pairs) assigned to the role.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the action (e.g., can_read, can_edit).",
						},
						"resource": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the resource (e.g., DAGs, Dashboard).",
						},
					},
				},
			},
		},
	}
}

func (d *AirflowFabRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*airflow.Client)
}

func (d *AirflowFabRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}
	
	var config AirflowFabRoleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	state, err := NewRoleDataSourceModelFromApi(ctx, d.client, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
