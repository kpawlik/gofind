package gofind

import (
	"os"
	"path"
	"testing"
)


var filesToCreate = []string{"fileOne.txt", "fileOne2.txt", "fff.txt", "fileO/fileO.txt"}
const searchPattern = "^.*fileO.*$"
const expectedCount = 4

//return current dir
func getCd(t *testing.T) string{
	if cwd, err := os.Getwd(); err != nil {
		t.Error(err)
	}else{
		return cwd
	}
	return ""
}

//create temp files
func setup(t *testing.T) {
	cwd := getCd(t)
	for _, f := range filesToCreate {
		if dir, _ := path.Split(f); dir != "" {
			os.MkdirAll(path.Join(cwd, dir), os.ModePerm)
		}
		os.Create(f)
	}
}
//remove temp files
func tearDown(t *testing.T) {
	cwd := getCd(t)
	for _, f := range filesToCreate {
		if err := os.RemoveAll(f); err != nil{
				t.Error(err)
		}
		if dir, _ := path.Split(f); dir != "" {
			if err := os.Remove(path.Join(cwd, dir)); err!=nil{
				t.Error(err)
			}
		}
	}
	
}
//test static search function
func TestFileSearch(t *testing.T) {
	var (
		cwd string
		err error
	)
	setup(t)
	defer tearDown(t)

	if cwd, err = os.Getwd(); err != nil {
		t.Error(err)
	}
	cfg := NewConfig(cwd, searchPattern, "")
	list := Find(cfg)
	if len(list) != expectedCount {
		t.Error("Serach gofiles by name should be ", expectedCount, " not ", len(list))
	}
}
