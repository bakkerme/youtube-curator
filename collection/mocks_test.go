package collection

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

type mockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (f mockFileInfo) Name() string {
	return f.name
}
func (f mockFileInfo) Size() int64 {
	return f.size
}
func (f mockFileInfo) Mode() os.FileMode {
	return 0
}
func (f mockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f mockFileInfo) IsDir() bool {
	return f.isDir
}
func (f mockFileInfo) Sys() interface{} {
	return nil
}

type mockDirReaderProvider struct {
	expectedDirname            *string
	expectedFilename           *string
	t                          *testing.T
	shouldErrorReadDir         bool
	shouldErrorReadFile        bool
	returnReadDirValue         *[]os.FileInfo
	returnReadFileValue        *[]byte
	returnReadFileValueForPath map[string][]byte
}

func (mdr *mockDirReaderProvider) ReadDir(dirname string) ([]os.FileInfo, error) {
	if mdr.shouldErrorReadDir {
		return nil, errors.New("oh the humanity")
	}

	if mdr.expectedDirname != nil {
		if mdr.t == nil {
			panic("Please provide mockDirReaderProvider a T")
		}

		if dirname != *mdr.expectedDirname {
			mdr.t.Errorf("dirname provided to ReadDir was not the expected value. Expected %s, got %s", *mdr.expectedDirname, dirname)
		}
	}

	if mdr.returnReadDirValue != nil {
		return *mdr.returnReadDirValue, nil
	}

	return nil, nil
}

func (mdr *mockDirReaderProvider) ReadFile(path string) ([]byte, error) {
	if mdr.shouldErrorReadFile {
		return nil, errors.New("oh the humanity")
	}

	if mdr.expectedFilename != nil {
		if mdr.t == nil {
			panic("Please provide mockDirReaderProvider a T")
		}

		if path != *mdr.expectedFilename {
			mdr.t.Errorf("path provided to ReadFile was not the expected value. Expected %s, got %s", *mdr.expectedFilename, path)
		}
	}

	if mdr.returnReadFileValue != nil {
		return *mdr.returnReadFileValue, nil
	}

	if mdr.returnReadFileValueForPath[path] != nil {
		return mdr.returnReadFileValueForPath[path], nil
	}

	return nil, nil
}

type mockConfigProvider struct {
	videoDirPath  string
	youtubeAPIKey string
}

func (cp mockConfigProvider) GetValue(key string) (string, error) {
	if key == "VIDEO_DIR_PATH" {
		return cp.videoDirPath, nil
	}

	if key == "YOUTUBE_API_KEY" {
		return cp.youtubeAPIKey, nil
	}

	return "", fmt.Errorf("mockConfigProvider could not find key %s", key)
}
