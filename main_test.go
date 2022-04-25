package main

import (
	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
	"fmt"
	"testing"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)

	err := operation.SetStaticWebsiteConfig()
	assert.NoError(err)

	fmt.Println("---- [incremental] ---->")
	incremental, err := operation.GetIncrementalConfig(config.Bucket)
	// assert.NoError(err)
	fmt.Println("<---- [incremental] ----")

	if incremental == nil {
		fmt.Println("---- [delete] ---->")
		errs := operation.DeleteObjects(config.Bucket)
		fmt.Println("<---- [delete] ----")
		assert.Equal(len(errs), 0)
	}

	records := utils.WalkDir(config.Folder)

	// overwrite, since dotenv doesn't support multiline
	config.Exclude = []string{"exclude.txt", "exclude/"}
	fmt.Println("---- [upload]  ---->")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records, incremental)
	fmt.Println("<---- [upload]  ----")
	assert.Equal(len(uploadErrs), 0, uploadErrs)

	// test incremental
	if incremental != nil {
		fmt.Println("---- [delete] ---->")
		errs := operation.DeleteObjectsIncremental(config.Bucket, incremental)
		fmt.Println("<---- [delete] ----")
		assert.Equal(len(errs), 0)
	}

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(len(uploaded), len(lor.Objects))

	fmt.Println("---- [incremental] ---->")
	operation.UploadIncrementalConfig(config.Bucket, uploaded)
	fmt.Println("<---- [incremental] ----")

	// test exclude
	lor, err = config.Bucket.ListObjects(oss.Prefix("exclude"))
	assert.NoError(err)
	assert.Empty(lor.Objects)

	// test cache-control
	for _, u := range uploaded {
		// 如果自定义域名解析到了cdn, 这个接口会报错, 但是上面的测试流程正常
		// 避开方法: env中endpoint使用bucket的endpoint或者bucket域名, 而不是自定义域名
		props, err := config.Bucket.GetObjectDetailedMeta(u.ObjectKey)
		assert.NoError(err)
		cacheControl := props.Get("Cache-Control")
		if utils.IsImage(u.ObjectKey) {
			assert.Equal(cacheControl, config.ImageCacheControl)
		} else if utils.IsHTML(u.ObjectKey) {
			assert.Equal(cacheControl, config.HTMLCacheControl)
		} else {
			assert.Equal(cacheControl, config.OtherCacheControl)
		}
	}

}
