package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

var _ resource.Resource = (*AirflowFabUserResource)(nil)

type AirflowFabUserResource struct {
	client *airflow.Client
}

func NewAirflowFabUserResource() resource.Resource {
	return &AirflowFabUserResource{}
}

func (r *AirflowFabUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fab_user"
}

func (r *AirflowFabUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Airflow FAB users.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{Required: true},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"email":      schema.StringAttribute{Required: true},
			"first_name": schema.StringAttribute{Required: true},
			"last_name":  schema.StringAttribute{Required: true},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *AirflowFabUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = req.ProviderData.(*airflow.Client)
}

func (r *AirflowFabUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := plan.CreateUserInApi(r.client)
	if err != nil {
		resp.Diagnostics.AddError("Creation Failed", err.Error())
		return
	}

	state, diagErr := NewUserResourceModelFromApi(r.client, plan.Username.ValueString())
	if diagErr != nil {
		resp.Diagnostics.AddError("State Refresh Failed", diagErr.Error())
		return
	}
	state.Password = plan.Password

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AirflowFabUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	
	newState, err := NewUserResourceModelFromApi(r.client, state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}

	newState.Password = state.Password
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *AirflowFabUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := plan.UpdateUserInApi(r.client, state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", err.Error())
		return
	}

	newState, err := NewUserResourceModelFromApi(r.client, plan.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("State Refresh Failed", err.Error())
		return
	}
	newState.Password = plan.Password

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *AirflowFabUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := r.client.DeleteUser(state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Failed", err.Error())
		return
	}
}