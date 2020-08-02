package main

import (
	"fmt"
	"os"

	"github.com/aliyun-oss-website-action/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {
	endpoint := os.Getenv("ENDPOINT")
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	folder := os.Getenv("FOLDER")
	bucketName:= os.Getenv("BUCKET")
	fmt.Println(endpoint, bucketName, folder)

	defer utils.TimeCost()()
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
