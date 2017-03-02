package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kpawlik/gofind"
)

var (
	// TODO:
	// quiet errors permision denite
	// print found context
	dir                  string
	searchNamePattern    string
	searchContentPattern string
	version              bool
	timeStats            bool
	quiet                bool
	help                 bool
	context              int
	stype                int
	searches             = map[int]findFunc{0: gofind.Find, 1: gofind.WalkFind}
)

type findFunc func(gofind.Config) ([]string, int32)

func init() {

	flag.StringVar(&dir, "dir", "", "Start directory (current by default)")
	flag.StringVar(&searchContentPattern, "content", "", "File content seaerch pattern (regexp)")
	flag.StringVar(&searchNamePattern, "name", "", "File name pattern (regexp)")
	flag.BoolVar(&timeStats, "stat", false, "Print time stats")
	flag.BoolVar(&help, "help", false, "Print help")
	flag.BoolVar(&quiet, "quiet", false, "Quiet permission denite errors")
	flag.IntVar(&context, "context", 0, "Number of chars of find context ")
	flag.IntVar(&stype, "type", 0, `Search types
	0 - (defualt) concurrent (fastest, but on Linux, for large no of files to search, may coused error 'to many open files'
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
	fmt.Println(`\n Hint: To case insensitive seach use (?i) prefix in regexp pattern`)
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
	fconf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern, quiet, context)
	s := time.Now()
	if searchFunc, ok = searches[stype]; !ok {
		printHelp()
		return
	}
	results, counter = searchFunc(fconf)
	if timeStats {
		fmt.Printf("Searched: %d. Found: %d\nTime: %s\n", counter, len(results), time.Now().Sub(s))
	}
}
