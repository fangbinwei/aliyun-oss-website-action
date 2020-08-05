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
	IndexPage       string
	NotFoundPage    string
)

func init() {

	godotenv.Load(".env.local")

	Endpoint = os.Getenv("ENDPOINT")
	AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	Folder = os.Getenv("FOLDER")
	BucketName = os.Getenv("BUCKET")
	IndexPage = os.Getenv("IndexPage")
	NotFoundPage = os.Getenv("NotFoundPage")

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("current directory: %s\n", currentPath)
	fmt.Printf("endpoint: %s\nbucketName: %s\nfolder: %s\nindexPage: %s\nnotFoundPage: %s\n",
		Endpoint, BucketName, Folder, IndexPage, NotFoundPage)
}
