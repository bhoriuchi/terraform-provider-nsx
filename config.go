package main

import (
	"gopkg.in/resty.v0"
	"crypto/tls"
	"time"
)

type Config struct {
	User string
	Password string
	NSXManager string
	InsecureFlag bool
}

func (c *Config) Client() (*Config, error) {
	if c.InsecureFlag == true {
		resty.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	}
	resty.SetBasicAuth(c.User, c.Password)
	resty.SetHeader("Accept", "application/xml")
	resty.SetHeader("Content-Type", "application/xml")
	resty.SetTimeout(30 * time.Second)
	return c, nil
}