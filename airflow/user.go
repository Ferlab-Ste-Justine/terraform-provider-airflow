package airflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type UserRole struct {
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Roles     []UserRole `json:"roles"`
	Password  string     `json:"password"`
}

type UpdateUserRequest struct {
	Username  *string     `json:"username,omitempty"`
	Email     *string     `json:"email,omitempty"`
	FirstName *string     `json:"first_name,omitempty"`
	LastName  *string     `json:"last_name,omitempty"`
	Roles     *[]UserRole `json:"roles,omitempty"`
	Password  *string     `json:"password,omitempty"`
}

type GetUserResponse struct {
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Roles          []UserRole `json:"roles"`
	Active         *bool      `json:"active"`
	LastLogin      *string    `json:"last_login"`
	LoginCount     *int64     `json:"login_count"`
	FailLoginCount *int64     `json:"fail_login_count"`
	CreatedOn      *string    `json:"created_on"`
	ChangedOn      *string    `json:"changed_on"`
}

func (cli *Client) CreateUser(ctx context.Context, userReq CreateUserRequest) error {
	body, err := json.Marshal(userReq)
	if err != nil {
		return fmt.Errorf("Failed to marshal create user request: %w", err)
	}

	usersUrl, err := cli.BuildUrl("/auth/fab/v1/users")
	if err != nil {
		return err
	}

	cli.LogRequest(ctx, http.MethodPost, usersUrl, string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, usersUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to create create user request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Create user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return fmt.Errorf("create user request failed with status %d: %s", resp.StatusCode, msg)
		}
		return fmt.Errorf("create user request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (cli *Client) GetUser(ctx context.Context, username string) (*GetUserResponse, error) {
	userUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/users/%s", url.PathEscape(username)))
	if err != nil {
		return nil, err
	}

	cli.LogRequest(ctx, http.MethodGet, userUrl, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create get user request: %w", err)
	}
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Get user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return nil, fmt.Errorf("get user request failed with status %d: %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("get user request failed with status %d", resp.StatusCode)
	}

	var result GetUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response from get user request: %w", err)
	}

	return &result, nil
}

func (cli *Client) UpdateUser(ctx context.Context, username string, updateReq UpdateUserRequest) error {
	userUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/users/%s", url.PathEscape(username)))
	if err != nil {
		return err
	}

	body, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("Failed to marshal update user request: %w", err)
	}

	cli.LogRequest(ctx, http.MethodPatch, userUrl, string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, userUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to create update user request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Update user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return fmt.Errorf("update user request failed with status %d: %s", resp.StatusCode, msg)
		}
		return fmt.Errorf("update user request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (cli *Client) DeleteUser(ctx context.Context, username string) error {
	userUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/users/%s", url.PathEscape(username)))
	if err != nil {
		return err
	}

	cli.LogRequest(ctx, http.MethodDelete, userUrl, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, userUrl, nil)
	if err != nil {
		return fmt.Errorf("Failed to create delete user request: %w", err)
	}
	cli.setAuthHeader(req)

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Delete user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := getErrorBody(resp)
		if msg != "" {
			return fmt.Errorf("delete user request failed with status %d: %s", resp.StatusCode, msg)
		}
		return fmt.Errorf("delete user request failed with status %d", resp.StatusCode)
	}

	return nil
}
