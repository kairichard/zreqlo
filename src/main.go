package main

import (
  "os"
  "fmt"
  "time"
  "log"
  "path"
  "runtime"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "github.com/fzzy/radix/redis"
)


var client *redis.Client

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
  _, filename, _, _ := runtime.Caller(0)
  beacon, err := ioutil.ReadFile(path.Join(path.Dir(filename), "../assets/1x1.gif"))
	if err != nil {
		log.Fatal(err)
	}
  res.Header().Set("Content-Type", "image/gif")
  res.Write(beacon)
}

func main(){
  initStorage(1)
  http.HandleFunc("/", httpStore)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
  defer client.Close()
}
