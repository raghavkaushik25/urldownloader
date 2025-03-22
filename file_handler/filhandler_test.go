package filehandler

import (
	"os"
	"sync"
	"testing"
	urlhandler "url-downloader/url_handler"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getwd() string {
	wd, _ := os.Getwd()
	return wd
}
func TestWriteDate(t *testing.T) {
	wg := &sync.WaitGroup{}
	output := make(chan *urlhandler.UrlHandler)
	logger := logrus.New()
	folderPath := getwd() + "/" + "test_output"
	urls := []string{
		"https://a.com",
		"https://b.com",
		// invalid url; but file should be created.
		"c",
	}
	fh := NewFileHandler("", "test_output", logger)
	wg.Add(1)
	go fh.WriteData(wg, output)
	for _, url := range urls {
		u := urlhandler.NewUrlHandler(url, logger)
		output <- u
	}
	close(output)
	wg.Wait()
	_, err := os.Stat(folderPath)
	assert.NotEqual(t, true, os.IsNotExist(err))
	assert.NoError(t, err)
	files, err := os.ReadDir(folderPath)
	assert.NoError(t, err)
	assert.Equal(t, len(files), 3)
	err = os.RemoveAll(folderPath)
	assert.NoError(t, err)
}
