package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	gohead "github.com/CasperGN/GoHead"
)

var (
	target  string
	targets string
	exclude string
	outdir  string
	threads int
	silent  bool
)

func init() {
	flag.StringVar(&target, "target", "", "Supply single target for probing.")
	flag.StringVar(&targets, "targets", "", "Supply a file of targets seperated by newlines.")
	flag.StringVar(&exclude, "exclude", "", "Supply a file of headers to exclude seperated by newlines.")
	flag.StringVar(&outdir, "outdir", "", "Supply a directory to output the result to. Writes 1 file per supplied target.")
	flag.IntVar(&threads, "threads", 5, "Number of threads")
	flag.BoolVar(&silent, "silent", false, "Print header (default false).")
}

func main() {

	var (
		asset           io.ReadCloser
		file            io.ReadCloser
		excludedHeaders []string
		exclusion       bool
		output          bool
		err             error
	)
	exclusion = false
	output = false

	flag.Parse()

	if target == "" && targets == "" {
		printHeader()
		flag.Usage()
		return
	}

	if outdir != "" {
		output = true
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

	intIP, _ := regexp.Compile(`(10(\.[1-9]{1}[0-9]{1,2}){3})|(172(\.[1-9]{1}[0-9]{1,2}){3})|(192(\.[1-9]{1}[0-9]{1,2}){3})`)

	var waitGroup sync.WaitGroup
	assets := make(chan string)

	for i := 0; i < threads; i++ {
		waitGroup.Add(1)
		go func() {
			for asset := range assets {
				data := ""
				var regMatch []string
				body, headers, target := gohead.Probe(asset)
				if len(headers) > 0 {
					for _, match := range intIP.FindAllString(body, -1) {
						regMatch = append(regMatch, match)
					}
					fmt.Printf("%s\n", target)
					for key, value := range headers {
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
						for _, match := range intIP.FindAllString(option, -1) {
							regMatch = append(regMatch, match)
						}
						fmt.Printf("%s: %s\n", key, option)
						data += key + ": " + option + "\n"
					}
				}
				ips := ""
				if len(regMatch) > 0 {
					for _, ip := range regMatch {
						if strings.Contains(ips, ip) {
							continue
						}
						ips += ip + " "
					}
					fmt.Println("")
					fmt.Printf("Internal IPs: %s", ips)
				}
				data += ips

				if output {
					filename := strings.Replace(strings.Split(asset, ":")[1], "/", "", -1)
					if outdir[len(outdir)-1:] != "/" {
						outdir += "/"
					}
					f, err := os.Create(outdir + filename)
					if err != nil {
						fmt.Println("Couldn't write to directory specified")
					}
					defer f.Close()

					_, err = f.WriteString(data)
					if err != nil {
						fmt.Println("Couldn't write to directory specified")
					}
					f.Sync()
				}
				fmt.Println("")
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
