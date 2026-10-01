package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Ferlab-Ste-Justine/terraform-provider-airflow/airflow"
)

type AirflowFabUserDataSourceModel struct {
	Username       types.String `tfsdk:"username"`
	Email          types.String `tfsdk:"email"`
	FirstName      types.String `tfsdk:"first_name"`
	LastName       types.String `tfsdk:"last_name"`
	Roles          types.List   `tfsdk:"roles"`
	Active         types.Bool   `tfsdk:"active"`
	CreatedOn      types.String `tfsdk:"created_on"`
	ChangedOn      types.String `tfsdk:"changed_on"`
	LastLogin      types.String `tfsdk:"last_login"`
	LoginCount     types.Int64  `tfsdk:"login_count"`
	FailLoginCount types.Int64  `tfsdk:"fail_login_count"`
}

type AirflowFabUserResourceModel struct {
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Roles     types.List   `tfsdk:"roles"`
}

func NewUserResourceModelFromApi(ctx context.Context, cli *airflow.Client, username string) (*AirflowFabUserResourceModel, error) {
	user, err := cli.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	roles := convertRolesToStringList(user.Roles)

	model := &AirflowFabUserResourceModel{
		Username:  types.StringValue(user.Username),
		Email:     types.StringValue(user.Email),
		FirstName: types.StringValue(user.FirstName),
		LastName:  types.StringValue(user.LastName),
		Roles:     roles,
	}

	return model, nil
}

func (model *AirflowFabUserResourceModel) CreateUserInApi(ctx context.Context, cli *airflow.Client) error {
	roles := convertStringListToRoles(model.Roles)

	return cli.CreateUser(ctx, airflow.CreateUserRequest{
		Username:  model.Username.ValueString(),
		Password:  model.Password.ValueString(),
		Email:     model.Email.ValueString(),
		FirstName: model.FirstName.ValueString(),
		LastName:  model.LastName.ValueString(),
		Roles:     roles,
	})
}

func (model *AirflowFabUserResourceModel) UpdateUserInApi(ctx context.Context, cli *airflow.Client, preExistingUsername string) error {
	roles := convertStringListToRoles(model.Roles)

	username := model.Username.ValueString()
	password := model.Password.ValueString()
	email := model.Email.ValueString()
	firstName := model.FirstName.ValueString()
	lastName := model.LastName.ValueString()

	return cli.UpdateUser(ctx, preExistingUsername, airflow.UpdateUserRequest{
		Username:  &username,
		Password:  &password,
		Email:     &email,
		FirstName: &firstName,
		LastName:  &lastName,
		Roles:     &roles,
	})
}

func convertRolesToStringList(roles []airflow.UserRole) types.List {
	var roleStrings []attr.Value
	for _, r := range roles {
		roleStrings = append(roleStrings, types.StringValue(r.Name))
	}
	list, _ := types.ListValue(types.StringType, roleStrings)

	return list
}


func convertStringListToRoles(list types.List) []airflow.UserRole {
	var roles []airflow.UserRole
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	
	var rawRoles []attr.Value
	list.ElementsAs(context.Background(), &rawRoles, false)
	for _, v := range rawRoles {
		roles = append(roles, airflow.UserRole{Name: v.String()})
	}
	return roles
}

func NewUserDataSourceModelFromApi(ctx context.Context, cli *airflow.Client, username string) (*AirflowFabUserDataSourceModel, error) {
	user, err := cli.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	model := &AirflowFabUserDataSourceModel{
		Username:  types.StringValue(user.Username),
		Email:     types.StringValue(user.Email),
		FirstName: types.StringValue(user.FirstName),
		LastName:  types.StringValue(user.LastName),
	}
	
	roles := convertRolesToStringList(user.Roles)
	model.Roles = roles

	model.Active = types.BoolPointerValue(user.Active)

	model.CreatedOn = types.StringPointerValue(user.CreatedOn)
	model.ChangedOn = types.StringPointerValue(user.ChangedOn)
	model.LastLogin = types.StringPointerValue(user.LastLogin)

	model.LoginCount = types.Int64PointerValue(user.LoginCount)
	model.FailLoginCount = types.Int64PointerValue(user.FailLoginCount)

	return model, nil
}