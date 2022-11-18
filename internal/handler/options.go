package handler

import "os"

type HttpServerOpt struct {
	ListenAddr string
}

func ParseHTTPServerOptions() (*HttpServerOpt, error) {
	opt := HttpServerOpt{
		ListenAddr: "127.0.0.1:8181",
	}
	if addr, ok := os.LookupEnv("HTTP_LISTEN_ADDR"); ok {
		opt.ListenAddr = addr
	}

	return &opt, nil
}
