package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getFromUrl(url string) (taskLines []string) {
	cookie := "session=53616c7465645f5f07b1603c10dd6df60ca1b1d38501c9be6e50e16153cb8acd51f7cb08f25cfac4cdf4d6556df456401023ba1ea5d9b431bbfc52be2a17b045"

	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Cookie", cookie)

	client := &http.Client{}

	resp, err := client.Do(request)
	fmt.Println(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	task := string(body)
	return strings.Split(task, "\n")
}

func getFromFile(fn string) (taskLines []string) {

	readFile, err := os.Open(fn)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		taskLines = append(taskLines, fileScanner.Text())
	}

	defer readFile.Close()

	return taskLines
}
