package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// FileInfoType is a type which contains dir and os.FileInfo
type FileInfoType struct {
	Dir     string
	Path    string
	PathOSS string
	Info    os.FileInfo
}

// WalkDir get sub files of target dir
func WalkDir(root string, fileInfos chan FileInfoType) {
	var sw sync.WaitGroup
	sw.Add(1)
	go func() {
		walkDir(root, &sw, fileInfos)
	}()
	go func() {
		sw.Wait()
		close(fileInfos)
	}()
}

func walkDir(dir string, sw *sync.WaitGroup, fileInfos chan<- FileInfoType) {
	defer sw.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			sw.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, sw, fileInfos)
		} else {
			p := filepath.Join(dir, entry.Name())
			fileInfos <- FileInfoType{
				Dir:     dir,
				Path:    p,
				PathOSS: path.Join(strings.Split(p, string(filepath.Separator))...),
				Info:    entry,
			}
		}
	}
}

var sema = make(chan struct{}, 20)

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
