package main

import (
	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
	"fmt"
	"os"
	"testing"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/assert"
)

var headersForTest = []utils.HeadersConfig{
	{
		Path: "\\.html",
		Headers: map[string]string{
			"Cache-Control": "public, max-age=0, must-revalidate",
		},
	},
	{
		Path: "\\.(jpg|png|gif|webp|ico)$",
		Headers: map[string]string{
			"Cache-Control": "max-age=31536000",
		},
	},
	{
		Path: "\\.pdf$",
		Headers: map[string]string{
			"Cache-Control": "max-age=11536000",
		},
	},
	{
		Path: "",
		Headers: map[string]string{
			"Cache-Control": "no-cache",
		},
	},
}

func testSetStaticWebsiteConfig(t *testing.T) {
	assert := assert.New(t)

	err := operation.SetStaticWebsiteConfig()
	assert.NoError(err)
}

func testUpload(t *testing.T) {
	config.Headers = headersForTest
	assert := assert.New(t)
	fmt.Println("---- [delete] ---->")
	errs := operation.DeleteObjects(config.Bucket)
	fmt.Println("<---- [delete end] ----")
	assert.Equal(len(errs), 0)

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(0, len(lor.Objects))

	records := utils.WalkDir(config.Folder)

	// overwrite, since dotenv doesn't support multiline
	config.Exclude = []string{"exclude.txt", "exclude/"}
	fmt.Println("---- [upload]  ---->")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records, nil)
	fmt.Println("<---- [upload end]  ----")
	assert.Equal(0, len(uploadErrs), uploadErrs)

	lor, err = config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(len(uploaded), len(lor.Objects))

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
		headers := utils.MatchHeadersConfig(u.ObjectKey, config.Headers)
		assert.Equal(headers["Cache-Control"], props.Get("Cache-Control"))
	}
}

func testUploadIncrementalFirst(t *testing.T) {
	config.Headers = headersForTest
	assert := assert.New(t)

	fmt.Println("---- [incremental] ---->")
	incremental, err := operation.GetRemoteIncrementalConfig(config.Bucket)
	assert.Error(err)
	assert.Empty(incremental)
	fmt.Println("<---- [incremental end] ----")

	if incremental == nil {
		fmt.Println("---- [delete] ---->")
		errs := operation.DeleteObjects(config.Bucket)
		fmt.Println("<---- [delete end] ----")
		assert.Equal(len(errs), 0)
	}

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(0, len(lor.Objects))

	records := utils.WalkDir(config.Folder)

	// overwrite, since dotenv doesn't support multiline
	config.Exclude = []string{"exclude.txt", "exclude/"}
	fmt.Println("---- [upload]  ---->")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records, incremental)
	fmt.Println("<---- [upload end]  ----")
	assert.Equal(0, len(uploadErrs), uploadErrs)

	lor, err = config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(len(uploaded), len(lor.Objects))

	fmt.Println("---- [incremental] ---->")
	err = operation.UploadIncrementalConfig(config.Bucket, uploaded)
	fmt.Println("<---- [incremental end] ----")
	assert.NoError(err)

	incremental, err = operation.GetRemoteIncrementalConfig(config.Bucket)
	assert.NoError(err)
	assert.Equal(len(uploaded), len(incremental.M))
	for _, v := range uploaded {
		assert.False(v.Incremental)
		assert.Equal(v.ContentMD5, incremental.M[v.ObjectKey].ContentMD5)
		assert.Equal(v.HeadersMD5, incremental.M[v.ObjectKey].HeadersMD5)
	}

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
		headers := utils.MatchHeadersConfig(u.ObjectKey, config.Headers)
		assert.Equal(headers["Cache-Control"], props.Get("Cache-Control"))
	}
}

