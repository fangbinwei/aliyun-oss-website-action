package main

import (
	"fmt"

	"github.com/aliyun-oss-website-action/utils"
	"github.com/aliyun-oss-website-action/config"
	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

func main() {
	defer utils.TimeCost()()

	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		utils.HandleError(err)
	}

	bEnable := true
	supportSubDirType := 0
	wxml := oss.WebsiteXML{}
	wxml.IndexDocument.Suffix = config.IndexPage
	wxml.ErrorDocument.Key = config.NotFoundPage
	wxml.IndexDocument.SupportSubDir = &bEnable
	wxml.IndexDocument.Type = &supportSubDirType

	err = client.SetBucketWebsiteDetail(config.BucketName, wxml)
	if err != nil {
		fmt.Printf("Failed to set website detail: %v", err)
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		utils.HandleError(err)
	}

	fmt.Println("---- delete start ----")
	deleteErrs := utils.DeleteObjects(bucket)
	utils.LogErrors(deleteErrs)
	fmt.Println("---- delete end ----")

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- upload start ----")
	uploadErrs := utils.UploadObjects(config.Folder, bucket, records)
	utils.LogErrors(uploadErrs)
	fmt.Println("---- upload end ----")
}
