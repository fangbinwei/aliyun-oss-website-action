package config

import (
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
	Prefix          string
	Exclude         []string
	BucketName      string
	IsCname         bool
	Client          *oss.Client
	Bucket          *oss.Bucket
	SkipSetting     bool
	IsIncremental   bool

	IndexPage         string
	NotFoundPage      string
	HTMLCacheControl  string
	ImageCacheControl string
	OtherCacheControl string
	PDFCacheControl   string
)

func init() {
	godotenv.Load(".env")
	godotenv.Load(".env.local")

	Endpoint = os.Getenv("ENDPOINT")
	IsCname = os.Getenv("CNAME") == "true"
	AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	Folder = os.Getenv("FOLDER")
	Prefix = os.Getenv("PREFIX")
	Exclude = utils.GetActionInputAsSlice(os.Getenv("EXCLUDE"))
	BucketName = os.Getenv("BUCKET")
	SkipSetting = os.Getenv("SKIP_SETTING") == "true"
	IsIncremental = os.Getenv("INCREMENTAL") == "true"

	IndexPage = utils.Getenv("INDEX_PAGE", "index.html")
	NotFoundPage = utils.Getenv("NOT_FOUND_PAGE", "404.html")
	HTMLCacheControl = utils.Getenv("HTML_CACHE_CONTROL", "no-cache")
	ImageCacheControl = utils.Getenv("IMAGE_CACHE_CONTROL", "max-age=864000")
	OtherCacheControl = utils.Getenv("OTHER_CACHE_CONTROL", "max-age=2592000")
	PDFCacheControl = utils.Getenv("PDF_CACHE_CONTROL", "max-age=2592000")

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("current directory: %s\n", currentPath)
	fmt.Printf("push to directory: %s\n", Prefix)
	fmt.Printf("endpoint: %s\nbucketName: %s\nfolder: %s\nincremental: %t\nexclude: %v\nindexPage: %s\nnotFoundPage: %s\nisCname: %t\nskipSetting: %t\n",
		Endpoint, BucketName, Folder, IsIncremental, Exclude, IndexPage, NotFoundPage, IsCname, SkipSetting)
	fmt.Printf("HTMLCacheControl: %s\nimageCacheControl: %s\notherCacheControl: %s\npdfCacheControl: %s\n",
		HTMLCacheControl, ImageCacheControl, OtherCacheControl, PDFCacheControl)

	Client, err = oss.New(Endpoint, AccessKeyID, AccessKeySecret, oss.UseCname(IsCname))
	if err != nil {
		utils.HandleError(err)
	}

	Bucket, err = Client.Bucket(BucketName)
	if err != nil {
		utils.HandleError(err)
	}
}
