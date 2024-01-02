package adventutils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetFromUrl(url string) (taskLines []string) {
	cookie := ""

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
