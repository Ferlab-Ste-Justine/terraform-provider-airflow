package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

var _ resource.Resource = (*AirflowFabRoleResource)(nil)

type RequireCanReadWebsiteInRoleValidator struct{}

func (val RequireCanReadWebsiteInRoleValidator) Description(_ context.Context) string {
	return "Validates that the permissions set contains 'can_read' on 'Website'."
}

func (val RequireCanReadWebsiteInRoleValidator) MarkdownDescription(ctx context.Context) string {
	return val.Description(ctx)
}

func (val RequireCanReadWebsiteInRoleValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	found := false
	for _, elem := range req.ConfigValue.Elements() {
		obj, ok := elem.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}

		var perm PermissionModel
		resp.Diagnostics.Append(obj.As(ctx, &perm, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		if perm.Action.IsUnknown() || perm.Resource.IsUnknown() {
			continue
		} 

		if perm.Action.ValueString() == "can_read" && perm.Resource.ValueString() == "Website" {
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Permission",
			"Airflow 3 automatically assigns 'can_read' on 'Website' to all custom roles. This permission must be explicitly included in the configuration of your roles to be consistent with what is in airflow.",
		)
	}
}

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
		Description: "An Airflow FAB role.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the role.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"permissions": schema.SetNestedAttribute{
				Description: "A list of permissions (action/resource pairs) assigned to the role. Note that because of limitations with airflow's PATCH api, changing this is a replacement operation. The role will be re-created and needs to be re-assigned to users. You need to grant `can_read` on `Website` to all custom roles. Airflow implicitly adds it and this provider enforces it.",
				Required: true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					RequireCanReadWebsiteInRoleValidator{},
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Description: "The name of the action (e.g., can_read, can_edit).",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"resource": schema.StringAttribute{
							Description: "The name of the resource (e.g., DAGs, Dashboard).",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
