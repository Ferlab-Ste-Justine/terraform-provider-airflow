package airflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (cli *Client) GetAccessToken() error {
	payload := TokenRequest{
		Username: cli.conf.Auth.Username,
		Password: cli.conf.Auth.Password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal token request body: %w", err)
	}

	authUrl, err := cli.BuildUrl("/auth/token") 
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, authUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.handle.Do(req)
	if err != nil {
		return fmt.Errorf("Token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	cli.accessToken = tokenResp.AccessToken
	
	return nil
}

func (cli *Client) setAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+cli.accessToken)
}