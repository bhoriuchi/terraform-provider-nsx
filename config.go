package main

import (
	"gopkg.in/resty.v0"
	"crypto/tls"
)

type Config struct {
	User string
	Password string
	NSXManager string
	InsecureFlag bool
}

func (c *Config) Client() (interface{}, error) {
	if c.InsecureFlag == true {
		resty.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	}
	resty.SetBasicAuth(c.User, c.Password)
	resty.SetHeader("Accept", "application/xml")

	client := resty.R()
	return &client, nil
}