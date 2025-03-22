package urlhandler

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

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
	urlsSucceded  int32
}

func NewStats() {
	once.Do(func() {
		st = &stats{}
	})
}

func GetStats() (int32, int32, int32) {
	return atomic.LoadInt32(&st.urlsProcessed), atomic.LoadInt32(&st.urlsSucceded), atomic.LoadInt32(&st.urlsFailed)
}

func NewUrlHandler(u string, log *logrus.Logger) *UrlHandler {
	return &UrlHandler{
		url:    u,
		data:   nil,
		logger: log,
	}
}

func (u *UrlHandler) Download(wg *sync.WaitGroup, semaphore chan struct{}, output chan *UrlHandler) {
	start := time.Now()
	defer wg.Done()
	defer func() { <-semaphore }()
	atomic.AddInt32(&st.urlsProcessed, 1)
	u.logger.Debugf("processing: %v", u.url)
	res, err := http.Get(u.url)
	if err != nil {
		atomic.AddInt32(&st.urlsFailed, 1)
		u.logger.Errorf("err : %v while fetching the url %v; failed count %v", err, u.url, atomic.LoadInt32(&st.urlsFailed))
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		atomic.AddInt32(&st.urlsFailed, 1)
		u.logger.Errorf("err : %v while while reading the response body for url %v; failed count %v  ", err, u.url, atomic.LoadInt32(&st.urlsFailed))
		return
	}
	atomic.AddInt32(&st.urlsSucceded, 1)
	u.data = b
	output <- u
	u.logger.Debugf("stats : time taken to process url %v is %v  ", u.url, time.Since(start))
}

func (u *UrlHandler) GetUrl() string {
	return u.url
}

func (u *UrlHandler) GetData() []byte {
	return u.data
}
