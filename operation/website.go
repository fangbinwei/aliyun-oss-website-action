package operation

import (
	"fmt"

	"aliyun-oss-website-action/config"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

// SetStaticWebsiteConfig is used to set some option of website, like redirect strategy, index page, 404 page.
func SetStaticWebsiteConfig(client *oss.Client, o *config.WebsiteOptions) error {
	bEnable := true
	supportSubDirType := 0
	websiteDetailConfig, err := client.GetBucketWebsite(config.Bucket.BucketName)
	if err != nil {
		serviceError, ok := err.(oss.ServiceError)
		// 404 means NoSuchWebsiteConfiguration
		if !ok || serviceError.StatusCode != 404 {
			fmt.Println("Failed to get website detail configuration, skip setting", err)
			return err
		}
	}
	wxml := oss.WebsiteXML(websiteDetailConfig)
	wxml.IndexDocument.Suffix = o.IndexPage
	wxml.ErrorDocument.Key = o.NotFoundPage
	wxml.IndexDocument.SupportSubDir = &bEnable
	wxml.IndexDocument.Type = &supportSubDirType

	err = client.SetBucketWebsiteDetail(config.BucketName, wxml)
	if err != nil {
		fmt.Printf("Failed to set website detail configuration: %v\n", err)
		return err
	}
	return nil
}
