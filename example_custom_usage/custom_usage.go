package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/maxwangnan005/exponent_retry"
)

func Get(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	opt := &exponent_retry.Options{
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		Fn: func(err error) bool {
			s := err.Error()
			if strings.HasPrefix(s, "something") {
				return true // need to retry
			}

			if strings.HasPrefix(s, "duplicate") {
				return false // no need to retry
			}

			return false
		},
	}

	err := exponent_retry.New(opt).Do(context.TODO(), func() error {
		return Get("https://github.com/maxwangnan005/exponent_retry")
	})

	fmt.Println(err)
}
