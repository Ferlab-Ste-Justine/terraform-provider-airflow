package airflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type RoleActionResource struct {
	Name string `json:"name"`
}

type RoleAction struct {
	Action   RoleActionResource `json:"action"`
	Resource RoleActionResource `json:"resource"`
}

type CreateRoleRequest struct {
	Name    string       `json:"name"`
	Actions []RoleAction `json:"actions,omitempty"`
}

type UpdateRoleRequest struct {
	Name    *string       `json:"name,omitempty"`
	Actions *[]RoleAction `json:"actions,omitempty"`
}

type GetRoleResponse struct {
	Name    string       `json:"name"`
	Actions []RoleAction `json:"actions"`
}

func (cli *Client) GetRole(ctx context.Context, name string) (*GetRoleResponse, error) {
	roleUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/roles/%s", url.PathEscape(name)))
	if err != nil {
		return nil, err
	}

	cli.LogRequest(ctx, http.MethodGet, roleUrl, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, roleUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create get role request: %w", err)
	}
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Get role request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return nil, fmt.Errorf("get role request failed with status %d: %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("get role request failed with status %d", resp.StatusCode)
	}

	var result GetRoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response from get role request: %w", err)
	}

	return &result, nil
}

func (cli *Client) CreateRole(ctx context.Context, roleReq CreateRoleRequest) error {
	body, err := json.Marshal(roleReq)
	if err != nil {
		return fmt.Errorf("Failed to marshal create role request: %w", err)
	}

	rolesUrl, err := cli.BuildUrl("/auth/fab/v1/roles")
	if err != nil {
		return err
	}

	cli.LogRequest(ctx, http.MethodPost, rolesUrl, string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rolesUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to create create role request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Create role request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return fmt.Errorf("create role request failed with status %d: %s", resp.StatusCode, msg)
		}
		return fmt.Errorf("create role request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (cli *Client) UpdateRole(ctx context.Context, name string, updateReq UpdateRoleRequest) error {
	roleUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/roles/%s", url.PathEscape(name)))
	if err != nil {
		return err
	}

	body, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("Failed to marshal update role request: %w", err)
	}

	cli.LogRequest(ctx, http.MethodPatch, roleUrl, string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, roleUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to create update role request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Update role request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return fmt.Errorf("update role request failed with status %d: %s", resp.StatusCode, msg)
		}
		return fmt.Errorf("update role request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (cli *Client) DeleteRole(ctx context.Context, name string) error {
	roleUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/roles/%s", url.PathEscape(name)))
	if err != nil {
		return err
	}

	cli.LogRequest(ctx, http.MethodDelete, roleUrl, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, roleUrl, nil)
	if err != nil {
		return fmt.Errorf("Failed to create delete role request: %w", err)
	}
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Delete role request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return fmt.Errorf("delete role request failed with status %d: %s", resp.StatusCode, msg)
		}
		return fmt.Errorf("delete role request failed with status %d", resp.StatusCode)
	}

	return nil
}
