package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fangbinwei/aliyun-oss-go-sdk/oss"
)

const INCREMENTAL_CONFIG = ".actioninfo"

type IncrementalConfig map[string]struct {
	ContentMD5 string
}

func (i *IncrementalConfig) stringify() ([]byte, error) {
	j, err := json.Marshal(i)
	return j, err
}

func (i *IncrementalConfig) parse(raw []byte) error {
	err := json.Unmarshal(raw, i)
	return err
}

func generateIncrementalConfig(uploaded []UploadedObject) ([]byte, error) {
	config := make(IncrementalConfig)
	for _, u := range uploaded {
		config[u.ObjectKey] = struct{ ContentMD5 string }{
			ContentMD5: u.ContentMD5,
		}
	}
	j, err:= config.stringify()
	return j, err

}

func UploadIncrementalConfig(bucket *oss.Bucket, records []UploadedObject) {
	config, err := generateIncrementalConfig(records)
	if err != nil {
		fmt.Printf("Failed to generate incremental config: %v\n", err)
		return
	}

	options := []oss.Option{
		oss.ObjectACL(oss.ACLPrivate),
	}
	err = bucket.PutObject(INCREMENTAL_CONFIG, bytes.NewReader(config), options...)
	if err != nil {
		fmt.Printf("Failed to upload incremental config: %v\n", err)
		return
	}

	fmt.Printf("Upload incremental info: %s\n", INCREMENTAL_CONFIG)
}

func GetIncrementalConfig(bucket *oss.Bucket) (IncrementalConfig, error) {
	c := new(bytes.Buffer)
	body, err := bucket.GetObject(INCREMENTAL_CONFIG)
	if err != nil {
		fmt.Println("Failed to get incremental config")
		return nil, err
	}
	io.Copy(c, body)
	body.Close()
	var config IncrementalConfig
	err = config.parse(c.Bytes())
	if err != nil {
		fmt.Printf("Failed to parse incremental config: %v\n", err)
		return nil, err
	}
	fmt.Printf("Get incremental info: %s\n", INCREMENTAL_CONFIG)

	return config, nil
}
