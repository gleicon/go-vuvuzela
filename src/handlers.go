// Copyright 2014 go-vuvuzela authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/fiorix/go-redis/redis"
	"net/http"
	"regexp"
)

var urlRe = regexp.MustCompile("^([a-zA-Z0-9]+)$")
var urlMatchRe = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)

func (s *httpServer) route() {
	http.HandleFunc("/u/", s.sendPageHandler)
	http.HandleFunc("/api/proxy", s.proxyHandler)
	http.HandleFunc("/api/stats/", s.statsHandler)
	http.Handle("/", http.FileServer(http.Dir(s.config.DocumentRoot)))
}

func cacheAndReplace(rc *redis.Client, eid string, url string, config *configFile) string {
	content, err := rc.Get("proxy:cache:" + eid)
	if content == "" || err != nil {
		var contents string
		if url == "" {
			url, _ = checkEid(rc, url)
		}
		contents, _ = ReplaceAndAppend(url, config.VUVUZELA.Img, config.VUVUZELA.Swf)
		_ = rc.Set("proxy:cache:"+eid, contents)
		_, _ = rc.Expire("proxy:cache:"+eid, 60000)
		return contents
	}
	return content
}

func (s *httpServer) sendPageHandler(w http.ResponseWriter, r *http.Request) {
	qn := r.URL.Path[len("/u/"):]
	if !urlRe.MatchString(qn) {
		http.Error(w, "Invalid url id: "+qn, 400)
		return
	}
	url, _ := checkEid(s.redis, qn)

	if url == "" {
		http.Error(w, "Id not found: "+qn, 404)
		return
	}
	_, _ = s.redis.HIncrBy("proxy:url:"+qn, "clicks", 1)
	fmt.Fprintf(w, cacheAndReplace(s.redis, qn, url, s.config))
}

func (s *httpServer) statsHandler(w http.ResponseWriter, r *http.Request) {
	qn := r.URL.Path[len("/api/stats/"):]
	if !urlRe.MatchString(qn) {
		http.Error(w, "Invalid url id: "+qn, 400)
		return
	}
	stats, _ := s.redis.HGetAll("proxy:url:" + qn)
	t, _ := json.Marshal(stats)
	fmt.Fprintf(w, string(t))
}

func (s *httpServer) proxyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		uuid int
		eid  string
	)

	url := r.PostFormValue("url")

	if !urlMatchRe.MatchString(url) {
		http.Error(w, "Invalid url format: "+url, 400)
		return
	}

	eid, _ = checkUrl(s.redis, url)

	if eid == "" {
		uuid = getUUID(s.redis)
		eid = base62FromUUID(uuid)
		_ = updateLookupTable(s.redis, url, eid)
		_ = s.redis.HSet("proxy:url:"+eid, "uuid", string(uuid))
		_ = s.redis.HSet("proxy:url:"+eid, "url", url)
		_ = s.redis.HSet("proxy:url:"+eid, "clicks", "0")
		_ = cacheAndReplace(s.redis, eid, url, s.config)
	}

	http.Redirect(w, r, "/u/"+eid, 302)
}

func updateLookupTable(rc *redis.Client, url string, eid string) error {
	err := rc.HSet("proxy:url:lookup", url, eid)
	if err != nil {
		return err
	}

	err = rc.HSet("proxy:eid:lookup", eid, url)
	if err != nil {
		return err
	}
	return nil
}

func checkUrl(rc *redis.Client, url string) (string, error) {
	return rc.HGet("proxy:url:lookup", url)
}

func checkEid(rc *redis.Client, eid string) (string, error) {
	return rc.HGet("proxy:eid:lookup", eid)
}
