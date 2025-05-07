package sftp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nice-pink/goutil/pkg/data"
	"github.com/sftpgo/sdk"
)

const (
	USERS_PATH string = "/users"
)

func (c *Client) GetUsers(limit int) []sdk.User {
	path := usersPath("")
	if limit > -1 {
		path += "?" + strconv.Itoa(limit)
	}

	var users []sdk.User
	_, err := c.RequestPath(http.MethodGet, path, nil, &users)
	if err != nil {
		return nil
	}
	return users
}

func (c *Client) AddUser(template string, patch map[string]any) (*sdk.User, error) {
	path := usersPath("")

	// patch
	userMap, _ := data.GetJson(template)
	newUserMap := data.PatchMap(userMap, patch)
	data, err := json.Marshal(newUserMap)
	if err != nil {
		return nil, err
	}

	// add user
	var user *sdk.User
	_, err = c.RequestPath(http.MethodPost, path, bytes.NewReader(data), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) UpdateUser(username string, patch map[string]any) (*sdk.User, error) {
	path := usersPath(username)

	// get current user
	var userMap map[string]any
	_, err := c.RequestPath(http.MethodGet, path, nil, &userMap)
	if err != nil {
		return nil, err
	}

	// patch user with map
	newUserMap := data.PatchMap(userMap, patch)
	data, err := json.Marshal(newUserMap)
	if err != nil {
		return nil, err
	}

	// update user
	var user *sdk.User
	_, err = c.RequestPath(http.MethodPut, path, bytes.NewReader(data), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) GetUser(username string) *sdk.User {
	path := usersPath(username)

	// get user
	var user *sdk.User
	_, err := c.RequestPath(http.MethodGet, path, nil, user)
	if err != nil {
		return nil
	}
	return user
}

func (c *Client) DeleteUser(username string) error {
	path := usersPath(username)

	// delete user
	_, err := c.RequestPath(http.MethodDelete, path, nil, nil)
	return err
}

// helper

func usersPath(suffix string) string {
	if suffix == "" {
		return USERS_PATH
	}
	return USERS_PATH + "/" + suffix
}
