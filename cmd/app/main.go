package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/sftpgo-http-client/pkg/sftp"
)

func main() {
	// config := flag.String("config", "", "Path to config file.")
	url := flag.String("url", "http://localhost:8080", "Url of sftpgo rest api.")
	user := flag.String("user", "admin", "User name for loging.")
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	password := os.Getenv("SFTPGO_PASSWORD")
	client := sftp.NewSftpClientForUser(*url, *user, password)
	users := client.GetUsers(-1)
	for _, u := range users {
		log.Info(u.BaseUser.Username)
	}
}
