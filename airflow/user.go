package airflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type UserRole struct {
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	FirstName  string       `json:"first_name"`
	LastName   string       `json:"last_name"`
	Roles      []UserRole   `json:"roles"`
	Password   string       `json:"password"`
}

type UpdateUserRequest struct {
	Username   *string     `json:"username,omitempty"`
	Email      *string     `json:"email,omitempty"`
	FirstName  *string     `json:"first_name,omitempty"`
	LastName   *string     `json:"last_name,omitempty"`
	Roles      *[]UserRole `json:"roles,omitempty"`
	Password   *string     `json:"password,omitempty"`
}

type GetUserResponse struct {
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Roles            []UserRole `json:"roles"`
	Active           *bool      `json:"active"`
	LastLogin        *string    `json:"last_login"` // date-time
	LoginCount       *int64     `json:"login_count"`
	FailLoginCount   *int64     `json:"fail_login_count"`
	CreatedOn        *string    `json:"created_on"` // date-time
	ChangedOn        *string    `json:"changed_on"` // date-time
}

func (cli *Client) CreateUser(userReq CreateUserRequest) error {
	body, err := json.Marshal(userReq)
	if err != nil {
		return fmt.Errorf("Failed to marshal create user request: %w", err)
	}

	usersUrl, err := cli.BuildUrl("/auth/fab/v1/users") 
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, usersUrl, bytes.NewBuffer(body))
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
		return fmt.Errorf("create user request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (cli *Client) GetUser(username string) (*GetUserResponse, error) {
	userUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/users/%s", url.PathEscape(username))) 
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, userUrl, nil)
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
		return nil, fmt.Errorf("Get user request failed with status: %d", resp.StatusCode)
	}

	var result GetUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response from get user request: %w", err)
	}

	return &result, nil
}

func (cli *Client) UpdateUser(username string, updateReq UpdateUserRequest) error {
	userUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/users/%s", url.PathEscape(username))) 
	if err != nil {
		return err
	}

	body, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("Failed to marshal update user request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, userUrl, bytes.NewBuffer(body))
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
		return fmt.Errorf("Update user request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (cli *Client) DeleteUser(username string) error {
	userUrl, err := cli.BuildUrl(fmt.Sprintf("/auth/fab/v1/users/%s", url.PathEscape(username))) 
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, userUrl, nil)
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
		return fmt.Errorf("Delete user request failed with status: %d", resp.StatusCode)
	}

	return nil
}