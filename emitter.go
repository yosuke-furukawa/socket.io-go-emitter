package SocketIO

import (
    "github.com/garyburd/redigo/redis"
    "github.com/vmihailenco/msgpack"
    "strconv"
    "bytes"
    _ "fmt"
)

const (
  EVENT = 2
  BINARY_EVENT = 5
)

type EmitterOpts struct {
  // Host means hostname like localhost
  Host string
  // Port means port number, like 6379
  Port int
  // Key means redis subscribe key
  Key string
  // Protocol, like tcp
  Protocol string
  // Address, like localhost:6379
  Addr string
}

type Emitter struct {
  Redis redis.Conn
  Key string
  rooms []string
  flags map[string]interface{}
}

// Emitter constructor
// Usage:
// SocketIO.NewEmitter(&SocketIO.EmitterOpts{
//    Host:"localhost", 
//    Port:6379,
// })
func NewEmitter(opts *EmitterOpts) (*Emitter, error) {
  var addr string
  if opts.Addr != "" {
    addr = opts.Addr
  } else if opts.Host != "" && opts.Port > 0 {
  addr = opts.Host + ":" + strconv.Itoa(opts.Port)
  } else {
    addr = "localhost:6379"
  }
  var protocol string
  if opts.Protocol == "" {
    protocol = "tcp"
  } else {
    protocol = opts.Protocol
  }
  conn, err := redis.Dial(protocol, addr)
  if err != nil {
    return nil, err;
  }
  
  var key string
  if opts.Key == "" {
    key = "socket.io#emitter"
  } else {
    key = opts.Key + "#emitter"
  }

  emitter := &Emitter{
    Redis:conn, 
    Key: key, 
  }
  return emitter, nil;
}

func (emitter *Emitter) Join() *Emitter {
  emitter.flags["join"] = true;
  return emitter;
}

func (emitter *Emitter) Volatile() *Emitter {
  emitter.flags["volatile"] = true;
  return emitter;
}

func (emitter *Emitter) Broadcast() *Emitter {
  emitter.flags["broadcast"] = true;
  return emitter;
}

/**
 * Limit emission to a certain `room`.
 *
 * @param {String} room
 */
func (emitter *Emitter) In(room string) *Emitter {
  for _, r := range emitter.rooms {
    if r == room {
      return emitter
    }
  }
  emitter.rooms = append(emitter.rooms, room)
  return emitter
}

func (emitter *Emitter) To(room string) *Emitter {
  return emitter.In(room)
}


/**
 * Limit emission to certain `namespace`.
 *
 * @param {String} namespace
 */
func (emitter *Emitter) Of(namespace string) *Emitter {
  emitter.flags["nsp"] = namespace
  return emitter
}

// send the packet
// Usage:
// Emit("event name", "data")
func (emitter *Emitter) Emit(event string, data ...interface{}) (*Emitter, error) {
  d := []interface{}{event}
  d = append(d, data...)
  packet := map[string]interface{}{
    "type" : EVENT,
    "data" : d,
  }
  if (emitter.flags["nsp"] != nil) {
    packet["nsp"] = emitter.flags["nsp"]
    delete(emitter.flags, "nsp")
  } else {
    packet["nsp"] = "/"
  }
  var pack []interface{} = make([]interface{}, 0)
  pack = append(pack, packet)
  pack = append(pack, map[string]interface{}{
    "rooms" : emitter.rooms,
    "flags" : emitter.flags,
  })
  buf := &bytes.Buffer{}
	enc := msgpack.NewEncoder(buf)
  error := enc.Encode(pack)
  if error != nil {
    return nil, error
  }
  emitter.Redis.Do("PUBLISH", emitter.Key, buf)
  emitter.rooms = make([]string, 0, 0)
  emitter.flags = make(map[string]interface{})
  return emitter, nil
}

func (emitter *Emitter) Close() {
  if (emitter.Redis != nil) {
    defer emitter.Redis.Close();
  }
}
