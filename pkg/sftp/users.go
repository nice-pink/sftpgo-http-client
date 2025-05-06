package sftp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/nice-pink/goutil/pkg/data"
	"github.com/sftpgo/sdk"
)

const (
	USERS_PATH string = "/users"
)

type UserSimple struct {
	Username    string
	Password    string
	Email       string
	Description string
}

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

func (c *Client) AddUser(simple UserSimple, template string) *sdk.User {
	path := usersPath("")

	// add user
	var user *sdk.User
	userReader := GetReaderFromUpdatedTemplate(template, simple)
	_, err := c.RequestPath(http.MethodPost, path, userReader, user)
	if err != nil {
		return nil
	}
	return user
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

	// add user
	var user *sdk.User
	_, err := c.RequestPath(http.MethodGet, path, nil, user)
	if err != nil {
		return nil
	}
	return user
}

func (c *Client) DeleteUser(username string) error {
	path := usersPath(username)

	// add user
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

func UpdateUser(user *sdk.User, simple UserSimple) {
	if simple.Username != "" {
		user.Username = simple.Username
	}
	if simple.Password != "" {
		user.Password = simple.Password
	}
	if simple.Email != "" {
		user.Email = simple.Email
	}
	if simple.Description != "" {
		user.Description = simple.Description
	}
}

func GetReaderFromUpdatedTemplate(template string, simple UserSimple) io.Reader {
	// get map from template
	userMap, err := data.GetJson(template)
	if err != nil {
		return nil
	}

	// update
	if simple.Username != "" {
		userMap["username"] = simple.Username
	}
	if simple.Password != "" {
		userMap["password"] = simple.Password
	}
	if simple.Email != "" {
		userMap["email"] = simple.Email
	}
	if simple.Description != "" {
		userMap["description"] = simple.Description
	}

	// return reader
	data, err := json.Marshal(userMap)
	if err != nil {
		return nil
	}
	return bytes.NewReader(data)
}
