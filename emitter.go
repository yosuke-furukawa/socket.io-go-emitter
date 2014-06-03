package SocketIO

import (
    "gopkg.in/redis.v1"
    "github.com/vmihailenco/msgpack"
    "strconv"
)

const (
  EVENT = 2
  BINARY_EVENT = 5
)

type EmitterOpts struct {
  Host string
  Port int
  Key string
  Socket string
  Password string
  DB int64
}

type Emitter struct {
  Redis *redis.Client
  Key string
  rooms []string
  flags map[string]interface{}
}

func (emitter *Emitter) initializeFlags() {
  emitter.flags = map[string]interface{}{
    "json" : true,
    "volatile" : true,
    "broadcast" : true,
  }
}

func NewEmitter(opts *EmitterOpts) *Emitter {
  var client *redis.Client
  if opts.DB == 0 {
    opts.DB = int64(-1)
  }
  if opts.Socket != "" {
    client = redis.NewUnixClient(&redis.Options{
    Addr: opts.Socket, 
    Password: opts.Password, 
    DB : opts.DB,
    })
  } else {
    addr := opts.Host + ":" + strconv.Itoa(opts.Port)
    client = redis.NewTCPClient(&redis.Options{
    Addr: addr, 
    Password: opts.Password, 
    DB: opts.DB,
    })
  }
  
  var key string
  if opts.Key == "" {
    key = "socket.io#emitter"
  } else {
    key = opts.Key + "#emitter"
  }

  emitter := &Emitter{
    Redis:client, 
    Key: key, 
  }
  emitter.initializeFlags()
  return emitter
}

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

func (emitter *Emitter) Of(namespace string) *Emitter {
  emitter.flags["nsp"] = namespace
  return emitter
}

func (emitter *Emitter) Emit(data ...string) *Emitter {
  packet := map[string]interface{}{
    "type" : EVENT,
    "data" : data,
  }
  if (emitter.flags["nsp"] != nil) {
    packet["nsp"] = emitter.flags["nsp"]
  } else {
    packet["nsp"] = "/"
  }
  var pack []interface{} = make([]interface{}, 3)
  pack = append(pack, packet)
  pack = append(pack, emitter.rooms)
  pack = append(pack, emitter.flags)
  message, _ := msgpack.Marshal(pack)
  emitter.Redis.Publish(emitter.Key, string(message))
  emitter.rooms = make([]string, 5, 5)
  emitter.flags = make(map[string]interface{})
  return emitter
}

func (emitter *Emitter) Close() {
  if (emitter.Redis != nil) {
    defer emitter.Redis.Close();
  }
}

//  func (emitter *Emitter) EmitBinary(buf []byte) *Emitter {
//    packet := map[string]interface{}{
//      "type" : BINARY_EVENT,
//      "data" : buf,
//    } 
//    if (emitter.flags["nsp"] != nil) {
//      packet["nsp"] = emitter.flags["nsp"]
//    } else {
//      packet["nsp"] = "/"
//    }
//    var pack []interface{} = make([]interface{}, 3)
//    pack = append(pack, packet)
//    pack = append(pack, emitter.rooms)
//    pack = append(pack, emitter.flags)
//    message, _ := msgpack.Marshal(pack)
//    emitter.Redis.Publish(emitter.Key, string(message))
//    emitter.rooms = make([]string)
//    emitter.flags = make(map[string]interface{})
//    return emitter
//  }
