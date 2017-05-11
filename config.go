package main

import (
	"log"
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

	resp, err := client.Get(c.NSXManager + "/2.0/services/securitytags/tag/securitytag-9")
	log.Printf("!!! %+v", resp)

	return &client, err
}