* 基于golang gin框架， websocket协议的 im聊天服务， 支持分布式跨服务器通信（基于redis，pub sub）

* easy chatting server , cross server by redis pubsub
* use gin and redis 
*start redis server ,  redis-server

`go run main.go $port` 
example :  

run server1

`go run main.go 8888`


run server2

`go run main.go 9999`

visit http://localhost:8888 

visit http://localhost:9999 


set uid  ,  "uid1"
set send msg ,  {"to":"uid1","content":"XXX"}

![截图](/server/imgs/screenshoot.png)



