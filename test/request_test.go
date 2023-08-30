package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	resp, err := http.Get("https://123.sogou.com/")
	assert.Equal(t, err, nil)
	bytes, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil)
	t.Log(string(bytes))
}

func TestWeekTime(t *testing.T) {
	testmap := map[interface{}]int{}
	testmap[2] = 1
	t.Log(testmap[2])
}

func TestStrSplit(t *testing.T) {
	testmap := make(map[int]int)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	//write
	go func() {
		for i := 0; i < 1000; i++ {
			testmap[i] = i
			time.Sleep(time.Millisecond)
		}
		wg.Done()
	}()
	//read
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Millisecond)
			fmt.Println(i, testmap[i])
		}
		wg.Done()
	}()
	wg.Wait()
}
