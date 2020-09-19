package main

import (
	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
	"fmt"
	"strings"
	"testing"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)

	err := operation.SetStaticWebsiteConfig()
	assert.NoError(err)

	errs := operation.DeleteObjects(config.Bucket)
	assert.Equal(len(errs), 0)

	records := utils.WalkDir(config.Folder)

	// overwrite, since dotenv doesn't support multiline
	config.Exclude = []string{"exclude.txt", "exclude/"}
	fmt.Println("---- upload start ----")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records)
	assert.Equal(len(uploadErrs), 0, uploadErrs)

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)

	assert.Equal(len(uploaded), len(lor.Objects))

	// test exclude
	lor, err = config.Bucket.ListObjects(oss.Prefix("exclude"))
	assert.NoError(err)
	assert.Empty(lor.Objects)

	// test cache-control
	prefix := config.Folder
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	for _, u := range uploaded {
		// 如果自定义域名解析到了cdn, 这个接口会报错, 但是上面的测试流程正常
		// 避开方法: env中endpoint使用bucket的endpoint或者bucket域名, 而不是自定义域名
		props, err := config.Bucket.GetObjectDetailedMeta(strings.TrimPrefix(u.PathOSS, prefix))
		assert.NoError(err)
		cacheControl := props.Get("Cache-Control")
		if utils.IsImage(u.PathOSS) {
			assert.Equal(cacheControl, config.ImageCacheControl)
		} else if utils.IsHTML(u.PathOSS) {
			assert.Equal(cacheControl, config.HTMLCacheControl)
		} else {
			assert.Equal(cacheControl, config.OtherCacheControl)
		}
	}

}
