package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	secrets bool
)

func init() {
	flag.StringVar(&target, "target", "", "Supply single target for probing.")
	flag.StringVar(&targets, "targets", "", "Supply a file of targets seperated by newlines.")
	flag.StringVar(&exclude, "exclude", "", "Supply a file of headers to exclude seperated by newlines.")
	flag.StringVar(&outdir, "outdir", "", "Supply a directory to output the result to. Writes 1 file per supplied target.")
	flag.IntVar(&threads, "threads", 5, "Number of threads")
	flag.BoolVar(&silent, "silent", false, "Print header (default false).")
	flag.BoolVar(&secrets, "secrets", false, "Search JavaScript files for keys, passwords or secrets (default false)")
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
	jsFiles, _ := regexp.Compile(`src="([a-zA-Z0-9\./-_]+.js)"`)
	urlFind, _ := regexp.Compile(`((/[a-zA-Z0-9\.]+){2,10})`)
	/*awsKey, _ := regexp.Compile(`"(AKIA[0-9A-Z]{16})"`)
	awsSecKey, _ := regexp.Compile(`"([0-9a-zA-Z/+]{40})"`)
	azureSub, _ := regexp.Compile(`"([a-z0-9]{32})"`)
	mailGun, _ := regexp.Compile(`"(key-[0-9a-zA-Z]{32})"`)
	fbKey, _ := regexp.Compile(`"([0-9]{13,17})"`)
	fbSecKey, _ := regexp.Compile(`"([0-9a-f]{32})"`)
	passwd, _ := regexp.Compile(`"((Pass:|pass:|Pass=|pass=|Password:|password:|Password=|password=)(.*[!@#$%^&*a-zA-Z0-9]))"`)*/

	var waitGroup sync.WaitGroup
	assets := make(chan string)

	for i := 0; i < threads; i++ {
		waitGroup.Add(1)
		go func() {
			for asset := range assets {
				data := ""
				var regMatch []string
				urlMap := make(map[string][]string)
				//keysFound := make(map[string]map[string]string)
				body, headers, target := gohead.Probe(asset)
				if len(headers) > 0 {
					if secrets {
						for _, match := range jsFiles.FindAllStringSubmatch(body, -1) {
							//fmt.Println(body)
							var matchedUrl = ""
							if strings.HasPrefix(match[1], "/") {
								matchedUrl += match[1]
							} else {
								matchedUrl += "/"
								matchedUrl += match[1]
							}
							jsBody, err := http.Get(asset + matchedUrl)
							if err != nil {
								//fmt.Printf("Error on GET for %s%s: %s", asset, match[1], err)
								continue
							}

							jsFile, err := ioutil.ReadAll(jsBody.Body)
							jsBody.Body.Close()
							if err != nil {
								//fmt.Printf("Error on reading body from %s%s: %s", asset, match[1], err)
								continue
							}
							urlMap[asset+matchedUrl] = append(urlMap[asset+matchedUrl], asset+matchedUrl)
							for _, match := range urlFind.FindAllStringSubmatch(string(jsFile), -1) {
								//fmt.Println(match)
								//fmt.Println(match[0])
								urlMap[asset+matchedUrl] = append(urlMap[asset+matchedUrl], match[0])
								//for _, url := range match {
								//	urlMap[asset+matchedUrl] = append(urlMap[asset+matchedUrl], url)
								//}

							}
							/*for _, match := range awsKey.FindAllStringSubmatch(string(jsFile), -1) {
								keysFound[asset+matchedUrl][match[1]] = "awskey"
							}
							for _, match := range awsSecKey.FindAllStringSubmatch(string(jsFile), -1) {
								fmt.Println("AwsSecKey:")
								fmt.Println(match[1])
							}
							for _, match := range azureSub.FindAllStringSubmatch(string(jsFile), -1) {
								fmt.Println("Azure:")
								fmt.Println(match[1])
							}
							for _, match := range mailGun.FindAllStringSubmatch(string(jsFile), -1) {
								fmt.Println("MailGun:")
								fmt.Println(match[1])
							}
							for _, match := range fbKey.FindAllStringSubmatch(string(jsFile), -1) {
								fmt.Println("FbKey:")
								fmt.Println(match[1])
							}
							for _, match := range fbSecKey.FindAllStringSubmatch(string(jsFile), -1) {
								fmt.Println("FbSecKey:")
								fmt.Println(match[1])
							}
							for _, match := range passwd.FindAllStringSubmatch(string(jsFile), -1) {
								fmt.Println("Passwd:")
								fmt.Println(match[1])
							}*/
						}
					}
					for _, match := range intIP.FindAllString(body, -1) {
						regMatch = append(regMatch, match)
					}
					//fmt.Printf("%s\n", target)
					data += "target: " + target

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
						//fmt.Printf("header: %s: %s\n", key, option)
						data += "\nheader: " + key + ": " + option
						//fmt.Printf("%s: %s\n", key, option)
						//data += key + ": " + option + "\n"
					}
				}
				ips := ""
				if len(regMatch) > 0 {
					//data += "\n"
					ips += "\ninternal IPs:"
					for _, ip := range regMatch {
						if strings.Contains(ips, ip) {
							continue
						}
						ips += ip + " "
					}
					//fmt.Println("")
					//fmt.Printf(ips)
					data += ips
				}

				if secrets {
					if len(urlMap) > 0 {
						//data += "\nPossible endpoints:\n"
						for key, value := range urlMap {
							for _, url := range value {
								if suffixes(url) {
									data += "\nfile: " + key + ": " + url
								} else {
									data += "\nendpoint: " + key + ": " + url
								}
							}
							//data += key + ": " + value + "\n"
						}
					}
					/*if len(keysFound) > 0 {
						data += "\nPossible keys:\n"
						for key, value := range keysFound {
							for k, v := range value {
								data += v + ": " + k + " in " + key
							}
						}
					}*/
				}

				fmt.Println(data)

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
				//fmt.Println("")
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

func suffixes(url string) bool {
	suffixArray := [11]string{".js", ".html", ".asp", ".aspx", ".php", ".htm", ".gif", ".jpg", ".jpeg", ".png", ".txt"}
	for _, suffix := range suffixArray {
		if strings.HasSuffix(url, suffix) {
			return true
		}
	}
	return false
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
