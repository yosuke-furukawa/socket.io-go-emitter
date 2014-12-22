package SocketIO

import (
  "time"
	"bytes"
	"github.com/garyburd/redigo/redis"
	"strings"
	"testing"
)

func TestHasBinary(t *testing.T) {
  hasBin := HasBinary("string")
  if hasBin != false {
		t.Error("string is not binary")
  }
  hasBin = HasBinary(123)
  if hasBin != false {
		t.Error("integer is not binary")
  }
  hasBin = HasBinary([]byte("abc"))
  if hasBin != true {
		t.Error("[]byte is binary")
  }
  var data []interface{}
  data = append(data, 1)
  data = append(data, "2")
  data = append(data, 3)
  hasBin = HasBinary(data)
  if hasBin != false {
		t.Error("string array is not binary")
  }
  data = append(data, []byte("aaa"))
  hasBin = HasBinary(data)
  if hasBin != true {
		t.Error("data has binary")
  }
  dataMap := make(map[string]interface{})
  dataMap["hoge"] = "bbb"
  hasBin = HasBinary(dataMap)
  if hasBin != false {
		t.Error("data doesnot have binary")
  }
  dataMap["fuga"] = []byte("bbb")
  hasBin = HasBinary(dataMap)
  if hasBin != true {
		t.Error("data has binary")
  }
  dataMap["fuga"] = data
  hasBin = HasBinary(dataMap)
  if hasBin != true {
		t.Error("data has binary")
  }

}

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
	emitter.Emit("text", "hogefuga")
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

func TestPublishMultipleTimes(t *testing.T) {
	time.Sleep(1 * time.Second)
	emitter, _ := NewEmitter(&EmitterOpts{
		Host: "localhost",
		Port: 6379,
	})

	if emitter == nil {
		t.Error("emitter is nil")
	}
	defer emitter.Close()
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("socket.io#emitter")
	emitter.Emit("text", "hogefuga")
	emitter.Emit("text", "foobar")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "foobar")
			isFirst := strings.Contains(string(v.Data), "hogefuga")
			if !isContain {
				if !isFirst {
					t.Errorf("%s not contains foobar", v.Data)
				}
			} else {
				return
			}
		}
	}
}

func TestPublishJson(t *testing.T) {
  time.Sleep(1 * time.Second)
	emitter, _ := NewEmitter(&EmitterOpts{
		Host: "localhost",
		Port: 6379,
	})

	if emitter == nil {
		t.Error("emitter is nil")
	}
	defer emitter.Close()
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("socket.io#emitter")
	emitter.Emit("jsondata", []byte(`{"name":"a","age":1,"bin":"abc"}`))
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "abc")
			if !isContain {
				t.Errorf("%s not contains abc", v.Data)
				return
			} else {
				return
			}
		}
	}
}

func TestPublishBinary(t *testing.T) {
  time.Sleep(1 * time.Second)
	emitter, _ := NewEmitter(&EmitterOpts{
		Host: "localhost",
		Port: 6379,
	})

	if emitter == nil {
		t.Error("emitter is nil")
	}
	defer emitter.Close()
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("socket.io#emitter")
	val := bytes.NewBufferString("aaabbbccc")
	emitter.EmitBinary("bin", val.Bytes())
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

func TestPublishEnd(t *testing.T) {
	time.Sleep(1 * time.Second)
	emitter, _ := NewEmitter(&EmitterOpts{
		Host: "localhost",
		Port: 6379,
	})
	defer emitter.Close()
	c, _ := redis.Dial("tcp", "localhost:6379")
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("socket.io#emitter")
	emitter.Emit("finish")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			isContain := strings.Contains(string(v.Data), "finish")
			if !isContain {
				t.Errorf("%s not contains end", v.Data)
				return
			} else {
				return
			}
		}
	}
}
