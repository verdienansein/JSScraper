package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

const BANNER = `
 _____ _______ _______ _______  ______ _______  _____  _______  ______
   |   |______ |______ |       |_____/ |_____| |_____] |______ |_____/
 __|   ______| ______| |_____  |    \_ |     | |       |______ |    \_

`

var client = &http.Client{
	Timeout: time.Duration(3 * time.Second),
}

func main() {
	fmt.Println(BANNER)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := sc.Text()
		color.Blue("Scraping URL: " + url)
		getJavascriptsFromUrl(url)
	}
}

func getJavascriptsFromUrl(url string) {
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyString := string(bodyBytes)
		re := regexp.MustCompile(`(https?://[a-zA-Z0-9\.\-/_]+?\.js)`)
		match := re.FindAllStringSubmatch(bodyString, -1)
		for _, element := range match {
			color.Green(element[1])
			getSecretsFromJS(element[1])
			fmt.Println("")
		}
		re = regexp.MustCompile(`["'](/[a-zA-Z0-9\.\-/_]+?\.js)`)
		match = re.FindAllStringSubmatch(bodyString, -1)
		for _, element := range match {
			absoluteUrl := url + element[1]
			color.Green(absoluteUrl)
			getSecretsFromJS(absoluteUrl)
			fmt.Println("")
		}
	}
}

func getSecretsFromJS(jsUrl string) {
	resp, err := client.Get(jsUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyString := string(bodyBytes)
		re := regexp.MustCompile(`(.{30})(token|auth|pass)(.{30})`)
		match := re.FindAllStringSubmatch(bodyString, -1)
		for _, element := range match {
			fmt.Print(element[1])
			color.New(color.FgRed).Print(element[2])
			fmt.Println(element[3])
		}
	}
}
