Request Logger - [![Build Status](https://travis-ci.org/kairichard/zreqlo.png?branch=master)](https://travis-ci.org/kairichard/request_logger)
======
... is a simple service which logs `RawQuery` and `Useragent` into redis
```
Synopsis
  -bind="127.0.0.1:5000": Location server should listen at
  -redis="127.0.0.1:6379": Location of redis instance
  -mount="/": Relative path where handler should be at
```
## Todo

  * responde with 204 when accept is Text
  * rewrite tests to use ginko or testify
