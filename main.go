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
			re := regexp.MustCompile(`(http.*\.js)`)
			match := re.FindStringSubmatch(bodyString)
			fmt.Println(match[1])
		}
	}
}
