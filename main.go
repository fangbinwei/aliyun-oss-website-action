package main

import (
	"fmt"
	"os"

	"github.com/aliyun-oss-website-action/utils"
	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
)

func main() {
	defer utils.TimeCost()()
	godotenv.Load(".env.local")

	endpoint := os.Getenv("ENDPOINT")
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	folder := os.Getenv("FOLDER")
	bucketName := os.Getenv("BUCKET")

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("endpoint: %s\nbucketName: %s\nfolder: %s\ncurrent directory: %s\n", endpoint, bucketName, folder, currentPath)

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		utils.HandleError(err)
	}

	bEnable := true
	supportSubDirType := 0
	wxml := oss.WebsiteXML{}
	wxml.IndexDocument.Suffix = "index.html"
	wxml.ErrorDocument.Key = "404.html"
	wxml.IndexDocument.SupportSubDir = &bEnable
	wxml.IndexDocument.Type = &supportSubDirType

	err = client.SetBucketWebsiteDetail(bucketName, wxml)
	if err != nil {
		fmt.Printf("Failed to set website detail: %v", err)
	}

	records := make(chan utils.FileInfoType)
	utils.WalkDir(folder, records)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		utils.HandleError(err)
	}

	fmt.Println("----upload info-----")

	errs := utils.UploadFiles(folder, bucket, records)
	if errs != nil {
		fmt.Println("errors:")
		for i, err := range errs {
			fmt.Printf("%d\n%v\n", i, err)
		}
	}

}
