package main

import (
  "os"
  "fmt"
  "time"
  //"io/ioutil"
  "net/http"
  "encoding/json"
  "github.com/fzzy/radix/redis"
)


var client *redis.Client
// var image, err := ioutil.ReadFile("assets/1x1.gif")

type RequestInfo struct {
  Query string
  UserAgent string
  Time int64
}

func initStorage(db int){
  var err error
  client, err = redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
  errHndlr(err)
  client.Cmd("SELECT", db)
}

func errHndlr(err error) {
  if err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }
}

func httpStore(res http.ResponseWriter, req *http.Request) {
  ri := RequestInfo{
    Query: req.URL.RawQuery,
    UserAgent: req.UserAgent(),
    Time: time.Now().Unix(),
  }
  b, err := json.Marshal(ri)
  if err != nil {
    fmt.Print(err)
  }
  client.Cmd("RPUSH", "incoming", string(b))
}

func main(){
  initStorage(1)
  defer client.Close()
}
