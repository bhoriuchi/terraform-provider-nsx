package main

import (
	"time"
	"crypto/tls"

	"gopkg.in/resty.v0"
)

type Semver struct {
	Major int
	Minor int
	Patch int
}

type Config struct {
	User string
	Password string
	NSXManager string
	TagEndpoint string
	InsecureFlag bool
	NSXVersion Semver
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