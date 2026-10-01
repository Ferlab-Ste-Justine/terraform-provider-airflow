package airflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getErrorBody(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// Try to parse as JSON and look for common error fields
	var errResp map[string]interface{}
	if json.Unmarshal(body, &errResp) == nil {
		if d, ok := errResp["detail"]; ok {
			return fmt.Sprintf("%v", d)
		}
		if m, ok := errResp["message"]; ok {
			return fmt.Sprintf("%v", m)
		}
		if e, ok := errResp["error"]; ok {
			return fmt.Sprintf("%v", e)
		}
	}

	// Fallback to raw body
	return string(body)
}
