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
	RequestedTagList bool
	TagsList NSXTagList
}

type NSXTagType struct {
	TypeName string `xml:"typeName"`
}

type NSXTag struct {
	ObjectId string `xml:"objectId"`
	ObjectTypeName string `xml:"objectTypeName"`
	VsmUuid string `xml:"vsmUuid"`
	NodeId string `xml:"nodeId"`
	Revision int `xml:"revision"`
	Type NSXTagType `xml:"type"`
	Name string `xml:"name"`
	Description string `xml:"description"`
	IsUniversal bool `xml:"isUniversal"`
	UniversalRevision int `xml:"universalRevision"`
	SystemResource bool `xml:"systemResource"`
	VmCount int `xml:"vmCount"`
}

type NSXTagList struct {
	SecurityTags []NSXTag `xml:"securityTag"`
}

func (c *Config) Client() (*Config, error) {
	if c.InsecureFlag == true {
		resty.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	}
	resty.SetBasicAuth(c.User, c.Password)
	resty.SetHeader("Accept", "application/xml")
	resty.SetTimeout(30 * time.Second)
	return c, nil
}