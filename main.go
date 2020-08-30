package main

import (
	"fmt"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
)

func main() {
	defer utils.TimeCost()()

	if !config.SkipSetting {
		operation.SetStaticWebsiteConfig(config.Client, config.Website)
	} else {
		fmt.Println("skip setting static pages related configuration")
	}

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
