var assert = require('assert');
var client = require('socket.io-client');
var io = require('socket.io')(8080);
var socket = client.connect('http://localhost:8080');
var sioredis = require('socket.io-redis');

io.adapter(sioredis({
  host : 'localhost',
  port : 6379
}));

socket.on('text', function(data){
  console.log("text", data);
  assert.equal(data, "hogefuga");
});
socket.on('jsondata', function(data){
  console.log("jsondata",data);
  assert.deepEqual(JSON.parse(data), {"name":"a","age":1,"bin":"abc"});
});
socket.on("bin", function(data){
  console.log("bin", "" + data);
  assert.equal("" + data, "aaabbbccc");
});
socket.on("finish", function(){
  console.log("finish");
  process.exit();
});
