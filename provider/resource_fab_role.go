package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

var _ resource.Resource = (*AirflowFabRoleResource)(nil)

type AirflowFabRoleResource struct {
	client *airflow.Client
}

func NewAirflowFabRoleResource() resource.Resource {
	return &AirflowFabRoleResource{}
}

func (r *AirflowFabRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fab_role"
}

func (r *AirflowFabRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Airflow FAB roles.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"permissions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "A list of permissions (action/resource pairs) assigned to the role.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Required:    true,
							Description: "The name of the action (e.g., can_read, can_edit).",
						},
						"resource": schema.StringAttribute{
							Required:    true,
							Description: "The name of the resource (e.g., DAGs, Dashboard).",
						},
					},
				},
			},
		},
	}
}

func (r *AirflowFabRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*airflow.Client)
}

func (r *AirflowFabRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var plan AirflowFabRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := plan.CreateRoleInApi(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Creation Failed", err.Error())
		return
	}

	state, diagErr := NewRoleResourceModelFromApi(ctx, r.client, plan.Name.ValueString())
	if diagErr != nil {
		resp.Diagnostics.AddError("State Refresh Failed", diagErr.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AirflowFabRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var state AirflowFabRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	newState, err := NewRoleResourceModelFromApi(ctx, r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *AirflowFabRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var plan, state AirflowFabRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := plan.UpdateRoleInApi(ctx, r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", err.Error())
		return
	}

	newState, err := NewRoleResourceModelFromApi(ctx, r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("State Refresh Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *AirflowFabRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var state AirflowFabRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := r.client.DeleteRole(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Failed", err.Error())
		return
	}
}
