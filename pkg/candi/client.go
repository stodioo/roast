package candi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	BEARER_AUTH_SCHEME = "Bearer"
	BASIC_AUTH_SCHEME  = "Default"
)

type Client struct {
	tokenType string
	token     string
	url       string
	username  string
	password  string
}

func NewClient(envPrefixVar, basePath string) (*Client, error) {
	cfg, err := parseFromEnv(envPrefixVar)
	if err != nil {
		return nil, err
	}
	cl := &Client{
		url: fmt.Sprintf("%s/%s", cfg.ApiBaseURL, basePath) + "%s",
	}
	if cfg.BasicAuth != "" {
		username, password, ok := parseBasicAuth(cfg.BasicAuth)
		if ok {
			cl.username = username
			cl.password = password
		}
	}
	return cl, nil
}

type ClientOption func(client *Client)

func (cl *Client) WithBearerToken(token string) {
	cl.token = token
}

func (cl *Client) TokenType(tokenType string) {
	cl.tokenType = tokenType
}

func (cl *Client) GetUserInfo(userID string) (*User, error) {
	user := &User{}
	err := cl.doRequest(http.MethodGet, "/users/"+userID, nil, user)

	return user, err
}

func (cl *Client) GetUserByPhoneNumber(phoneNumbers string) (*UserPhoneNumberList, error) {
	userPhoneNumberList := &UserPhoneNumberList{}
	req := url.Values{}
	req.Set("phone_numbers", phoneNumbers)
	query := req.Encode()
	err := cl.doRequest(http.MethodGet, "/users/by_phone_numbers?"+query, nil, userPhoneNumberList)

	return userPhoneNumberList, err
}

func (cl *Client) Register(val interface{}) (*TerminalRegisterPostResponse, error) {
	terminalRegister := &TerminalRegisterPostResponse{}
	req, _ := json.Marshal(val)

	err := cl.doRequest(http.MethodPost, "/terminals/register", req, terminalRegister)
	return terminalRegister, err
}
