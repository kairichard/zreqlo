package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	initStorage(99)
}

func tearDown() {
	client.Cmd("FLUSHDB")
}

func TestHandlerRespondsWithAGif(t *testing.T) {
	defer tearDown()
	q := "a=test&b=test"
	request, _ := http.NewRequest("GET", "/?"+q, nil)
	response := httptest.NewRecorder()

	httpStore(response, request)
	ctype := response.HeaderMap.Get("Content-Type")
	if ctype != "image/gif" {
		t.Fatalf("Content-Type expected %v:\n Got: %v", "image/gif", ctype)
	}
	rb, err := ioutil.ReadAll(response.Body)
	errHndlr(err)
	if bytes.Compare(rb, beacon) != 0 {
		t.Fatalf("beacon was not in body: %v - %v", rb, beacon)
	}
}

func TestHandlerRecordsRequestedPath(t *testing.T) {
	defer tearDown()
	q := "a=test"
	path := "/somewhere/test"
	request, _ := http.NewRequest("GET", path+"?"+q, nil)
	response := httptest.NewRecorder()

	httpStore(response, request)

	res, _ := client.Cmd("LPOP", "incoming").Bytes()
	var ri RequestInfo
	json.Unmarshal(res, &ri)
	if ri.Path != path {
		t.Fatalf("path does not match %+v", ri)
	}
}

func TestHandlerDoesNotRecordAnythingWhenThereIsNoQuery(t *testing.T) {
	defer tearDown()
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	httpStore(response, request)
	length, _ := client.Cmd("LLEN", "incoming").Int()
	if length != 0 {
		t.Fatalf("List should not contain entry")
	}
}

func TestAllQueryParamsAreStored(t *testing.T) {
	defer tearDown()
	q := "a=test&b=test"
	request, _ := http.NewRequest("GET", "/?"+q, nil)
	response := httptest.NewRecorder()

	httpStore(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
	}

	length, _ := client.Cmd("LLEN", "incoming").Int()
	if length != 1 {
		t.Fatalf("List should not have zero entries")
	}

	res, _ := client.Cmd("LPOP", "incoming").Bytes()
	var ri RequestInfo
	json.Unmarshal(res, &ri)
	if ri.Query != q {
		t.Fatalf("Query does not match %+v", ri)
	}
}
