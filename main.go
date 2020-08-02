package main

import (
	"fmt"
	"os"

	"github.com/aliyun-oss-website-action/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {
	defer utils.TimeCost()()
	endpoint := os.Getenv("ENDPOINT")
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	folder := os.Getenv("FOLDER")
	bucketName := os.Getenv("BUCKET")

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("endpoint: %s\n bucketName: %s\n folder: %s\n current directory: %s\n", endpoint, bucketName, folder, currentPath)

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		utils.HandleError(err)
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
	utils.WalkDir(folder, records)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		utils.HandleError(err)
	}

	errs := utils.UploadFiles(folder, bucket, records)
	if errs != nil {
		fmt.Println("errors:")
		for i, err := range errs {
			fmt.Printf("%d\n%v\n", i, err)
		}
	}

}
