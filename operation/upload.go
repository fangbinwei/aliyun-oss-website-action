package operation

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/utils"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

type UploadedObject struct {
	ObjectKey   string
	Incremental bool
	utils.FileInfoType
}

// UploadObjects upload files to OSS
func UploadObjects(root string, bucket *oss.Bucket, records <-chan utils.FileInfoType, i *IncrementalConfig) ([]UploadedObject, []error) {
	if root == "/" {
		fmt.Println("You should not upload the root directory, use ./ instead. 通常来说, 你不应该上传根目录, 也许你是要配置 ./")
		os.Exit(1)
	} else {
		root = path.Clean(root) + "/"
	}

	var sw sync.WaitGroup
	var errorMutex sync.Mutex
	var uploadedMutex sync.Mutex
	var errs []error
	uploaded := make([]UploadedObject, 0, 20)
	var tokens = make(chan struct{}, 30)
	for item := range records {
		sw.Add(1)
		go func(item utils.FileInfoType) {
			defer sw.Done()
			fPath := item.Path
			objectKey := strings.TrimPrefix(item.PathOSS, root)
			options := getHTTPHeader("/"+objectKey, &item)

			if shouldExclude(objectKey) {
				fmt.Printf("[EXCLUDE] objectKey: %s\n\n", objectKey)
				return
			}
			if shouldSkip(item, objectKey, i) {
				fmt.Printf("[SKIP] objectKey: %s \n\n", objectKey)
				uploadedMutex.Lock()
				uploaded = append(uploaded, UploadedObject{ObjectKey: objectKey, Incremental: true, FileInfoType: item})
				uploadedMutex.Unlock()
				return
			}

			tokens <- struct{}{}
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
			uploaded = append(uploaded, UploadedObject{ObjectKey: objectKey, FileInfoType: item})
			uploadedMutex.Unlock()
		}(item)
	}
	sw.Wait()
	if len(errs) > 0 {
		return uploaded, errs
	}
	return uploaded, nil
}

func getHTTPHeader(path string, item *utils.FileInfoType) (option []oss.Option) {
	headersConfig := utils.MatchHeadersConfig(path, config.Headers)
	for k, v := range headersConfig {
		option = append(option, oss.SetHeader(k, v))
	}
	headersMD5, err := utils.GetHeadersConfigMD5(headersConfig)
	if err != nil {
		fmt.Printf("Failed to get headers config MD5: %v\n", err)
	}
	item.HeadersMD5 = headersMD5
	return
}

func shouldExclude(objectKey string) bool {
	return utils.Match(config.Exclude, objectKey)
}

func shouldSkip(item utils.FileInfoType, objectKey string, i *IncrementalConfig) bool {
	if i == nil {
		return false
	}
	i.RLock()
	remoteConfig, ok := i.M[objectKey]
	i.RUnlock()
	if !ok {
		return false
	}
	// delete existed objectKey in incremental map, the left is what we should delete
	i.Lock()
	delete(i.M, objectKey)
	i.Unlock()
	if item.ValidHash && item.ContentMD5 == remoteConfig.ContentMD5 && item.HeadersMD5 == remoteConfig.HeadersMD5 {
		return true
	}
	return false
}
