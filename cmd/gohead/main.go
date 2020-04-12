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
	exclude string
	threads int
	silent  bool
)

func init() {
	flag.StringVar(&target, "target", "", "Supply single target for probing.")
	flag.StringVar(&targets, "targets", "", "Supply a file of targets seperated by newlines.")
	flag.StringVar(&exclude, "exclude", "", "Supply a file of headers to exclude seperated by newlines.")
	flag.IntVar(&threads, "threads", 5, "Number of threads")
	flag.BoolVar(&silent, "silent", false, "Print header (default false).")
}

func main() {

	var (
		asset           io.ReadCloser
		file            io.ReadCloser
		excludedHeaders []string
		exclusion       bool
		err             error
	)
	exclusion = false

	flag.Parse()

	if target == "" && targets == "" {
		printHeader()
		flag.Usage()
		return
	}

	if exclude != "" {
		exclusion = true
		file, err = os.Open(exclude)
		if err != nil {
			fmt.Printf("Error: Cannot read targets file. Name: %s, error: %s", targets, err)
		}
		fRead := bufio.NewScanner(file)
		for fRead.Scan() {
			excludedHeaders = append(excludedHeaders, fRead.Text())
		}
		file.Close()
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
						if exclusion {
							if contains(excludedHeaders, key) {
								// Header is irrelevant
								continue
							}
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

func contains(excludedHeaders []string, header string) bool {
	for _, excluded := range excludedHeaders {
		if excluded == header {
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
