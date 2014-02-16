package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fzzy/radix/redis"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	redis_location = flag.String("redis", "127.0.0.1:6379", "Location of redis instance")
	server_bind    = flag.String("bind", "127.0.0.1:5000", "Location server should listen at")
	handler_mount  = flag.String("mount", "/", "Relative path where handler should be at")
)

var client *redis.Client
var beacon []byte

type RequestInfo struct {
	Query     string
	Path      string
	UserAgent string
	Time      int64
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	var e error
	beacon, e = ioutil.ReadFile(path.Join(path.Dir(filename), "assets/1x1.gif"))
	errHndlr(e)
}

func initStorage(db int) {
	var err error
	client, err = redis.DialTimeout("tcp", *redis_location, time.Duration(10)*time.Second)
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
		Query:     req.URL.RawQuery,
		Path:      req.URL.Path,
		UserAgent: req.UserAgent(),
		Time:      time.Now().Unix(),
	}

	if ri.Query != "" {
		b, err := json.Marshal(ri)
		errHndlr(err)
		client.Cmd("RPUSH", "incoming", string(b))
	}

	res.Header().Set("Content-Type", "image/gif")
	res.Write(beacon)
}

func main() {
	flag.Parse()
	initStorage(1)
	http.HandleFunc(*handler_mount, httpStore)
	err := http.ListenAndServe(*server_bind, nil)
	errHndlr(err)
	defer client.Close()
}
