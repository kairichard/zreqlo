Goreep
======
... is a simple service which logs `RawQuery` and `Useragent` into redis

Synopsis
  -bind="127.0.0.1:5000": Location server should listen at
  -redis="127.0.0.1:6379": Location of redis instance
  -mount="/": Relative path where handler should be at

## Todo
  * responde with 204 when accept is Text
  * dont read file on every request
  * clean up file structure to be more like the recommended on
