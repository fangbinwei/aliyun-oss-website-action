package operation

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/utils"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

type UploadedObject struct {
	ObjectKey   string
	Incremental bool
	utils.FileInfoType
}
type UploadError struct {
	detail    utils.FileInfoType
	objectKey string
	fPath     string
	msg       string
	rawError  error
}
type UploadOptions struct {
	I           *IncrementalConfig
	Concurrency int
}

func (e *UploadError) Error() string {
	return fmt.Sprintf("[FAILED] objectKey: %s\nfilePath: %s\nDetail: %v", e.objectKey, e.fPath, e.msg)
}

// UploadObjects upload files to OSS
func UploadObjects(root string, bucket *oss.Bucket, records <-chan utils.FileInfoType, uploadOptions UploadOptions) ([]UploadedObject, []error) {
	root = path.Clean(root) + "/"
	concurrency := uploadOptions.Concurrency
	if concurrency <= 0 {
		concurrency = 30
	}

	var sw sync.WaitGroup
	var errorMutex sync.Mutex
	var uploadedMutex sync.Mutex
	var errs []error
	uploaded := make([]UploadedObject, 0, 50)
	var tokens = make(chan struct{}, concurrency)
	for item := range records {
		sw.Add(1)
		go func(item utils.FileInfoType) {
			defer sw.Done()
			fPath := item.Path
			objectKey := strings.TrimPrefix(item.PathOSS, root)
			options := getHTTPHeader(&item)

			if shouldExclude(objectKey) {
				// fmt.Printf("[EXCLUDE] objectKey: %s\n\n", objectKey)
				return
			}
			if shouldSkip(item, objectKey, uploadOptions.I) {
				// fmt.Printf("[SKIP] objectKey: %s \n\n", objectKey)
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
				uErr := &UploadError{detail: item, fPath: fPath, objectKey: objectKey, rawError: err, msg: err.Error()}
				errs = append(errs, uErr)
				errorMutex.Unlock()
				return
			}
			// fmt.Printf("objectKey: %s\nfilePath: %s\n\n", objectKey, fPath)
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

func UploadRetry(errs []error, times int) ([]UploadedObject, []error) {
	if len(errs) == 0 {
		return []UploadedObject{}, nil
	}
	uploadedResult := make([]UploadedObject, 0, 50)

	retry := func(e []error) ([]UploadedObject, []error) {
		time.Sleep(time.Second * 3)
		records := make(chan utils.FileInfoType, 20)
		go func() {
			defer close(records)
			for _, item := range e {
				if uploadError, ok := item.(*UploadError); ok {
					records <- uploadError.detail
				}
			}
		}()
		return UploadObjects(config.Folder, config.Bucket, records, UploadOptions{Concurrency: 20})
	}
	for i := 0; i < times; i++ {
		if len(errs) == 0 {
			return uploadedResult, nil
		}
		uploaded, uploadError := retry(errs)
		fmt.Printf("\n[RETRY %v]", i+1)
		LogUploadedResult(uploaded, uploadError)
		uploadedResult = append(uploadedResult, uploaded...)
		errs = uploadError
	}
	return uploadedResult, errs
}

func LogUploadedResult(result []UploadedObject, errs []error) {
	if result == nil {
		return
	}
	uploadedCount := 0
	skippedCount := 0

	for _, v := range result {
		if v.Incremental {
			skippedCount++
		} else {
			uploadedCount++
		}
	}
	fmt.Printf("\nuploaded %v object(s), skipped %v object(s), %v error(s).\n", uploadedCount, skippedCount, len(errs))
}

func getHTTPHeader(item *utils.FileInfoType) []oss.Option {
	return []oss.Option{
		getCacheControlOption(item),
	}
}

func getCacheControlOption(item *utils.FileInfoType) oss.Option {
	var value string
	filename := item.Name

	if utils.IsHTML(filename) {
		value = config.HTMLCacheControl
	} else if utils.IsImage(filename) {
		// pic name may not contains hash, so use different strategy
		// 10 days
		value = config.ImageCacheControl
	} else if utils.IsPDF(filename) {
		value = config.PDFCacheControl
	} else {
		// static assets like .js .css, use contentHash in file name, so html can update these files.
		// 30 days
		value = config.OtherCacheControl
	}
	item.CacheControl = value
	return oss.CacheControl(value)
}

func shouldExclude(objectKey string) bool {
	if utils.Match(config.Exclude, objectKey) {
		return true
	}
	return false
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
	if item.ValidHash && item.ContentMD5 == remoteConfig.ContentMD5 && item.CacheControl == remoteConfig.CacheControl {
		return true
	}
	return false
}
