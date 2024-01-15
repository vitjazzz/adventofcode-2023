package adventutils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetFromUrl(url string, filterEmpty bool) (taskLines []string) {
	cookie := ""

	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Cookie", cookie)

	client := &http.Client{}

	resp, err := client.Do(request)
	fmt.Println(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	task := string(body)
	lines := strings.Split(task, "\n")
	if !filterEmpty {
		return lines
	}
	var filteredLines []string
	for _, line := range lines {
		if line != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
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
