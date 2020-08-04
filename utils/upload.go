package utils

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

// UploadObjects upload files to OSS
func UploadObjects(root string, bucket *oss.Bucket, records chan FileInfoType) []error {
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
			options := getHTTPHeader(&item)
			err := bucket.PutObjectFromFile(objectKey, fPath, options...)
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

func getHTTPHeader(item *FileInfoType) []oss.Option {
	return []oss.Option{
		getCacheControlOption(item.Info.Name()),
	}
}

func getCacheControlOption(filename string) oss.Option {
	var value string
	if isHTML(filename) {
		value = "no-cache"
	} else if isImage(filename) {
		// pic name may not contains hash, so use different strategy
		// 10 days
		value = "max-age=864000"
	} else {
		// static assets like .js .css, use contentHash in file name, so html can update these files.
		// 30 days
		value = "max-age=2592000"
	}
	return oss.CacheControl(value)
}

func isHTML(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".html")
}

func isImage(filename string) bool {
	imageExts := []string{
		".png",
		".jpg",
		".jpeg",
		".webp",
		".gif",
		".bmp",
		".tiff",
		".ico",
		".svg",
	}
	return func() bool {
		ext := path.Ext(filename)
		for _, e := range imageExts {
			if e == ext {
				return true
			}
		}
		return false
	}()
}
