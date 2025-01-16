package operation

import (
	"fmt"
	"sync"

	"aliyun-oss-website-action/utils"

	// "github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const maxKeys = 100

// DeleteObjects is used to delete all objects of the bucket
func DeleteObjects(bucket *oss.Bucket) []error {
	var errs []error
	objKeyCollection := make(chan string, maxKeys)
	go listObjects(bucket, objKeyCollection)

	var sw sync.WaitGroup
	var mutex sync.Mutex
	tokens := make(chan struct{}, 10)
	for k := range objKeyCollection {
		sw.Add(1)
		go func(key string) {
			defer sw.Done()
			defer func() {
				<-tokens
			}()
			tokens <- struct{}{}
			err := deleteObject(bucket, key)
			if err != nil {
				mutex.Lock()
				errs = append(errs, fmt.Errorf("[FAILED] objectKey: %s\nDetail: %v", key, err))
				mutex.Unlock()
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

func DeleteObjectsIncremental(bucket *oss.Bucket, i *IncrementalConfig) []error {
	if i == nil {
		return nil
	}
	// delete incremental info
	i.M[INCREMENTAL_CONFIG] = struct {
		ContentMD5   string
		CacheControl string
	}{}

	// TODO: optimize
	var errs []error

	var sw sync.WaitGroup
	var mutex sync.Mutex
	tokens := make(chan struct{}, 10)
	for k := range i.M {
		sw.Add(1)
		go func(key string) {
			defer sw.Done()
			defer func() {
				<-tokens
			}()
			tokens <- struct{}{}
			err := deleteObject(bucket, key)
			if err != nil {
				mutex.Lock()
				errs = append(errs, fmt.Errorf("[FAILED] objectKey: %s\nDetail: %v", key, err))
				mutex.Unlock()
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

func deleteObject(bucket *oss.Bucket, key string) error {
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
			utils.HandleError(err)
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
