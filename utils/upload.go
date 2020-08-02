package utils

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// UploadFiles upload files to OSS
func UploadFiles(root string, bucket *oss.Bucket, records chan FileInfoType) []error {
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	var sw sync.WaitGroup
	var errs []error
	for item := range records {
		sw.Add(1)
		var tokens = make(chan struct{}, 10)
		go func(item FileInfoType) {
			defer sw.Done()
			fPath := item.Path
			objectKey := strings.TrimPrefix(item.PathOSS, root)
			tokens <- struct{}{}
			err := bucket.PutObjectFromFile(objectKey, fPath)
			<-tokens
			if err != nil {
				errs = append(errs, fmt.Errorf("objectKey: %s\nfilePath: %s\nerror: %v", objectKey, fPath, err))
				return
			}
			fmt.Printf("objectKey: %s\nfilePath: %s\n", objectKey, fPath)
			fmt.Println()
		}(item)
	}
	sw.Wait()
	if len(errs) > 0 {
		return errs
	}
	return nil
}
