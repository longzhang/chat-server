* easy chatting server , cross server by redis pubsub
* use gin and redis 

start redis server ,  redis-server
go run main.go $port 
example :  
run server1
go run main.go 8888
run server2
go run main.go 8080

visit http://localhost:0000 

set uid 
set send msg ,  {"to":"uid","content":"XXX"}

