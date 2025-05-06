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
	FOLDERS_PATH string = "/folders"
)

type FolderSimple struct {
	Name        string
	Description string
	Path        string
}

func (c *Client) GetFolders(limit int) []sdk.BaseVirtualFolder {
	path := foldersPath("")
	if limit > -1 {
		path += "?" + strconv.Itoa(limit)
	}

	var folders []sdk.BaseVirtualFolder
	_, err := c.RequestPath(http.MethodGet, path, nil, &folders)
	if err != nil {
		return nil
	}
	return folders
}

func (c *Client) AddFolder(simple FolderSimple, template string) *sdk.BaseVirtualFolder {
	path := foldersPath("")

	// add folder
	var folder *sdk.BaseVirtualFolder
	folderReader := GetReaderFromUpdatedFolderTemplate(template, simple)
	_, err := c.RequestPath(http.MethodPost, path, folderReader, folder)
	if err != nil {
		return nil
	}
	return folder
}

func (c *Client) UpdateFolder(name string, patch map[string]any) (*sdk.BaseVirtualFolder, error) {
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
	var folder *sdk.BaseVirtualFolder
	_, err = c.RequestPath(http.MethodPut, path, bytes.NewReader(data), folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (c *Client) GetFolder(name string) *sdk.BaseVirtualFolder {
	path := foldersPath(name)

	var folder *sdk.BaseVirtualFolder
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

func UpdateFolder(folder *sdk.BaseVirtualFolder, simple FolderSimple) {
	if simple.Name != "" {
		folder.Name = simple.Name
	}
	if simple.Description != "" {
		folder.Description = simple.Description
	}
	if simple.Path != "" {
		folder.MappedPath = simple.Path
	}
}

func GetReaderFromUpdatedFolderTemplate(template string, simple FolderSimple) io.Reader {
	// get map from template
	folderMap, err := data.GetJson(template)
	if err != nil {
		return nil
	}

	// update
	if simple.Name != "" {
		folderMap["name"] = simple.Name
	}
	if simple.Description != "" {
		folderMap["description"] = simple.Description
	}
	if simple.Path != "" {
		folderMap["mapped_path"] = simple.Path
	}

	// return reader
	data, err := json.Marshal(folderMap)
	if err != nil {
		return nil
	}
	return bytes.NewReader(data)
}
