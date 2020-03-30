package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var client http.Client

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := sc.Text()
		jsList := getJavascriptsFromUrl(url)
		fmt.Println(jsList)
	}
}

func getJavascriptsFromUrl(url string) []string {
	var jsList []string
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		re := regexp.MustCompile(`(https?://\S*?\.js)`)
		match := re.FindAllStringSubmatch(bodyString, -1)
		for _, element := range match {
			jsList = append(jsList, element[1])
		}
	}
	return jsList
}
