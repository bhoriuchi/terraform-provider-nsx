package main

import (
	"encoding/xml"
)

type NSXTagType struct {
	TypeName string `xml:"typeName"`
}

type NSXTag struct {
	XMLName 	xml.Name 	`xml:"securityTag"`
	ObjectId 	string 		`xml:"objectId,omitempty"`
	ObjectTypeName 	string 		`xml:"objectTypeName,omitempty"`
	Type 		NSXTagType 	`xml:"type,omitempty"`
	Name 		string 		`xml:"name,omitempty"`
	Description 	string 		`xml:"description,omitempty"`
	IsUniversal 	bool 		`xml:"isUniversal,omitempty"`
	VmCount		int		`xml:"vmCount,omitempty"`
}

type NSXTagList struct {
	SecurityTags []NSXTag `xml:"securityTag"`
}