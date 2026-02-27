package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type traceContextKeyring string

const (
	startTimeKey traceContextKeyring = "startTime"
	methodKey    traceContextKeyring = "method"
	pathKey      traceContextKeyring = "path"
)

func TraceTripperMiddleware() (func(*http.Request) error, func(*http.Request, *http.Response, error) error) {
	before := func(req *http.Request) error {
		startTime := time.Now()
		method := req.Method
		path := req.URL.Path

		c := req.Context()
		c = context.WithValue(c, startTimeKey, startTime)
		c = context.WithValue(c, methodKey, method)
		c = context.WithValue(c, pathKey, path)

		*req = *req.WithContext(c)

		return nil
	}
	after := func(req *http.Request, resp *http.Response, err error) error {
		if resp != nil {
			c := req.Context()
			startTime, ok := c.Value(startTimeKey).(time.Time)
			if !ok {
				log.Println("|WARNING| start time is not found in conext")
				startTime = time.Now()
			}
			latency := time.Since(startTime)

			statusCode := resp.StatusCode

			method, _ := c.Value(methodKey).(string)
			path, _ := c.Value(pathKey).(string)

			if err != nil {
				log.Printf("[%s] %s - error: %d - latency: %s", method, path, err, latency)
			} else {
				log.Printf("[%s] %s - satus: %d - latency: %s", method, path, statusCode, latency)
			}
			return nil
		}
		log.Println("request is nil")
		return nil
	}

	return before, after
}
