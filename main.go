package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	filehandler "url-downloader/file_handler"
	urlhandler "url-downloader/url_handler"

	"github.com/sirupsen/logrus"
)

var (
	output        = make(chan *urlhandler.UrlHandler, 50)
	maxGoroutines = 50
	semaphore     = make(chan struct{}, maxGoroutines)
	//logger        = logrus.New()
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	logger := logrus.New()
	go func() {
		<-stop
		ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelTimeout()
		<-ctx.Done()
		logger.Info("Exiting now.")
	}()
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
