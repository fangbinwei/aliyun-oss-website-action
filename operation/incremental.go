package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	// "github.com/fangbinwei/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const INCREMENTAL_CONFIG = ".actioninfo"

type IncrementalConfig struct {
	sync.RWMutex
	M map[string]struct {
		ContentMD5   string
		CacheControl string
	}
}

func (i *IncrementalConfig) stringify() ([]byte, error) {
	j, err := json.Marshal(i.M)
	return j, err
}

func (i *IncrementalConfig) parse(raw []byte) error {
	err := json.Unmarshal(raw, &(i.M))
	return err
}

func generateIncrementalConfig(uploaded []UploadedObject) ([]byte, error) {
	i := new(IncrementalConfig)
	i.M = make(map[string]struct {
		ContentMD5   string
		CacheControl string
	})
	for _, u := range uploaded {
		if !u.ValidHash {
			continue
		}
		i.M[u.ObjectKey] = struct {
			ContentMD5   string
			CacheControl string
		}{
			ContentMD5:   u.ContentMD5,
			CacheControl: u.CacheControl,
		}
	}
	j, err := i.stringify()
	return j, err

}

func UploadIncrementalConfig(bucket *oss.Bucket, records []UploadedObject) error {
	j, err := generateIncrementalConfig(records)
	if err != nil {
		fmt.Printf("Failed to generate incremental info: %v\n", err)
		return err
	}

	options := []oss.Option{
		oss.ObjectACL(oss.ACLPrivate),
	}
	err = bucket.PutObject(INCREMENTAL_CONFIG, bytes.NewReader(j), options...)
	if err != nil {
		fmt.Printf("Failed to upload incremental info: %v\n", err)
		return err
	}

	fmt.Printf("Update & Upload incremental info: %s\n", INCREMENTAL_CONFIG)
	return nil
}

func GetRemoteIncrementalConfig(bucket *oss.Bucket) (*IncrementalConfig, error) {
	c := new(bytes.Buffer)
	body, err := bucket.GetObject(INCREMENTAL_CONFIG)
	if err != nil {
		fmt.Printf("Failed to get remote incremental info: %v\n", err)
		return nil, err
	}
	io.Copy(c, body)
	body.Close()
	i := new(IncrementalConfig)
	err = i.parse(c.Bytes())
	if err != nil {
		fmt.Printf("Failed to parse remote incremental info: %v\n", err)
		return nil, err
	}
	fmt.Printf("Get remote incremental info: %s\n", INCREMENTAL_CONFIG)

	return i, nil
}
