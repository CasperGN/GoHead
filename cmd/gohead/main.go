package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	gohead "github.com/CasperGN/GoHead"
)

var (
	target  string
	targets string
	threads int
	silent  bool
)

func init() {
	flag.StringVar(&target, "target", "", "Supply single target for probing.")
	flag.StringVar(&targets, "targets", "", "Supply a file of targets seperated by newlines.")
	flag.IntVar(&threads, "threads", 5, "Number of threads")
	flag.BoolVar(&silent, "silent", false, "Print header (default false).")
}

func main() {

	var (
		asset io.ReadCloser
		err   error
	)

	flag.Parse()

	if target == "" && targets == "" {
		printHeader()
		flag.Usage()
		return
	}

	if target != "" {
		asset = ioutil.NopCloser(strings.NewReader(target))
	} else {
		asset, err = os.Open(targets)
		if err != nil {
			fmt.Printf("Error: Cannot read targets file. Name: %s, error: %s", targets, err)
		}
	}
	defer asset.Close()

	if !silent {
		printHeader()
	}

	var waitGroup sync.WaitGroup
	assets := make(chan string)

	for i := 0; i < threads; i++ {
		waitGroup.Add(1)
		go func() {
			for asset := range assets {
				result, target := gohead.Probe(asset)
				if len(result) > 0 {
					fmt.Printf("%s\n", target)
					for key, value := range result {
						if contains(key) {
							// Header is irrelevant
							continue
						}
						option := ""
						for _, val := range value {
							option += val
						}
						fmt.Printf("%s: %s\n", key, option)
					}
				}

			}
			waitGroup.Done()
		}()
	}

	fRead := bufio.NewScanner(asset)
	for fRead.Scan() {
		assets <- fRead.Text()
	}
	close(assets)
	waitGroup.Wait()
}

func contains(header string) bool {
	benignHeaders := [12]string{
		"Content-Language", "X-Ua-Compatible", "Last-Modified", "Date", "Etag", "Connection", "Content-Length", "Pragma",
		"Expires", "Server-Timing", "X-Content-Type-Options", "Cache-Control"
	}
	for _, benign := range benignHeaders {
		if benign == header {
			return true
		}
	}
	return false
}

func printHeader() {
	header :=
		`
              ______      __  __               __ 
             / ____/___  / / / /__  ____ _____/ / 
            / / __/ __ \/ /_/ / _ \/ __ \/ __  /  
        __ / /_/ / /_/ / __  /  __/ /_/ / /_/ / __
      _/_/ \____/\____/_/ /_/\___/\__,_/\__,_/_/_/
    _/_/___________________________________ _/_/  
  _/_//_____/_____/_____/_____/_____/_____//_/    
 /_/                                     /_/      			  
		`
	fmt.Println(header)
}
