package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "encoding/json"
)

func init() {
  initStorage(99)
}

func tearDown(){
  client.Cmd("FLUSHDB")
}

func TestAllQueryParamsAreStored(t *testing.T) {
    q := "a=test&b=test"
    request, _ := http.NewRequest("GET", "/?" + q, nil)
    response := httptest.NewRecorder()

    httpStore(response, request)

    if response.Code != http.StatusOK {
        t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
    }

    defer tearDown()
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
