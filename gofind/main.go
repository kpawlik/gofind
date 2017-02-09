package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kpawlik/gofind"
)

const VERSION = "0.2"

var (
	// Start dir. Root dir path to start searching
	//dir = flag.String("d", "", "Start directory")
	dir                  string
	searchNamePattern    string
	searchContentPattern string
	version              bool
	timeStats            bool
	// Patter to find in file name
	//searchNamePattern    = flag.String("np", "", "Pattern to seaerch in file name")
	//searchContentPattern = flag.String("cp", "", "Pattern to seaerch in file Content")
	// print version
	//version = flag.Bool("v", false, "Version")
)

func init() {
	var (
		help bool
	)
	flag.StringVar(&dir, "d", "", "Start directory. Default current")
	flag.StringVar(&searchContentPattern, "c", "", "File content seaerch pattern (regexp)")
	flag.StringVar(&searchNamePattern, "n", "", "File name patters to search (regexp)")
	flag.BoolVar(&timeStats, "s", false, "Print summary")
	flag.BoolVar(&help, "h", false, "Print help")

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

	fconf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern)
	s := time.Now()
	results := gofind.Find(fconf)
	if timeStats {
		fmt.Printf("Found: %d\n", len(results))
		fmt.Println("Time: ", time.Now().Sub(s))
	}
}
