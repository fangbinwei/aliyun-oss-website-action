package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliyun-oss-website-action/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {
	client, err := oss.New("", "", "")
	if err != nil {
		HandleError(err)
	}

	// bEnable := true
	// wxml := oss.WebsiteXML{}
	// wxml.IndexDocument.Suffix = "index.html"
	// wxml.IndexDocument.SupportSubDir = &bEnable
	// wxml.IndexDocument.Type = "0"
	// err = client.SetBucketWebsiteDetail("fangbinwei-blog", wxml)
	// if err != nil {
	// 	HandleError(err)
	// }

	records := make(chan utils.FileInfoType)
	utils.WalkDir("dist", records)

	bucket, err := client.Bucket("fangbinwei-blog")
	if err != nil {
		HandleError(err)
	}
	for item := range records {
		fPath := item.Path
		objectKey := strings.TrimPrefix(item.PathOSS, "dist/")
		err = bucket.PutObjectFromFile(objectKey, fPath)
		if err != nil {
			fmt.Println("occurred error:", err)
		}
		fmt.Printf("objectKey: %s\n filePath: %s\n", objectKey, fPath)
	}
}

// HandleError is error handling method, print error and exit
func HandleError(err error) {
	fmt.Println("occurred error:", err)
	os.Exit(1)
}
