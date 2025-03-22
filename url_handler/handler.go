package urlhandler

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var once sync.Once

var st *stats

type UrlHandler struct {
	url    string
	data   []byte
	logger *logrus.Logger
}

type stats struct {
	urlsProcessed int32
	urlsFailed    int32
	urlsSucceeded int32
}

func NewStats() {
	once.Do(func() {
		st = &stats{}
	})
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
	defer func() { <-semaphore }()
	atomic.AddInt32(&st.urlsProcessed, 1)
	u.logger.Infof("processing: %v", u.url)
	res, err := http.Get(u.url)
	if err != nil {
		u.logger.Errorf("%v while fetching the url %v", err, u.url)
		atomic.AddInt32(&st.urlsFailed, 1)
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		atomic.AddInt32(&st.urlsFailed, 1)
		u.logger.Errorf("%v while while reading the response body for url  %v", err, u.url)
		return
	}
	atomic.AddInt32(&st.urlsSucceeded, 1)
	u.data = b
	output <- u
	u.logger.Infof("stats : passed %v , failed %v, total %v", atomic.LoadInt32(&st.urlsSucceeded),
		atomic.LoadInt32(&st.urlsFailed),
		atomic.LoadInt32(&st.urlsProcessed))
}

func (u *UrlHandler) GetUrl() string {
	return u.url
}

func (u *UrlHandler) GetData() []byte {
	return u.data
}
