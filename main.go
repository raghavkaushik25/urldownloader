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
	logger := logrus.New()
	path := flag.String("path", "", "Path of the csv; -path=/path/to/csv")
	debug := flag.Bool("debug", false, "Run application in debug mode; -debug")
	flag.Parse()
	if *path == "" {
		logger.Fatal("path to the csv must be provided")
	}
	if *debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	urlhandler.NewStats()
	fh := filehandler.NewFileHandler(*path, "/output", logger)
	wg := &sync.WaitGroup{}
	downloadWg := &sync.WaitGroup{}
	wg.Add(1)
	go fh.ReadCsv(wg, downloadWg, semaphore, output)
	wg.Add(1)
	go fh.WriteData(wg, output)
	wg.Wait()
	close(semaphore)
}
