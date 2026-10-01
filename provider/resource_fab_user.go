package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

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
		Description: "An Airflow FAB user.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "User id for dashboard login and api token generation.",
				Required: true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
			},
			"password": schema.StringAttribute{
				Description: "User authentication secret for dashboard login and api token generation.",
				Required:  true,
				Sensitive: true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
			},
			"email":      schema.StringAttribute{
				Description: "User's email.",
				Required: true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
			},
			"first_name": schema.StringAttribute{
				Description: "User's first name.",
				Required: true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
			},
			"last_name":  schema.StringAttribute{
				Description: "User's last name.",
				Required: true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
			},
			"roles": schema.SetAttribute{
				Description: "Set of user's role. Note that because changing a role's permission causes a recreation and this field is evaluate for reconciliation during the plan phase, changing a role's permission and assigning it to a user in the same terraform apply will cause the user to temporarily lose the role until another terraform apply is performed.",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
		},
	}
}

func (r *AirflowFabUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*airflow.Client)
}

func (r *AirflowFabUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var plan AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := plan.CreateUserInApi(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Creation Failed", err.Error())
		return
	}

	state, diagErr := NewUserResourceModelFromApi(ctx, r.client, plan.Username.ValueString())
	if diagErr != nil {
		resp.Diagnostics.AddError("State Refresh Failed", diagErr.Error())
		return
	}
	state.Password = plan.Password

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AirflowFabUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var state AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	
	newState, err := NewUserResourceModelFromApi(ctx, r.client, state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}

	newState.Password = state.Password
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *AirflowFabUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var plan, state AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := plan.UpdateUserInApi(ctx, r.client, state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", err.Error())
		return
	}

	newState, err := NewUserResourceModelFromApi(ctx, r.client, plan.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("State Refresh Failed", err.Error())
		return
	}
	newState.Password = plan.Password

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *AirflowFabUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Airflow Client",
			"Expected configured Airflow client. Please report this issue to the provider developers.",
		)
		return
	}

	var state AirflowFabUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := r.client.DeleteUser(ctx, state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Failed", err.Error())
		return
	}
}
