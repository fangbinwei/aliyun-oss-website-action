package config

import (
	"fmt"
	"os"

	"github.com/aliyun-oss-website-action/utils"
	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
)

// WebsiteOptions contains options for static website setting in OSS
type WebsiteOptions = struct {
	IndexPage    string
	NotFoundPage string
}

var (
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Folder          string
	BucketName      string
	Client          *oss.Client
	Bucket          *oss.Bucket
	Website         *WebsiteOptions
)

func init() {

	godotenv.Load(".env.local")

	Endpoint = os.Getenv("ENDPOINT")
	AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	Folder = os.Getenv("FOLDER")
	BucketName = os.Getenv("BUCKET")
	Website = &WebsiteOptions{
		IndexPage:    os.Getenv("INDEX_PAGE"),
		NotFoundPage: os.Getenv("NOT_FOUND_PAGE"),
	}

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("current directory: %s\n", currentPath)
	fmt.Printf("endpoint: %s\nbucketName: %s\nfolder: %s\nindexPage: %s\nnotFoundPage: %s\n",
		Endpoint, BucketName, Folder, Website.IndexPage, Website.NotFoundPage)

	Client, err = oss.New(Endpoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		utils.HandleError(err)
	}

	Bucket, err = Client.Bucket(BucketName)
	if err != nil {
		utils.HandleError(err)
	}
}
