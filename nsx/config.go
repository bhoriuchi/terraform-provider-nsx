package nsx

import (
	"crypto/tls"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/go-resty/resty"
)

type Config struct {
	User         string
	Password     string
	NSXManager   string
	TagEndpoint  string
	InsecureFlag bool
	NSXVersion   semver.Version
}

func (c *Config) ClientInit() (*Config, error) {
	if c.InsecureFlag == true {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	resty.SetBasicAuth(c.User, c.Password)
	resty.SetHeader("Accept", "application/xml")
	resty.SetHeader("Content-Type", "application/xml")
	resty.SetTimeout(30 * time.Second)
	return c, nil
}
