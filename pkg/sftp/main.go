package sftp

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
)

type Client struct {
	url        string
	apiKey     string
	user       string
	password   string
	token      string
	httpClient *http.Client
}

func NewSftpClient(url, apiKey string) *Client {
	c := &Client{
		url:        url,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}

	return c
}

func NewSftpClientForUser(url, user, password string) *Client {
	c := &Client{
		url:        url,
		user:       user,
		password:   password,
		httpClient: &http.Client{},
	}

	err := c.getToken()
	if err != nil {
		return nil
	}

	return c
}

// token

func (c *Client) getToken() error {
	url := c.url + "/api/v2/token"

	// build request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Err(err, "get request", url)
		return err
	}

	// set basic auth header
	basicAuthString := c.user + ":" + c.password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(basicAuthString))
	req.Header.Add("Authorization", basicAuth)

	// request token
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Err(err, "get response", url)
		return err
	}
	defer resp.Body.Close()

	// read data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err, "read data", url)
		return err
	}

	// unmarshal response
	type RespBody struct {
		AccessToken string `json:"access_token"`
	}

	var respBody RespBody
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		log.Err(err, "unmarshal data", url, string(data))
		return err
	}

	// set token
	c.token = respBody.AccessToken
	return nil
}

// general

func (c *Client) RequestPath(method, path string, body io.Reader, responseBody any) ([]byte, error) {
	url := strings.TrimSuffix(c.url, "/api/v2") + "/api/v2" + strings.TrimPrefix(path, "/api/v2")

	// build request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Err(err, "request", url)
		return nil, err
	}

	if c.token != "" {
		// set bearer auth header
		bearer := "Bearer " + c.token
		req.Header.Add("Authorization", bearer)
	} else if c.apiKey != "" {
		// set api key auth header
		req.Header.Add("X-SFTPGO-API-KEY", c.apiKey)
	}

	// request token
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Err(err, "response", url)
		return nil, err
	}
	defer resp.Body.Close()

	// read data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err, "read data", url)
		return nil, err
	}

	// unmarshal response
	if responseBody != nil {
		err = json.Unmarshal(data, responseBody)
		if err != nil {
			log.Err(err, "unmarshal data", url, string(data))
			return nil, err
		}
	}

	// return body
	return data, nil
}
