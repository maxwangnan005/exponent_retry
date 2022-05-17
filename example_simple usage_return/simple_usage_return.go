package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/maxwangnan005/exponent_retry"
)

func Get() (int64, error) {
	rand.Seed(time.Now().UnixNano())
	number := rand.Int63n(28)

	if number != 0 {
		return number, errors.New("something is wrong")
	}

	return 0, nil
}

func main() {
	body, err := exponent_retry.DoWithReturn(context.TODO(), func() (interface{}, error) {
		return Get()
	})

	if v, ok := body.(int64); ok {
		fmt.Println(v)
	} else {
		fmt.Println(body)
	}

	fmt.Println(err)
}
