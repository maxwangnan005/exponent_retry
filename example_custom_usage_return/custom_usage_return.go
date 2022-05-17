package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/maxwangnan005/exponent_retry"
)

func Get() (int64, error) {
	rand.Seed(time.Now().UnixNano())
	number := rand.Int63n(4)

	if number == 1 {
		return number, errors.New("something is wrong")
	}

	if number == 2 {
		return number, errors.New("something is wrong")
	}

	if number == 3 {
		return number, errors.New("duplicate")
	}

	return 0, nil
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

	body, err := exponent_retry.New(opt).DoWithReturn(context.TODO(), func() (interface{}, error) {
		return Get()
	})

	if v, ok := body.(int64); ok {
		fmt.Println(v)
	} else {
		fmt.Println(body)
	}

	fmt.Println(err)
}
