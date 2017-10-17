package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kpawlik/gofind"
)

var (
	dir                  string
	searchNamePattern    string
	searchContentPattern string
	version              bool
	timeStats            bool
	quiet                bool
	help                 bool
	context              int
	showLine             bool
	includeSubdirs       bool
	searchType           int
	searches             = map[int]findFunc{0: gofind.Find, 1: gofind.WalkFind}
)

type findFunc func(gofind.Config) ([]string, int32)

func init() {

	flag.StringVar(&dir, "dir", "", "Start directory (current by default)")
	flag.StringVar(&searchContentPattern, "content", "", "File content seaerch pattern (regexp)")
	flag.StringVar(&searchNamePattern, "name", "", "File name pattern (regexp)")
	flag.BoolVar(&timeStats, "stat", false, "Print time stats")
	flag.BoolVar(&help, "help", false, "Print help")
	flag.BoolVar(&quiet, "quiet", false, "Quiet permission denied errors")
	flag.BoolVar(&showLine, "line", false, "Show line when pattent was found")
	flag.BoolVar(&includeSubdirs, "subdirs", true, "Search subdirs")
	flag.IntVar(&context, "context", 0, "Number of chars of find context ")
	flag.IntVar(&searchType, "type", 0, `Search types
	0 - (default) concurrent (fastest, but on Linux, for large no of files to search, may caused error 'to many open files'
	1 - walk from standard lib)`)
	flag.Parse()
	if help {
		printHelp()
	}
	if searchNamePattern == "" && searchContentPattern == "" {
		flag.PrintDefaults()
		os.Exit(0)
		return
	}
	if dir == "" {
		dir, _ = os.Getwd()
	}
}

func printHelp() {
	fmt.Println(`Find file with name matches to pattern or search file content for pattern.
gofind -d [SEARCH ROOT] -n [FILE NAME PATTERN] -c [CONTENT TO SEARCH] -type [0,1]

Params:
`)
	flag.PrintDefaults()
	fmt.Println(`\n Hint: To case insensitive search use (?i) prefix in regexp pattern`)
	os.Exit(0)
	return

}

func main() {
	var (
		results    []string
		counter    int32
		searchFunc findFunc
		ok         bool
	)
	conf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern, quiet, context, showLine, includeSubdirs)
	s := time.Now()
	if searchFunc, ok = searches[searchType]; !ok {
		printHelp()
		return
	}
	results, counter = searchFunc(conf)
	if timeStats {
		fmt.Printf("Searched: %d. Found: %d\nTime: %s\n", counter, len(results), time.Now().Sub(s))
	}
}
