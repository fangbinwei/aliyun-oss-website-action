package utils

import (
	"fmt"
	"sync"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

const maxKeys = 100

// DeleteObjects is used to delete all objects of the bucket
func DeleteObjects(bucket *oss.Bucket) []error {
	var errs []error
	objKeyCollection := make(chan string, maxKeys)
	go listObjects(bucket, objKeyCollection)

	var sw sync.WaitGroup
	tokens := make(chan struct{}, 10)
	for k := range objKeyCollection {
		sw.Add(1)
		go func(key string) {
			defer sw.Done()
			defer func() {
				<-tokens
			}()
			tokens <- struct{}{}
			err := delete(bucket, key)
			if err != nil {
				errs = append(errs, fmt.Errorf("[FAILED] objectKey: %s\nDetail: %v", key, err))
				return
			}
			fmt.Printf("objectKey: %s\n", key)
		}(k)
	}
	sw.Wait()

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func delete(bucket *oss.Bucket, key string) error {
	err := bucket.DeleteObject(key)
	if err != nil {
		return err
	}
	return nil
}

func listObjects(bucket *oss.Bucket, objKeyCollection chan<- string) {
	marker := oss.Marker("")
	for {
		lor, err := bucket.ListObjects(oss.MaxKeys(maxKeys), marker)
		if err != nil {
			HandleError(err)
		}
		for _, object := range lor.Objects {
			objKeyCollection <- object.Key
		}
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			close(objKeyCollection)
			break
		}
	}

}
