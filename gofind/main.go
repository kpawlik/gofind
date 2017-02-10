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
	context              int
)

func init() {
	var (
		help bool
	)
	flag.StringVar(&dir, "d", "", "Start directory. Default current")
	flag.StringVar(&searchContentPattern, "c", "", "File content seaerch pattern (regexp)")
	flag.StringVar(&searchNamePattern, "n", "", "File name patters to search (regexp)")
	flag.BoolVar(&timeStats, "stat", false, "Print summary")
	flag.BoolVar(&help, "h", false, "Print help")
	flag.BoolVar(&quiet, "q", false, "Quiet permission denite errors")
	flag.IntVar(&context, "context", 0, "Show result context (slower)")
	flag.Parse()
	if help {
		fmt.Println(`Find file / search file content
gofind -d [SEARCH ROOT] -n [FILE NAME PATTERN] -c [CONTENT TO SEARCH]
Params:
`)
		flag.PrintDefaults()
		fmt.Println(`To case insensitive seach use (?i) prefix in regexp pattern`)
		os.Exit(0)
		return
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

func main() {
	fconf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern, quiet, context)
	s := time.Now()
	results := gofind.Find(fconf)
	if timeStats {
		fmt.Printf("Found: %d\n", len(results))
		fmt.Println("Time: ", time.Now().Sub(s))
	}
}
