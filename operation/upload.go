package operation

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/utils"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

type UploadedObject struct {
	ObjectKey   string
	ContentMD5  string
	incremental bool
}

// UploadObjects upload files to OSS
func UploadObjects(root string, bucket *oss.Bucket, records <-chan utils.FileInfoType, m IncrementalConfig) ([]UploadedObject, []error) {
	root = path.Clean(root) + "/"

	var sw sync.WaitGroup
	var errorMutex sync.Mutex
	var uploadedMutex sync.Mutex
	var errs []error
	uploaded := make([]UploadedObject, 0, 20)
	for item := range records {
		sw.Add(1)
		var tokens = make(chan struct{}, 10)
		go func(item utils.FileInfoType) {
			defer sw.Done()
			fPath := item.Path
			objectKey := strings.TrimPrefix(item.PathOSS, root)
			if shouldExclude(objectKey) {
				fmt.Printf("[EXCLUDE] objectKey: %s\n\n", objectKey)
				return
			}
			if m != nil {
				// delete existed objectKey in incremental map, the left is what we should delete
				defer delete(m, objectKey)
			}
			if shouldSkip(item, objectKey, m) {
				fmt.Printf("[SKIP] objectKey: %s \n\n", objectKey)
				uploaded = append(uploaded, UploadedObject{ObjectKey: objectKey, ContentMD5: item.ContentMD5, incremental: true})
				return
			}

			tokens <- struct{}{}
			options := getHTTPHeader(&item)
			err := bucket.PutObjectFromFile(objectKey, fPath, options...)
			<-tokens
			if err != nil {
				errorMutex.Lock()
				errs = append(errs, fmt.Errorf("[FAILED] objectKey: %s\nfilePath: %s\nDetail: %v", objectKey, fPath, err))
				errorMutex.Unlock()
				return
			}
			fmt.Printf("objectKey: %s\nfilePath: %s\n\n", objectKey, fPath)
			uploadedMutex.Lock()
			uploaded = append(uploaded, UploadedObject{ObjectKey: objectKey, ContentMD5: item.ContentMD5})
			uploadedMutex.Unlock()
		}(item)
	}
	sw.Wait()
	if len(errs) > 0 {
		return uploaded, errs
	}
	return uploaded, nil
}

func getHTTPHeader(item *utils.FileInfoType) []oss.Option {
	return []oss.Option{
		getCacheControlOption(item.Info.Name()),
	}
}

func getCacheControlOption(filename string) oss.Option {
	var value string
	if utils.IsHTML(filename) {
		value = config.HTMLCacheControl
	} else if utils.IsImage(filename) {
		// pic name may not contains hash, so use different strategy
		// 10 days
		value = config.ImageCacheControl
	} else {
		// static assets like .js .css, use contentHash in file name, so html can update these files.
		// 30 days
		value = config.OtherCacheControl
	}
	return oss.CacheControl(value)
}

func shouldExclude(objectKey string) bool {
	if utils.Match(config.Exclude, objectKey) {
		return true
	}
	return false
}

func shouldSkip(item utils.FileInfoType, objectKey string, m IncrementalConfig) bool {
	if m == nil {
		return false
	}
	val, ok := m[objectKey]
	if !ok {
		return false
	}
	if item.ContentMD5 == val.ContentMD5 {
		return true
	}
	return false
}
