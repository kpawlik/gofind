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
	includeGzip          bool
	includeSubdirs       bool
	searchType           int
	searches             = map[int]findFunc{0: gofind.Find, 1: gofind.WalkFind}
)

type findFunc func(gofind.Config) ([]string, int32)

func init() {

	flag.StringVar(&dir, "dir", "", "Start directory (by default current directory)")
	flag.StringVar(&searchContentPattern, "content", "", "File content search pattern (regexp)")
	flag.StringVar(&searchNamePattern, "name", "", "File name pattern (regexp)")
	flag.BoolVar(&timeStats, "stat", false, "Print time stats")
	flag.BoolVar(&help, "help", false, "Print help")
	flag.BoolVar(&quiet, "quiet", false, "Quiet permission denied errors")
	flag.BoolVar(&showLine, "line", false, "Show line when pattent was found")
	flag.BoolVar(&version, "version", false, "Print version no")
	flag.BoolVar(&includeSubdirs, "subdirs", true, "Search subdirs")
	flag.BoolVar(&includeGzip, "gzip", false, "Include gzip")
	flag.IntVar(&context, "context", 0, "Number of chars of find context ")
	flag.IntVar(&searchType, "type", 0, `Search types
	0 - (default) concurrent (fastest, but on Linux, for large no of files to search, may caused error 'to many open files'
	1 - walk from standard lib)`)
	flag.Parse()
	if help {
		printHelp()
	}
	if version {
		fmt.Println("20180517A")
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
	fmt.Println(`Find file with name matches to pattern or/and search file content for pattern.

	./gofind --dir [SEARCH ROOT] --name [FILE NAME PATTERN] --content [CONTENT TO SEARCH]

To case insensitive search use (?i) prefix in regexp pattern. 
Example:
This command will lists all files with extensions csv, CSV, Csv etc., which contains text "case ins text". Matching will be case insensitive. 
	./gofind --name "(?).*csv$" --content "(?)case ins text"

To search only directory, without subdirectories  use --subdirs=false
Example:
This command will lists all files with extension csv, only from directory /tmp/files

	./gofind --dir /tmp/files/ --name .*csv$ --subdirs=false

Parameters:
`)
	flag.PrintDefaults()
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
	conf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern, quiet, context, showLine, includeSubdirs, includeGzip)
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
