package main

import (
	"flag"
	"sync"
	filehandler "url-downloader/file_handler"
	urlhandler "url-downloader/url_handler"

	"github.com/sirupsen/logrus"
)

var (
	output        = make(chan *urlhandler.UrlHandler, 50)
	maxGoroutines = 50
	semaphore     = make(chan struct{}, maxGoroutines)
	logger        = logrus.New()
)

func main() {

	path := flag.String("path", "", "-path=/path/to/csv")
	flag.Parse()
	if *path == "" {
		logger.Fatal("path to the csv must be provided")
	}
	urlhandler.NewStats()
	fh := filehandler.NewFileHandler(*path, "/output")
	wg := &sync.WaitGroup{}
	downloadWg := &sync.WaitGroup{}
	wg.Add(1)
	go fh.ReadCsv(wg, downloadWg, semaphore, output)
	wg.Add(1)
	go fh.WriteData(wg, output)
	wg.Wait()
	close(semaphore)
}
