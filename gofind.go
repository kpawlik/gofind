// gofind2.go.go
package gofind

import (
	"os"
	"path"
	"regexp"
)

// Config struct of find funcionalitu
//    StartPath - root directory to start search
//    SearchPatter - regular expression with patter to search
//    PrintErrors - print errors about opennung files, access denied etc
//    PrintStatistics - print stats how long search was made
type Config struct {
	StartPath            string
	SearchNamePattern    *regexp.Regexp
	SearchContentPattern *regexp.Regexp
}

//
// NewConfig create new instance of Config 
//
func NewConfig(startDir string, snp, scp string) Config {
	snpRe := regexp.MustCompile(snp)
	scpRe := regexp.MustCompile(scp)
	return Config{startDir, snpRe, scpRe}
}

//
// Func returns list of files which maches to patterns from conf
//
func Find(conf Config) []string {
	// Chanel to recive names of matching files
	var namesCh = make(chan string)
	// Chanel to recive names of dir names 
	var dirCh = make(chan string)
	// chanel to recive finished dirs
	var processedCh = make(chan string)
	files := process(conf, namesCh, dirCh, processedCh)
	return files
}

//
// Function gets channel as parameters. Sends to outCh files
// which matches to patters from conf
func GoFind(conf Config, outCh chan string) {
	// Chanel to recive names of dir names 
	var dirCh = make(chan string)
	// chanel to recive finished dirs
	var processedCh = make(chan string)

	goprocess(conf, outCh, dirCh, processedCh)

}

//
//
//
func goprocess(conf Config, namesCh, dirCh, processedCh chan string) {
	startPath := conf.StartPath
	re := conf.SearchNamePattern
	fi := processDir(startPath)
	balancer := make(map[string]bool, 10)
	go processList(startPath, fi, re, namesCh, dirCh, processedCh)
	balancer[startPath] = true
	for {
		if len(balancer) <= 0 {
			break
		}
		select {
		case dir := <-dirCh:
			fi := processDir(dir)
			go processList(dir, fi, re, namesCh, dirCh, processedCh)
			balancer[dir] = true
		case dir := <-processedCh:
			delete(balancer, dir)
		}
	}
	close(namesCh)
}

//
//
//
func process(conf Config, namesCh, dirCh, processedCh chan string) []string {
	startPath := conf.StartPath
	re := conf.SearchNamePattern
	fi := processDir(startPath)
	balancer := make(map[string]bool, 10)
	var fileList []string
	go processList(startPath, fi, re, namesCh, dirCh, processedCh)
	balancer[startPath] = true
	for {
		if len(balancer) <= 0 {
			break
		}
		select {
		case dir := <-dirCh:
			fi := processDir(dir)
			go processList(dir, fi, re, namesCh, dirCh, processedCh)
			balancer[dir] = true
		case dir := <-processedCh:
			delete(balancer, dir)
		case fpath := <-namesCh:
			fileList = append(fileList, fpath)
		}
	}
	close(namesCh)
	return fileList
}

//
// Process directory given as a path
//
func processDir(p string) []os.FileInfo {
	var (
		f   *os.File
		fi  []os.FileInfo
		err error
	)

	if f, err = os.Open(p); err != nil {
		return fi
	}
	defer f.Close()
	if fi, err = f.Readdir(-1); err != nil {
		return fi
	}
	return fi
}

//
// Proces list of file infos to search re in files
//
func processList(basePath string, list []os.FileInfo, re *regexp.Regexp, namesCh, dirCh, processedCh chan string) {
	for _, fileInfo := range list {
		if fileInfo.IsDir() {
			dirCh <- path.Join(basePath, fileInfo.Name())
		}
		if re.MatchString(fileInfo.Name()) {
			namesCh <- path.Join(basePath, fileInfo.Name())
		}
	}
	processedCh <- basePath
}
