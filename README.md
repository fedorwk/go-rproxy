# go-rproxy
go-rproxy is easy-to-use reverse proxy tool for your services. It can be used both as a standalone application and as a router for your http.Server.

## Usage
1. List your services in following format
`
[service_name] [address] *[protocol]
`
service_name is url path section which defines target service
protocol is optional. If it is not specified, the default protocol will be used, which is specified in srvlist.go DefaultProtocol ("http")
Example srvlist.txt:
```
    static 192.168.1.1:80 http
    api 127.0.0.1:8080 http
```
2. Pass this file to NewServiceListFromReader() 
```go
srvlistFile, _ := os.Open("path/to/file")
srvlist := rproxy.NewServiceListFromReader(srvlistFile)
```
3. Setup Gin router
```go
router := gin.Default()
router, err = rproxy.New(router, srvlist)
```
4. Setup http.Server using router
```go
server := &http.Server{
	Addr: ":8080",
	Handler: router,
}
err := server.ListenAndServe()
```

Now requests sent to proxy.address.com/static/any/path will be redirected to "static" service on 192.168.1.1:80. New path would be **/any/path**.
Same for proxy.address.com/api/. Requests will be redirected to localhost 8080 port. 
