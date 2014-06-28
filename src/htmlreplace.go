package main

import (
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	exphtml "code.google.com/p/go.net/html"

	"fmt"
	"net/http"
)

func vuvu(a string) string {
	return "http://vuvuzelr.7co.cc/vuvuzela.jpg"
}

func ReplaceAndAppend(url string, img string, swf string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching page: " + url)
		return "", err
	}
	tree, err := transform.NewFromReader(resp.Body)

	if err != nil {
		fmt.Println("Error creating Doc")
		return "", err
	}

	attr := []exphtml.Attribute{}
	attr = append(attr, exphtml.Attribute{Key: "data", Val: swf})
	attr = append(attr, exphtml.Attribute{Key: "width", Val: "1"})
	attr = append(attr, exphtml.Attribute{Key: "height", Val: "1"})

	object := h5.Element("object", attr)

	t := tree.Clone()
	tt, _ := transform.Trans(transform.TransformAttrib("src", func(string) string { return img }), "img")
	t.Apply(transform.AppendChildren(object), "body")

	t.ApplyAll(tt)
	return t.String(), nil
}
