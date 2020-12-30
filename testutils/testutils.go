package testutils

import (
	"errors"
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

// MockDirReader provides a mock DirReaderProvider with
// a number of properties available to tailor the output of
// filesystem operations supported through the DirReader
type MockDirReader struct {
	ExpectedDirname            *string
	ExpectedFilename           *string
	T                          *testing.T
	ShouldErrorReadDir         bool
	ShouldErrorReadFile        bool
	ReturnReadDirValue         *[]os.FileInfo
	ReturnReadFileValue        *[]byte
	ReturnReadFileValueForPath map[string][]byte
	ReturnHomeDirPath          *string
}

// ReadDir mocks the DirReader ReadDir function
func (mdr *MockDirReader) ReadDir(dirname string) ([]os.FileInfo, error) {
	if mdr.ShouldErrorReadDir {
		return nil, errors.New("oh the humanity")
	}

	if mdr.ExpectedDirname != nil {
		if mdr.T == nil {
			panic("Please provide MockDirReader a T")
		}

		if dirname != *mdr.ExpectedDirname {
			mdr.T.Errorf("dirname provided to ReadDir was not the expected value. Expected %s, got %s", *mdr.ExpectedDirname, dirname)
		}
	}

	if mdr.ReturnReadDirValue != nil {
		return *mdr.ReturnReadDirValue, nil
	}

	return nil, nil
}

// ReadFile mocks the DirReader ReadFile function
func (mdr *MockDirReader) ReadFile(path string) ([]byte, error) {
	if mdr.ShouldErrorReadFile {
		return nil, errors.New("oh the humanity")
	}

	if mdr.ExpectedFilename != nil {
		if mdr.T == nil {
			panic("Please provide MockDirReader a T")
		}

		if path != *mdr.ExpectedFilename {
			mdr.T.Errorf("path provided to ReadFile was not the expected value. Expected %s, got %s", *mdr.ExpectedFilename, path)
		}
	}

	if mdr.ReturnReadFileValue != nil {
		return *mdr.ReturnReadFileValue, nil
	}

	if mdr.ReturnReadFileValueForPath[path] != nil {
		return mdr.ReturnReadFileValueForPath[path], nil
	}

	return nil, nil
}

// GetHomeDirPath mocks the path to the user's home directory
func (mdr *MockDirReader) GetHomeDirPath() string {
	if mdr.ReturnHomeDirPath == nil {
		mdr.T.Error("Attempted to run GetHomeDirPath but no ReturnHomeDirPath was provided")
	}

	return *mdr.ReturnHomeDirPath
}
