# Install command:
```
   go get github.com/kpawlik/gofind
```   
Version 0.2 
- search files by name
- search files by content - not yet

   

### Example of usage:

```go
   package main

    import (
	    "github.com/kpawlik/gofind"
	    "flag"
	    "fmt"
	    "os"
	    "time"
    )

    const VERSION = "0.2"

    func main() {
	    var (
		    // Start dir. Root dir path to start searching
		    dir = flag.String("d", "", "Start directory")
		    // Patter to find in file name
		    searchNamePattern = flag.String("np", "", "Pattern to seaerch in file name")
		    searchContentPattern = flag.String("cp", "", "Pattern to seaerch in file Content")
		    // print version
		    version = flag.Bool("v", false, "Version")
	    )
	
	    flag.Parse()
	    if *version {
		    fmt.Println("Version: ", VERSION)
		    return
	    }
	    if *searchNamePattern == "" {
		    flag.PrintDefaults()
		    return
	    }
	    if *dir == "" {
		    *dir, _ = os.Getwd()
	    }
	    testFind(*dir, *searchNamePattern, *searchContentPattern)
	    testGoFind(*dir, *searchNamePattern, *searchContentPattern)
    }

    func testFind(dir, searchNamePattern, searchContentPattern string){
	
	    fconf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern)
	    s := time.Now()
	    fileList := gofind.Find(fconf)
	    for _, f := range(fileList){
		    fmt.Println(f)
		    _ = f
	    }
	    fmt.Println("TOTOAL 1: ", time.Now().Sub(s))
    }
    //
    //
    //
    func testGoFind(dir, searchNamePattern, searchContentPattern string){
	    fconf := gofind.NewConfig(dir, searchNamePattern, searchContentPattern)
	    var ch = make(chan string)
	    go gofind.GoFind(fconf, ch)
	    s := time.Now()
	    stop := false 
	    for{
		    if stop{
			    break
		    }
		    select{
			    case fp, err := <-ch:
				    if stop = !err; stop {
					    break
				    }
				    _ = fp 
				    fmt.Println(fp)
		    }
	    }
	    fmt.Println("TOTOAL 2: ", time.Now().Sub(s))
    }
```
