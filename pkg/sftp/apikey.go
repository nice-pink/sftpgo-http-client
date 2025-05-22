package sftp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nice-pink/goutil/pkg/data"
)

type ApiKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Scope       int    `json:"scope"` // 1: admin, 2: user
	Description string `json:"description"`
	User        string `json:"user"`  // leave empty to avoid user binding
	Admin       string `json:"admin"` // leave empty to avoid admin binding (any admin)
	// CreatedAt   int64  `json:"created_at"`
	// UpdatedAt   int64  `json:"updated_at"`
	// LastUseAt   int64  `json:"last_use_at"`
	// ExpiresAt   int64  `json:"expires_at"`
}

const (
	APIKEY_PATH string = "/apikeys"
)

func (c *Client) GetApiKeys(limit int) []ApiKey {
	path := apiKeysPath("")
	if limit > -1 {
		path += "?" + strconv.Itoa(limit)
	}

	var keys []ApiKey
	_, err := c.RequestPath(http.MethodGet, path, nil, &keys)
	if err != nil {
		return nil
	}
	return keys
}

func (c *Client) AddApiKey(template string, patch map[string]any) (*ApiKey, error) {
	path := apiKeysPath("")

	// patch
	keyMap, _ := data.GetJsonMap(template)
	keyMap = data.PatchMap(keyMap, patch)
	data, err := json.Marshal(keyMap)
	if err != nil {
		return nil, err
	}

	// add key
	var key ApiKey
	_, err = c.RequestPath(http.MethodPost, path, bytes.NewReader(data), &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *Client) UpdateApiKey(id string, patch map[string]any) (*ApiKey, error) {
	path := apiKeysPath(id)

	// get current api key
	var keyMap map[string]any
	_, err := c.RequestPath(http.MethodGet, path, nil, &keyMap)
	if err != nil {
		return nil, err
	}

	// patch api key with map
	keyMap = data.PatchMap(keyMap, patch)
	data, err := json.Marshal(keyMap)
	if err != nil {
		return nil, err
	}

	// update key
	var key *ApiKey
	_, err = c.RequestPath(http.MethodPut, path, bytes.NewReader(data), key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (c *Client) GetApiKey(id string) *ApiKey {
	path := apiKeysPath(id)

	// get key
	var key ApiKey
	_, err := c.RequestPath(http.MethodGet, path, nil, &key)
	if err != nil {
		return nil
	}
	return &key
}

func (c *Client) DeleteApiKey(id string) error {
	path := apiKeysPath(id)

	// delete key
	_, err := c.RequestPath(http.MethodDelete, path, nil, nil)
	return err
}

// helper

func apiKeysPath(suffix string) string {
	if suffix == "" {
		return APIKEY_PATH
	}
	return APIKEY_PATH + "/" + UrlEncode(suffix)
}
