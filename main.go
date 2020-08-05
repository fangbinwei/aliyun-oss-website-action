package main

import (
	"fmt"

	"github.com/aliyun-oss-website-action/config"
	"github.com/aliyun-oss-website-action/operation"
	"github.com/aliyun-oss-website-action/utils"
)

func main() {
	defer utils.TimeCost()()

	operation.SetStaticWebsiteConfig(config.Client, config.Website)

	fmt.Println("---- delete start ----")
	deleteErrs := operation.DeleteObjects(config.Bucket)
	utils.LogErrors(deleteErrs)
	fmt.Println("---- delete end ----")

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- upload start ----")
	uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records)
	utils.LogErrors(uploadErrs)
	fmt.Println("---- upload end ----")
}
