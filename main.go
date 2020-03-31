package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const BANNER = `
 _____ _______ _______ _______  ______ _______  _____  _______  ______
   |   |______ |______ |       |_____/ |_____| |_____] |______ |_____/
 __|   ______| ______| |_____  |    \_ |     | |       |______ |    \_

`

var client http.Client
var outputLength string
var keywordsList string

func main() {
	fmt.Println(BANNER)

	var keywords string
	flag.StringVar(&keywords, "k", "auth,pass,token", "comma separeted keywords to find in javascripts (Default: auth,pass,token)")
	flag.StringVar(&outputLength, "l", "30", "length of the grepped output (Default: 30)")
	var to int
	flag.IntVar(&to, "t", 10, "timeout (seconds)")
	flag.Parse()

	client = http.Client{
		Timeout: time.Duration(to) * time.Second,
	}

	keywordsList = strings.Replace(keywords, ",", "|", -1)
	color.Magenta("Using regex: " + keywordsList)
	fmt.Println("")

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
		getSecretsFromAbsoluteUrls(bodyString)
		getSecretsFromRelativeUrls(bodyString, url)
	}
}

func getSecretsFromAbsoluteUrls(bodyString string) {
	re := regexp.MustCompile(`(https?://[a-zA-Z0-9\.\-/_]+?\.js)`)
	match := re.FindAllStringSubmatch(bodyString, -1)
	for _, element := range match {
		color.Green(element[1])
		getSecretsFromJS(element[1])
		fmt.Println("")
	}

}

func getSecretsFromRelativeUrls(bodyString string, baseUrl string) {
	re := regexp.MustCompile(`["'](/[a-zA-Z0-9\.\-/_]+?\.js)`)
	match := re.FindAllStringSubmatch(bodyString, -1)
	for _, element := range match {
		absoluteUrl := baseUrl + element[1]
		color.Green(absoluteUrl)
		getSecretsFromJS(absoluteUrl)
		fmt.Println("")
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
		re := regexp.MustCompile(`(?i)(.{` + outputLength + `})(` + keywordsList + `)(.{` + outputLength + `})`)
		match := re.FindAllStringSubmatch(bodyString, -1)
		for _, element := range match {
			fmt.Print(element[1])
			color.New(color.FgRed).Print(element[2])
			fmt.Println(element[3])
		}
	}
}
