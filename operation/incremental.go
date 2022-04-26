package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

const INCREMENTAL_CONFIG = ".actioninfo"

type IncrementalConfig struct {
	sync.RWMutex
	m map[string]struct {
		ContentMD5 string
	}
}

func (i *IncrementalConfig) stringify() ([]byte, error) {
	j, err := json.Marshal(i.m)
	return j, err
}

func (i *IncrementalConfig) parse(raw []byte) error {
	err := json.Unmarshal(raw, &(i.m))
	return err
}

func generateIncrementalConfig(uploaded []UploadedObject) ([]byte, error) {
	i := new(IncrementalConfig)
	i.m = make(map[string]struct{ ContentMD5 string })
	for _, u := range uploaded {
		i.m[u.ObjectKey] = struct{ ContentMD5 string }{
			ContentMD5: u.ContentMD5,
		}
	}
	j, err := i.stringify()
	return j, err

}

func UploadIncrementalConfig(bucket *oss.Bucket, records []UploadedObject) {
	j, err := generateIncrementalConfig(records)
	if err != nil {
		fmt.Printf("Failed to generate incremental info: %v\n", err)
		return
	}

	options := []oss.Option{
		oss.ObjectACL(oss.ACLPrivate),
	}
	err = bucket.PutObject(INCREMENTAL_CONFIG, bytes.NewReader(j), options...)
	if err != nil {
		fmt.Printf("Failed to upload incremental info: %v\n", err)
		return
	}

	fmt.Printf("Upload incremental info: %s\n", INCREMENTAL_CONFIG)
}

func GetRemoteIncrementalConfig(bucket *oss.Bucket) (*IncrementalConfig, error) {
	c := new(bytes.Buffer)
	body, err := bucket.GetObject(INCREMENTAL_CONFIG)
	if err != nil {
		fmt.Println("Failed to get remote incremental info")
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