func testUploadIncrementalSecond(t *testing.T) {
	assert := assert.New(t)

	fmt.Println("---- [incremental] ---->")
	incremental, err := operation.GetRemoteIncrementalConfig(config.Bucket)
	assert.NoError(err)
	fmt.Println("<---- [incremental end] ----")

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(len(lor.Objects)-1, len(incremental.M))

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- [upload]  ---->")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records, incremental)
	fmt.Println("<---- [upload end]  ----")
	assert.Equal(0, len(uploadErrs), uploadErrs)

	lor, err = config.Bucket.ListObjects()

	assert.NoError(err)
	assert.Equal(len(uploaded), len(lor.Objects)-1)

	// incremental.M中剩余的项是待删除的, 数量为0, 因为此次上传和上次上传的文件一模一样
	assert.Equal(0, len(incremental.M))

	fmt.Println("---- [delete] ---->")
	// 只删除.actioninfo
	errs := operation.DeleteObjectsIncremental(config.Bucket, incremental)
	fmt.Println("<---- [delete end] ----")
	assert.Equal(0, len(errs))

	lor, err = config.Bucket.ListObjects()

	assert.NoError(err)
	assert.Equal(len(uploaded), len(lor.Objects))

	fmt.Println("---- [incremental] ---->")
	err = operation.UploadIncrementalConfig(config.Bucket, uploaded)
	fmt.Println("<---- [incremental end] ----")
	assert.NoError(err)

	incremental, err = operation.GetRemoteIncrementalConfig(config.Bucket)
	assert.NoError(err)
	assert.Equal(len(uploaded), len(incremental.M))
	for _, v := range uploaded {
		// 全都不需要上传, 命中incremental
		assert.True(v.Incremental)
		assert.Equal(v.ContentMD5, incremental.M[v.ObjectKey].ContentMD5)
		assert.Equal(v.HeadersMD5, incremental.M[v.ObjectKey].HeadersMD5)
	}

}

func testUploadIncrementalThird(t *testing.T) {
	assert := assert.New(t)
	folder := "testdata/group2"
	// 改变cache-control会让对应文件重新上传, 即使hash没变
	// 我们可以勇敢在config.Headers开头添加一个针对ico文件的cache-control来实现上述操作
	config.Headers = headersForTest
	config.Headers = append([]utils.HeadersConfig{
		{
			Path: "\\.ico$",
			Headers: map[string]string{
				"Cache-Control": "max-age=114514191",
			},
		},
	}, config.Headers...)

	fmt.Println("---- [incremental] ---->")
	incremental, err := operation.GetRemoteIncrementalConfig(config.Bucket)
	assert.NoError(err)
	fmt.Println("<---- [incremental end] ----")

	lor, err := config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(len(lor.Objects)-1, len(incremental.M))

	records := utils.WalkDir(folder)

	fmt.Println("---- [upload]  ---->")
	uploaded, uploadErrs := operation.UploadObjects(folder, config.Bucket, records, incremental)
	fmt.Println("<---- [upload end]  ----")
	assert.Equal(0, len(uploadErrs), uploadErrs)

	// incremental.M中剩余的项是待删除的, 大于0
	assert.Greater(len(incremental.M), 0)
	for _, v := range uploaded {
		if v.ObjectKey == "empty.js" {
			assert.True(v.Incremental)
			continue
		}
		if v.ObjectKey == "example.js" {
			assert.False(v.Incremental)
			continue
		}
		if v.ObjectKey == "favicon.ico" {
			assert.False(v.Incremental)
			continue
		}
		t.Fail()
	}

	fmt.Println("---- [delete] ---->")
	errs := operation.DeleteObjectsIncremental(config.Bucket, incremental)
	fmt.Println("<---- [delete end] ----")
	assert.Equal(0, len(errs))

	lor, err = config.Bucket.ListObjects()
	assert.NoError(err)
	assert.Equal(len(uploaded), len(lor.Objects))

	fmt.Println("---- [incremental] ---->")
	err = operation.UploadIncrementalConfig(config.Bucket, uploaded)
	fmt.Println("<---- [incremental end] ----")
	assert.NoError(err)

	incremental, err = operation.GetRemoteIncrementalConfig(config.Bucket)
	assert.NoError(err)
	assert.Equal(len(uploaded), len(incremental.M))
	for _, v := range uploaded {
		assert.Equal(v.ContentMD5, incremental.M[v.ObjectKey].ContentMD5)
		assert.Equal(v.HeadersMD5, incremental.M[v.ObjectKey].HeadersMD5)
	}

}

func TestAction(t *testing.T) {
	t.Run("SetStaticWebsiteConfig", testSetStaticWebsiteConfig)
	t.Run("First upload", testUpload)
	t.Run("Second upload", testUpload)
	t.Run("First incremental upload without .actioninfo", testUploadIncrementalFirst)
	t.Run("Second incremental upload", testUploadIncrementalSecond)
	t.Run("Third incremental upload, change cache-control", testUploadIncrementalThird)
}

func TestMain(m *testing.M) {
	code := m.Run()

	fmt.Println("Empty bucket after test")
	operation.DeleteObjects(config.Bucket)

	os.Exit(code)
}
