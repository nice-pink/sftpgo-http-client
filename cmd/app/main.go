package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/sftpgo-http-client/pkg/sftp"
)

func main() {
	// config := flag.String("config", "", "Path to config file.")
	url := flag.String("url", "http://localhost:8081", "Url of sftpgo rest api.")
	user := flag.String("user", "admin", "User name for loging.")
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	password := os.Getenv("SFTPGO_PASSWORD")
	client := sftp.NewSftpClientForUser(*url, *user, password)
	if client == nil {
		log.Error("Could not get sftp client!")
		os.Exit(2)
	}

	// users
	log.Info("Users:")
	users := client.GetUsers(-1)
	for _, u := range users {
		log.Info("-", u.BaseUser.Username)
	}
	if users != nil {
		data, _ := json.MarshalIndent(users[1], "", "  ")
		log.Info(string(data))
	}

	// groups
	log.Info("Groups:")
	groups := client.GetGroups(-1)
	for _, g := range groups {
		log.Info("-", g.Name, g.ID)
	}
	if groups != nil {
		data, _ := json.MarshalIndent(groups[7], "", "  ")
		log.Info(string(data))
	}

	// folders
	log.Info("Folders:")
	folders := client.GetFolders(-1)
	for _, f := range folders {
		log.Info("-", f.Name)
	}
	if folders != nil {
		data, _ := json.MarshalIndent(folders[5], "", "  ")
		log.Info(string(data))
	}

	log.Info("Folder:")
	folder := client.GetFolder(folders[5].Name)
	log.Info(folder)
}
