package main

import (
	"fmt"
	"os"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
)

func main() {
	code := 0
	defer func() {
		os.Exit(code)
	}()
	defer utils.TimeCost()()

	if !config.SkipSetting {
		operation.SetStaticWebsiteConfig()
	} else {
		fmt.Println("skip setting static pages related configuration")
	}

	var incremental operation.IncrementalConfig
	if config.IsIncremental {
		fmt.Println("---- [incremental] ---->")
		incremental, _ = operation.GetIncrementalConfig(config.Bucket)
		fmt.Println("<---- [incremental] ----")
	}
	if !config.IsIncremental || incremental == nil {
		// TODO: delete after upload
		fmt.Println("---- [delete] ---->")
		deleteErrs := operation.DeleteObjects(config.Bucket)
		utils.LogErrors(deleteErrs)
		fmt.Println("<---- [delete] ----")
	}

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- [upload] ---->")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records, incremental)
	utils.LogErrors(uploadErrs)
	fmt.Println("<---- [upload] ----")

	if config.IsIncremental && incremental != nil {
		fmt.Println("---- [delete] ---->")
		deleteErrs := operation.DeleteObjectsIncremental(config.Bucket, incremental)
		utils.LogErrors(deleteErrs)
		fmt.Println("<---- [delete] ----")
	}

	if config.IsIncremental {
		fmt.Println("---- [incremental] ---->")
		operation.UploadIncrementalConfig(config.Bucket, uploaded)
		fmt.Println("<---- [incremental] ----")
	}

	if len(uploadErrs) > 0 {
		code = 1
	}

}
