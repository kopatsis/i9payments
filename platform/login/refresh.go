package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type tokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type tokenResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

// GetNewIDToken exchanges a refresh token for a new ID token and refresh token.
func GetNewIDToken(refreshToken string) (string, string, error) {

	apiKey := os.Getenv("ADMIN_API")

	url := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s", apiKey)

	reqBody := tokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("error response from server: %s", body)
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", "", err
	}

	return tokenResp.IDToken, tokenResp.RefreshToken, nil
}
