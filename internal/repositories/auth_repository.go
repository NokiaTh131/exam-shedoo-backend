package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type AuthRepository struct {
	TokenURL     string
	RedirectURL  string
	ClientID     string
	ClientSecret string
	Scope        string
	BasicInfoURL string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type CmuEntraIDBasicInfo struct {
	CmuitaccountName   string `json:"cmuitaccount_name"`
	Cmuitaccount       string `json:"cmuitaccount"`
	StudentID          string `json:"student_id"`
	FirstnameTH        string `json:"firstname_TH"`
	FirstnameEN        string `json:"firstname_EN"`
	LastnameTH         string `json:"lastname_TH"`
	LastnameEN         string `json:"lastname_EN"`
	OrganizationNameTH string `json:"organization_name_TH"`
	OrganizationNameEN string `json:"organization_name_EN"`
	ItaccounttypeID    string `json:"itaccounttype_id"`
	ItaccounttypeTH    string `json:"itaccounttype_TH"`
	ItaccounttypeEN    string `json:"itaccounttype_EN"`
}

func NewAuthRepository(tokenURL, redirectURL, clientID, clientSecret, scope, basicInfoURL string) *AuthRepository {
	return &AuthRepository{
		TokenURL:     tokenURL,
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
		BasicInfoURL: basicInfoURL,
	}
}

func (r *AuthRepository) ExchangeCode(ctx context.Context, code string) (string, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("redirect_uri", r.RedirectURL)
	form.Set("client_id", r.ClientID)
	form.Set("client_secret", r.ClientSecret)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", r.TokenURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint: %s", res.Status)
	}

	var tr TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", errors.New("empty access token")
	}
	return tr.AccessToken, nil
}

func (r *AuthRepository) GetBasicInfo(ctx context.Context, accessToken string) (*CmuEntraIDBasicInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BasicInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint: %s", res.Status)
	}

	var info CmuEntraIDBasicInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
