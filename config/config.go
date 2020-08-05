package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Folder          string
	BucketName      string
)

func init() {

	godotenv.Load(".env.local")

	Endpoint = os.Getenv("ENDPOINT")
	AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	Folder = os.Getenv("FOLDER")
	BucketName = os.Getenv("BUCKET")
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("endpoint: %s\nbucketName: %s\nfolder: %s\ncurrent directory: %s\n", Endpoint, BucketName, Folder, currentPath)
}
