// gofind.go
package gofind

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sync"
)

var (
	fileMux    = &sync.Mutex{}
	dirMux     = &sync.Mutex{}
	counterMux = &sync.Mutex{}
	wg         = &sync.WaitGroup{}
	results    = []string{}
	counter    int
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
	SearchByName         bool
	SearchByContent      bool
	Quiet                bool
	ShowContext          bool
	ContextBuffer        int
}

//
// NewConfig create new instance of Config
//
func NewConfig(startDir, fileNamePattern, contentPattern string, quiet bool, contextBuffer int) Config {
	var (
		snpRe, scpRe                  *regexp.Regexp
		searchByName, searchByContent bool
	)
	if searchByName = fileNamePattern != ""; searchByName {
		snpRe = regexp.MustCompile(fileNamePattern)

	}
	if searchByContent = contentPattern != ""; searchByContent {
		scpRe = regexp.MustCompile(contentPattern)
	}

	return Config{StartPath: startDir,
		SearchNamePattern:    snpRe,
		SearchContentPattern: scpRe,
		SearchByName:         searchByName,
		SearchByContent:      searchByContent,
		Quiet:                quiet,
		ShowContext:          contextBuffer > 0,
		ContextBuffer:        contextBuffer,
	}
}

// Find returns list of files which maches to patterns from conf
func Find(conf Config) ([]string, int) {
	wg.Add(1)
	go searchDir(conf, conf.StartPath)
	wg.Wait()
	return results, counter
}

func searchDir(conf Config, dirPath string) {
	var (
		finfos   []os.FileInfo
		err      error
		fileName string
		filePath string
	)
	defer wg.Done()
	incCounter()
	if finfos, err = readDir(dirPath); err != nil {
		if conf.Quiet && os.IsPermission(err) {
			return
		}
		fmt.Printf("Error reading file: %s (%v)\n", dirPath, err)
	}

	for _, finfo := range finfos {
		fileName = finfo.Name()
		filePath = path.Join(dirPath, fileName)
		fileNameMatch := false
		if conf.SearchByName {
			fileNameMatch = conf.SearchNamePattern.MatchString(fileName)
		}
		if conf.SearchByName && fileNameMatch && !conf.SearchByContent {
			printRes(filePath)
		}
		if finfo.IsDir() {
			wg.Add(1)
			go searchDir(conf, filePath)
		} else {
			incCounter()
			if conf.SearchByName && fileNameMatch && conf.SearchByContent {
				err = searchFile(conf, filePath)
			}
			if !conf.SearchByName && conf.SearchByContent {
				err = searchFile(conf, filePath)
			}
			if err != nil && !conf.Quiet && !os.IsPermission(err) {
				fmt.Printf("Error reading file %s (%v)\n", filePath, err)
			}
		}
	}
}

func readDir(dirPath string) (finfos []os.FileInfo, err error) {
	var (
		file *os.File
	)
	defer dirMux.Unlock()
	dirMux.Lock()
	if file, err = os.Open(dirPath); err != nil {
		return
	}
	finfos, err = file.Readdir(-1)
	file.Close()
	return

}

func searchFile(conf Config, filePath string) (err error) {
	var (
		fileCnt []byte
	)

	fileMux.Lock()
	defer fileMux.Unlock()
	if fileCnt, err = ioutil.ReadFile(filePath); err != nil {
		return
	}
	if conf.SearchContentPattern.Match(fileCnt) {
		printRes(filePath)
		if conf.ShowContext {
			printMatchContext(fileCnt, conf.SearchContentPattern, conf.ContextBuffer)
		}
	}
	return
}

func printMatchContext(content []byte, re *regexp.Regexp, bufferSize int) {

	indexes := re.FindAllIndex(content, -1)
	for _, index := range indexes {
		start := index[0]
		if start > bufferSize {
			start = start - bufferSize
		} else {
			start = 0
		}
		end := index[1]
		if end+bufferSize < len(content) {
			end = end + bufferSize
		} else {
			end = len(content)
		}
		fmt.Printf("------\n%s\n------\n", content[start:end])
	}
}
func printRes(fileName string) {
	fmt.Println(fileName)
	results = append(results, fileName)
}

func incCounter() {
	counterMux.Lock()
	counter++
	counterMux.Unlock()
}
