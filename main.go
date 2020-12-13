package main

import (
	"fmt"
	"os"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
)

func main() {
	defer utils.TimeCost()()

	if !config.SkipSetting {
		operation.SetStaticWebsiteConfig()
	} else {
		fmt.Println("skip setting static pages related configuration")
	}

	fmt.Println("---- delete start ----")
	deleteErrs := operation.DeleteObjects(config.Bucket)
	utils.LogErrors(deleteErrs)
	fmt.Println("---- delete end ----")

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- upload start ----")
	_, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records)
	utils.LogErrors(uploadErrs)
	fmt.Println("---- upload end ----")

	if len(uploadErrs) > 0 {
		os.Exit(1)
	}

}
