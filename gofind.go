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

const (
	workersNo = 1
)

var (
	//dirQueue chan string = make(chan string, workersNo*100000)

	//fileQueue chan string = make(chan string, 10000)

	fileMux     *sync.Mutex     = &sync.Mutex{}
	dirMux      *sync.Mutex     = &sync.Mutex{}
	resultQueue *cn             = newCn()
	wg          *sync.WaitGroup = &sync.WaitGroup{}
	results     []string        = []string{}
)

type cn struct {
	open    bool
	c       chan string
	counter int
	mx      *sync.Mutex
}

func newCn() *cn {
	return &cn{true, make(chan string, workersNo*100000), 0, &sync.Mutex{}}
}

func (c *cn) Send(s string) {
	c.mx.Lock()
	c.counter++
	c.mx.Unlock()
	c.c <- s
}

func (c *cn) Get() (string, bool) {
	res, ok := <-c.c
	c.mx.Lock()
	c.counter--
	c.mx.Unlock()
	return res, ok
}

func (c *cn) Close() {
	if c.open {
		c.mx.Lock()
		close(c.c)
		c.open = false
		c.mx.Unlock()
	}
}

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
func Find(conf Config) []string {
	go print()
	wg.Add(1)
	go searchDir(conf, conf.StartPath)
	wg.Wait()
	return results
}

func searchDir(conf Config, dirPath string) {
	var (
		finfos   []os.FileInfo
		err      error
		fileName string
		filePath string
	)
	defer wg.Done()
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
			cbuff := conf.ContextBuffer
			indexes := conf.SearchContentPattern.FindAllIndex(fileCnt, -1)
			for _, index := range indexes {
				start := index[0]
				if start > cbuff {
					start = start - cbuff
				} else {
					start = 0
				}
				end := index[1]
				if end+cbuff < len(fileCnt) {
					end = end + cbuff
				} else {
					end = len(fileCnt)
				}

				fmt.Printf("------\n%s\n------\n", fileCnt[start:end])
			}
		}
	}
	return
}

func printRes(fileName string) {
	fmt.Println(fileName)
	results = append(results, fileName)
}
