package candi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var httpClient = &http.Client{}

type apiResponse struct {
	Code        int             `json:"code"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
}

func (c *Client) doRequest(httpMethod, method string, request []byte, response interface{}) error {
	endPoint := fmt.Sprintf(c.url, method)
	var resp *http.Response

	req, err := http.NewRequest(httpMethod, endPoint, strings.NewReader(string(request)))
	if err != nil {
		return err
	}
	if c.tokenType != "" {
		switch c.tokenType {
		case BEARER_AUTH_SCHEME:
			req.Header.Add("Authorization", "Bearer "+c.token)
		case BASIC_AUTH_SCHEME:
			req.SetBasicAuth(c.username, c.password)
		}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response status : %s", resp.Status)
	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	fmt.Println(string(body))

	if err != nil {
		return nil
	}

	return json.Unmarshal(body, response)
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	c, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	return cs[:s], cs[s+1:], true
}
