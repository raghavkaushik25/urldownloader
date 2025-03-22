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
	// channel to receive the downloaded content of URL. Size of ths channel can be adjusted.
	output = make(chan *urlhandler.UrlHandler, 50)
	// max  # of download goroutines
	maxGoroutines = 50
	// Buffered channel to limit the gororunts
	semaphore = make(chan struct{}, maxGoroutines)
)

func main() {
	// channel to recieve SIGTERM signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// logger with default level as INFO
	logger := logrus.New()
	// separate goroutine to listen to os signals as main will terminate immediately on SIGTERM.
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
	// initialize singleton stats.
	urlhandler.NewStats()
	fh := filehandler.NewFileHandler(*path, "/output", logger)
	// Waitgroup to maintain Read and write go routines.
	wg := &sync.WaitGroup{}
	// Waitgroup to maintain download go routines.
	downloadWg := &sync.WaitGroup{}
	wg.Add(1)
	// One separate goroutine to read csv
	go fh.ReadCsv(wg, downloadWg, semaphore, output)
	wg.Add(1)
	// One separate goroutine to write to disk
	go fh.WriteData(wg, output)
	wg.Wait()
	close(semaphore)

}
