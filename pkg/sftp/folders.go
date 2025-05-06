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
	FOLDERS_PATH string = "/folders"
)

func (c *Client) GetFolders(limit int) []sdk.VirtualFolder {
	path := foldersPath("")
	if limit > -1 {
		path += "?" + strconv.Itoa(limit)
	}

	var folders []sdk.VirtualFolder
	_, err := c.RequestPath(http.MethodGet, path, nil, &folders)
	if err != nil {
		return nil
	}
	return folders
}

func (c *Client) AddFolder(template string, patch map[string]any) (*sdk.VirtualFolder, error) {
	path := foldersPath("")

	// patch
	folderMap, _ := data.GetJson(template)
	newFolderMap := data.PatchMap(folderMap, patch)
	data, err := json.Marshal(newFolderMap)
	if err != nil {
		return nil, err
	}

	// add folder
	var folder *sdk.VirtualFolder
	_, err = c.RequestPath(http.MethodPost, path, bytes.NewReader(data), folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (c *Client) UpdateFolder(name string, patch map[string]any) (*sdk.VirtualFolder, error) {
	path := foldersPath(name)

	// get current folder
	var folderMap map[string]any
	_, err := c.RequestPath(http.MethodGet, path, nil, &folderMap)
	if err != nil {
		return nil, err
	}

	// patch folder with map
	newFolderMap := data.PatchMap(folderMap, patch)
	data, err := json.Marshal(newFolderMap)
	if err != nil {
		return nil, err
	}

	// update folder
	var folder *sdk.VirtualFolder
	_, err = c.RequestPath(http.MethodPut, path, bytes.NewReader(data), folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (c *Client) GetFolder(name string) *sdk.VirtualFolder {
	path := foldersPath(name)

	var folder *sdk.VirtualFolder
	_, err := c.RequestPath(http.MethodGet, path, nil, folder)
	if err != nil {
		return nil
	}
	return folder
}

func (c *Client) DeleteFolder(name string) error {
	path := foldersPath(name)

	_, err := c.RequestPath(http.MethodDelete, path, nil, nil)
	return err
}

// helper

func foldersPath(suffix string) string {
	if suffix == "" {
		return FOLDERS_PATH
	}
	return FOLDERS_PATH + "/" + suffix
}
