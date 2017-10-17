package gofind

import (
	"fmt"
	"os"
	"path/filepath"
)

// WalkFind search using path.Walk function
func WalkFind(conf Config) ([]string, int32) {
	filepath.Walk(conf.StartPath, func(filePath string, fi os.FileInfo, e error) (err error) {
		incCounter()
		fileNameMatch := false
		fileName := fi.Name()
		if conf.SearchByName {
			fileNameMatch = conf.SearchNamePattern.MatchString(fileName)
		}
		if conf.SearchByName && fileNameMatch {
			if conf.SearchByContent && !fi.IsDir() {
				err = searchFile(conf, filePath)
			}
		}
		if conf.SearchByName && fileNameMatch && !conf.SearchByContent {
			printRes(filePath)
		}
		if !conf.SearchByName && conf.SearchByContent && !fi.IsDir() {
			err = searchFile(conf, filePath)
		}
		if err != nil && !conf.Quiet && !os.IsPermission(err) {
			fmt.Printf("Error reading file %s (%v)\n", filePath, err)
		}
		return
	})
	return results, counter
}
