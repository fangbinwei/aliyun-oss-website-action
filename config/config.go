package config

import (
	"encoding/json"
	"fmt"
	"os"

	"aliyun-oss-website-action/utils"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
)

var (
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Folder          string
	Exclude         []string
	BucketName      string
	IsCname         bool
	Client          *oss.Client
	Bucket          *oss.Bucket
	SkipSetting     bool
	IsIncremental   bool

	IndexPage    string
	NotFoundPage string
	Headers      []utils.HeadersConfig
)

func init() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")

	Endpoint = os.Getenv("ENDPOINT")
	IsCname = os.Getenv("CNAME") == "true"
	AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	Folder = os.Getenv("FOLDER")
	Exclude = utils.GetActionInputAsSlice(os.Getenv("EXCLUDE"))
	BucketName = os.Getenv("BUCKET")
	SkipSetting = os.Getenv("SKIP_SETTING") == "true"
	IsIncremental = os.Getenv("INCREMENTAL") == "true"

	IndexPage = utils.Getenv("INDEX_PAGE", "index.html")
	NotFoundPage = utils.Getenv("NOT_FOUND_PAGE", "404.html")

	headersEnv := utils.Getenv("HEADERS", utils.DEFAULT_HEADERS_CONFIG)
	Headers = []utils.HeadersConfig{}
	if headersEnv != "" {
		err := json.Unmarshal([]byte(headersEnv), &Headers)
		if err != nil {
			utils.HandleError(fmt.Errorf("failed to unmarshal HEADERS: %w", err))
		}
	}

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("current directory: %s\n", currentPath)
	fmt.Printf("endpoint: %s\nbucketName: %s\nfolder: %s\nincremental: %t\nexclude: %v\nindexPage: %s\nnotFoundPage: %s\nisCname: %t\nskipSetting: %t\n",
		Endpoint, BucketName, Folder, IsIncremental, Exclude, IndexPage, NotFoundPage, IsCname, SkipSetting)
	fmt.Printf("headers: %s\n", headersEnv)

	Client, err = oss.New(Endpoint, AccessKeyID, AccessKeySecret, oss.UseCname(IsCname))
	if err != nil {
		utils.HandleError(err)
	}

	Bucket, err = Client.Bucket(BucketName)
	if err != nil {
		utils.HandleError(err)
	}
}
