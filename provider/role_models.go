package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

type PermissionModel struct {
	Action   types.String `tfsdk:"action"`
	Resource types.String `tfsdk:"resource"`
}

type AirflowFabRoleResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Permissions types.List   `tfsdk:"permissions"`
}

type AirflowFabRoleDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Permissions types.List   `tfsdk:"permissions"`
}

func NewRoleResourceModelFromApi(ctx context.Context, cli *airflow.Client, name string) (*AirflowFabRoleResourceModel, error) {
	role, err := cli.GetRole(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch role: %w", err)
	}

	model := &AirflowFabRoleResourceModel{
		Name:        types.StringValue(role.Name),
		Permissions: apiActionsToTfList(role.Actions),
	}

	return model, nil
}

func NewRoleDataSourceModelFromApi(ctx context.Context, cli *airflow.Client, name string) (*AirflowFabRoleDataSourceModel, error) {
	role, err := cli.GetRole(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch role: %w", err)
	}

	model := &AirflowFabRoleDataSourceModel{
		Name:        types.StringValue(role.Name),
		Permissions: apiActionsToTfList(role.Actions),
	}

	return model, nil
}

func (model *AirflowFabRoleResourceModel) CreateRoleInApi(ctx context.Context, cli *airflow.Client) error {
	return cli.CreateRole(ctx, airflow.CreateRoleRequest{
		Name:    model.Name.ValueString(),
		Actions: tfListToApiActions(model.Permissions),
	})
}

func (model *AirflowFabRoleResourceModel) UpdateRoleInApi(ctx context.Context, cli *airflow.Client, preExistingName string) error {
	name := model.Name.ValueString()
	actions := tfListToApiActions(model.Permissions)

	return cli.UpdateRole(ctx, preExistingName, airflow.UpdateRoleRequest{
		Name:    &name,
		Actions: &actions,
	})
}

func apiActionsToTfList(actions []airflow.RoleAction) types.List {
	if len(actions) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"action":   types.StringType,
				"resource": types.StringType,
			},
		})
	}

	attrTypes := map[string]attr.Type{
		"action":   types.StringType,
		"resource": types.StringType,
	}

	var attrValues []attr.Value
	for _, a := range actions {
		objVal, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
			"action":   types.StringValue(a.Action.Name),
			"resource": types.StringValue(a.Resource.Name),
		})
		attrValues = append(attrValues, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, attrValues)
	return list
}

func tfListToApiActions(list types.List) []airflow.RoleAction {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var permissions []PermissionModel
	list.ElementsAs(context.Background(), &permissions, false)

	var apiActions []airflow.RoleAction
	for _, p := range permissions {
		apiActions = append(apiActions, airflow.RoleAction{
			Action:   airflow.RoleActionResource{Name: p.Action.ValueString()},
			Resource: airflow.RoleActionResource{Name: p.Resource.ValueString()},
		})
	}
	return apiActions
}
