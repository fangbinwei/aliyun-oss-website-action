package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var sema = make(chan struct{}, 20)

// FileInfoType is a type which contains dir and os.FileInfo
type FileInfoType struct {
	Dir          string
	Path         string
	PathOSS      string
	Name         string
	CacheControl string // Complete 'CacheControl' when uploading files
	ContentMD5   string
	ValidHash    bool // if ContentMD5 is valid
}

// WalkDir get sub files of target dir
func WalkDir(root string) <-chan FileInfoType {
	fileInfos := make(chan FileInfoType, 100)
	var sw sync.WaitGroup
	sw.Add(1)
	go func() {
		walkDir(root, &sw, fileInfos)
	}()
	go func() {
		sw.Wait()
		close(fileInfos)
	}()
	return fileInfos
}

func walkDir(dir string, sw *sync.WaitGroup, fileInfos chan<- FileInfoType) {
	defer sw.Done()
	for _, entry := range dirents(dir) {
		entryName := entry.Name()
		if entry.IsDir() {
			sw.Add(1)
			subdir := filepath.Join(dir, entryName)
			go walkDir(subdir, sw, fileInfos)
		} else {
			p := filepath.Join(dir, entryName)
			contentMD5, _ := HashMD5(p)
			fileInfos <- FileInfoType{
				ValidHash:  contentMD5 != "",
				ContentMD5: contentMD5,
				Dir:        dir,
				Path:       p,
				PathOSS:    filepath.ToSlash(p),
				Name:       entryName,
			}
		}
	}
}

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("dirents error: %v\n", err)
		return nil
	}
	return entries
}
