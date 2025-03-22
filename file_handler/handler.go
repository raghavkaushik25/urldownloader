package filehandler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"sync"
	"time"
	urlhandler "url-downloader/url_handler"

	"github.com/sirupsen/logrus"
)

type FileHandler interface {
	ReadCsv(wg *sync.WaitGroup, downloadWg *sync.WaitGroup, sema chan struct{}, output chan *urlhandler.UrlHandler)
	WriteData(wg *sync.WaitGroup, output chan *urlhandler.UrlHandler)
}

type fileHanlder struct {
	logger *logrus.Logger
	path   string
}

func NewFileHandler(p string) FileHandler {
	return &fileHanlder{
		path:   p,
		logger: logrus.New(),
	}
}

func (fh *fileHanlder) ReadCsv(wg *sync.WaitGroup, downloadWg *sync.WaitGroup, sema chan struct{}, output chan *urlhandler.UrlHandler) {
	defer wg.Done()
	f, err := os.Open(fh.path)
	if err != nil {
		fh.logger.Fatalf("error %v reading csv at path %v", err, fh.path)
	}
	r := csv.NewReader(f)
	//Ignore the header
	_, err = r.Read()
	if err != nil {
		fh.logger.Fatalf("error %v reading csv at path %v", err, fh.path)
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			fmt.Printf("reached EOF\n")
			break
		}
		if err != nil {
			fmt.Printf("error %v reading the record %v", err, record)
			continue
		}
		url := record[0]
		fh.logger.Infof("got new url %v ", url)
		u := urlhandler.NewUrlHandler(url, fh.logger)
		downloadWg.Add(1)
		sema <- struct{}{}
		go u.Download(downloadWg, sema, output)
	}
	downloadWg.Wait()
	close(output)
}

func (fh *fileHanlder) WriteData(wg *sync.WaitGroup, output chan *urlhandler.UrlHandler) {
	defer wg.Done()
	for u := range output {
		fPrefix := ""
		fh.logger.Infof("writing file for url %v", u.GetUrl())
		url, err := url.Parse(u.GetUrl())
		if err != nil {
			fPrefix = "dummy"
		} else {
			fPrefix = url.Host
		}
		f, err := os.Create(fmt.Sprintf("/Users/raghav/go/src/url_downloader/output/%v_%v.txt", fPrefix, time.Now().String()))
		defer f.Close()
		if err != nil {
			fh.logger.Errorf("%v while creating file for url %v", err, u.GetUrl())
			continue
		}
		_, err = f.Write(u.GetData())
		if err != nil {
			fh.logger.Errorf("%v while writing file for url %v", err, u.GetUrl())
			continue
		}
	}
}
