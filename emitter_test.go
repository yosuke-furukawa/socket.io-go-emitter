package SocketIO

import (
	"bytes"
	_ "fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
	"testing"
)

func TestPublish(t *testing.T) {
	emitter, _ := NewEmitter(&EmitterOpts{
		Host: "localhost",
		Port: 6379,
	})

	if emitter == nil {
		t.Error("emitter is nil")
	}
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("socket.io#emitter")
	emitter.Emit("time", "hogefuga")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "hogefuga")
			if !isContain {
				t.Errorf("%s not contains hogefuga", v.Data)
				return
			} else {
				return
			}
		}
	}
}

func TestPublishBinary(t *testing.T) {
	emitter, _ := NewEmitter(&EmitterOpts{
		Host: "localhost",
		Port: 6379,
	})

	if emitter == nil {
		t.Error("emitter is nil")
	}
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("socket.io#emitter")
	val := bytes.NewBufferString("aaabbbccc")
	emitter.EmitBinary("time", val.Bytes())
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "aaabbbccc")
			if !isContain {
				t.Errorf("%s not contains aaabbbccc", v.Data)
				return
			} else {
				return
			}
		}
	}
}
