package urlhandler

import (
	"io"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type UrlHandler struct {
	url    string
	data   []byte
	logger *logrus.Logger
}

func NewUrlHandler(u string, log *logrus.Logger) *UrlHandler {
	return &UrlHandler{
		url:    u,
		data:   nil,
		logger: log,
	}
}

func (u *UrlHandler) Download(wg *sync.WaitGroup, semaphore chan struct{}, output chan *UrlHandler) {
	defer wg.Done()
	u.logger.Infof("processing: %v", u.url)
	res, err := http.Get(u.url)
	if err != nil {
		u.logger.Errorf("%v while fetching the url %v", err, u.url)
		return
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		u.logger.Errorf("%v while while reading the response body for url  %v", err, u.url)
		return
	}
	u.data = b
	<-semaphore
	output <- u

}

func (u *UrlHandler) GetUrl() string {
	return u.url
}

func (u *UrlHandler) GetData() []byte {
	return u.data
}
