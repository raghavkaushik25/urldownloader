package urlhandler

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestUrlHandler(t *testing.T) {
	wg := &sync.WaitGroup{}
	consumerWg := &sync.WaitGroup{}
	sema := make(chan struct{}, 1)
	output := make(chan *UrlHandler, 1)
	result := []*UrlHandler{}
	logger := logrus.New()
	urls := []string{
		"https://example.com",
		//Invalid URL
		"1234",
		"https://www.webpagetest.org/",
	}
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for op := range output {
			result = append(result, op)
		}
	}()
	for _, url := range urls {
		u := NewUrlHandler(url, logger)
		wg.Add(1)
		sema <- struct{}{}
		go u.Download(wg, sema, output)
	}
	wg.Wait()
	close(output)
	close(sema)
	consumerWg.Wait()
	assert.Equal(t, len(result), 2)
}
