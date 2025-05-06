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
	GROUPS_PATH string = "/groups"
)

type GroupSimple struct {
	Name        string
	Description string
}

func (c *Client) GetGroups(limit int) []sdk.Group {
	path := groupsPath("")
	if limit > -1 {
		path += "?" + strconv.Itoa(limit)
	}

	var groups []sdk.Group
	_, err := c.RequestPath(http.MethodGet, path, nil, &groups)
	if err != nil {
		return nil
	}
	return groups
}

func (c *Client) AddGroup(simple GroupSimple, template string) *sdk.Group {
	path := groupsPath("")

	// add group
	var group *sdk.Group
	groupReader := GetReaderFromUpdatedGroupTemplate(template, simple)
	_, err := c.RequestPath(http.MethodPost, path, groupReader, group)
	if err != nil {
		return nil
	}
	return group
}

func (c *Client) UpdateGroup(name string, patch map[string]any) (*sdk.Group, error) {
	path := groupsPath(name)

	// get current group
	var groupMap map[string]any
	_, err := c.RequestPath(http.MethodGet, path, nil, &groupMap)
	if err != nil {
		return nil, err
	}

	// patch group with map
	newGroupMap := data.PatchMap(groupMap, patch)
	data, err := json.Marshal(newGroupMap)
	if err != nil {
		return nil, err
	}

	// update group
	var group *sdk.Group
	_, err = c.RequestPath(http.MethodPut, path, bytes.NewReader(data), group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (c *Client) GetGroup(name string) *sdk.Group {
	path := groupsPath(name)

	var group *sdk.Group
	_, err := c.RequestPath(http.MethodGet, path, nil, group)
	if err != nil {
		return nil
	}
	return group
}

func (c *Client) DeleteGroup(name string) error {
	path := groupsPath(name)

	_, err := c.RequestPath(http.MethodDelete, path, nil, nil)
	return err
}

// helper

func groupsPath(suffix string) string {
	if suffix == "" {
		return GROUPS_PATH
	}
	return GROUPS_PATH + "/" + suffix
}

func UpdateGroup(group *sdk.Group, simple GroupSimple) {
	if simple.Name != "" {
		group.Name = simple.Name
	}
	if simple.Description != "" {
		group.Description = simple.Description
	}
}

func GetReaderFromUpdatedGroupTemplate(template string, simple GroupSimple) io.Reader {
	// get map from template
	groupMap, err := data.GetJson(template)
	if err != nil {
		return nil
	}

	// update
	if simple.Name != "" {
		groupMap["name"] = simple.Name
	}
	if simple.Description != "" {
		groupMap["description"] = simple.Description
	}

	// return reader
	data, err := json.Marshal(groupMap)
	if err != nil {
		return nil
	}
	return bytes.NewReader(data)
}
