package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	filehandler "url-downloader/file_handler"
	urlhandler "url-downloader/url_handler"

	"github.com/stretchr/testify/assert"
)

func getwd() string {
	wd, _ := os.Getwd()
	return wd
}

// createTempCSVFile creates a temporary CSV file for testing.
func createTempCSVFile(t *testing.T, urls [][]string) string {

	tmpFile, err := os.CreateTemp(getwd(), "url_list.csv")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	err = writer.WriteAll(urls)
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}
	writer.Flush()
	return tmpFile.Name()
}

func TestFlow(t *testing.T) {
	//Create Temp CSV file
	type result struct {
		resp []byte
	}
	testCases := make(map[string]*result)
	urls := [][]string{
		{"Header"},
		{"https://example.com"},
		{"www.someotherurl.com/api/v1"},
		{"www.anotherone.com"},
	}
	for i := 1; i < len(urls); i++ {
		var b []byte
		url := urls[i][0]
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		b, err = io.ReadAll(res.Body)
		if err != nil {
			continue
		}
		r := &result{
			resp: b,
		}
		testCases[url] = r
	}
	tempFile := createTempCSVFile(t, urls)

	sema := make(chan struct{}, 2)
	output := make(chan *urlhandler.UrlHandler, 1)
	fh := filehandler.NewFileHandler(tempFile, "")
	wg := &sync.WaitGroup{}
	dwg := &sync.WaitGroup{}
	wg.Add(1)
	go fh.ReadCsv(wg, dwg, sema, output)

	for op := range output {
		url := op.GetUrl()
		r, ok := testCases[url]
		assert.Equal(t, ok, true, fmt.Sprintf("url %v not found in test cases", url))
		assert.Equal(t, string(r.resp), string(op.GetData()))
	}
	wg.Wait()
	close(sema)
	err := os.Remove(tempFile)
	assert.NoError(t, err)
}
