package operation

import (
	"fmt"
	"strconv"

	"aliyun-oss-website-action/config"

	// "github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// SetStaticWebsiteConfig is used to set some option of website, like redirect strategy, index page, 404 page.
func SetStaticWebsiteConfig() error {
	// bEnable := true
	// supportSubDirType := 0
	websiteDetailConfig, err := config.Client.GetBucketWebsite(config.Bucket.BucketName)
	if err != nil {
		serviceError, ok := err.(oss.ServiceError)
		// 404 means NoSuchWebsiteConfiguration
		if !ok || serviceError.StatusCode != 404 {
			fmt.Println("Failed to get website detail configuration, skip setting", err)
			return err
		}
	}
	wxml := oss.WebsiteXML(websiteDetailConfig)
	wxml.IndexDocument.Suffix = config.IndexPage
	wxml.ErrorDocument.Key = config.NotFoundPage
	// wxml.IndexDocument.SupportSubDir = &bEnable
	// wxml.IndexDocument.Type = &supportSubDirType
	error_http_code, _ := strconv.Atoi(config.ErrorDocumentHTTPCode)

	// Define one website detail
	ruleOk := oss.RoutingRule{
		RuleNumber: 1,
		Condition: oss.Condition{
			KeyPrefixEquals:             "abc",
			HTTPErrorCodeReturnedEquals: error_http_code,
		},
	}

	if len(wxml.RoutingRules) == 0 {
		wxml.RoutingRules = append(wxml.RoutingRules, ruleOk)
	} else {
		wxml.RoutingRules = []oss.RoutingRule{ruleOk}
	}

	err = config.Client.SetBucketWebsiteDetail(config.BucketName, wxml)
	if err != nil {
		fmt.Printf("Failed to set website detail configuration: %v\n", err)
		return err
	}
	return nil
}
