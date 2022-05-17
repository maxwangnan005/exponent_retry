package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

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
	err := exponent_retry.Do(context.Background(), func() error {
		return Get("https://github.com/maxwangnan005/exponent_retry")
	})

	fmt.Println(err)
}
