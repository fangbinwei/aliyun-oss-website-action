package main

import (
	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	errs := operation.DeleteObjects(config.Bucket)
	assert.Equal(len(errs), 0)

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- upload start ----")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records)
	assert.Equal(len(uploadErrs), 0, uploadErrs)

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)

	assert.Equal(len(uploaded), len(lor.Objects))

	// test cache-control
	prefix := config.Folder
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	for _, u := range uploaded {
		props, err := config.Bucket.GetObjectDetailedMeta(strings.TrimPrefix(u.PathOSS, prefix))
		assert.NoError(err)
		cacheControl := props.Get("Cache-Control")
		if operation.IsImage(u.PathOSS) {
			assert.Equal(cacheControl, config.ImageCacheControl)
		} else if operation.IsHTML(u.PathOSS) {
			assert.Equal(cacheControl, config.HTMLCacheControl)
		} else {
			assert.Equal(cacheControl, config.OtherCacheControl)
		}
	}

}
