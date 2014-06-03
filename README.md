socket.io-go-emitter
=====================

A Golang implementation of [socket.io-emitter](https://github.com/Automattic/socket.io-emitter)

This project uses redis.
Make sure your environment has redis.

Install and development
--------------------

To install in your golang project.

```sh
$ go get github.com/yosuke-furukawa/socket.io-go-emitter
```

Usage
---------------------

Example:

```go
  emitter := SocketIO.NewEmitter(&SocketIO.EmitterOpts{
    Host:"localhost", 
    Port:6379,
  })
  emitter.Emit("message", "I love you!!")
```

### Broadcasting and other flags

Possible flags

- json
- volatile
- broadcast

```go
  emitter := SocketIO.NewEmitter(&SocketIO.EmitterOpts{
    Host:"localhost", 
    Port:6379,
  })
  emitter.Volatile().Emit("message", "I love you!!")
```


