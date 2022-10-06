Custom Transfer Protocol
======

Go implementation on a customized protocol set for file transfer.

## 1. Installation & usage

to install clone this repo:
```sh
git clone git@github.com:danmartinsa/custom-protocol.git
```

To start the server run:
```sh
chmod +x serverApp
./serverApp
```

Alternatively run from source: 
```sh
go run server/server.go
```

To start the client
```sh
chmod +x ./clientApp
./clientApp <server_host>

#Ex.For localhost:
./clientApp 0.0.0.0
```

Alternatively run from source:
```sh
go run client/client.go <server_host>
```



## 2. Server


- Internal Words listm each word is a operation or inforamtion request
``` go
const (
  DIR = "DIR"
  CD  = "CD"
  PWD = "PWD"
  GET = "GET"
  UPL = "UPL"
)
```

- DIR - Return directory list
- CD - Change dir
- PWD - Wcho current dir
- GET - Request for file download
- UPL - Request for file upload


## 3. Client

- internal command list, each command is a direct 
```Go
const (
  uiDir   = "ls"
  uiCd    = "cd"
  uiPwd   = "pwd"
  uiQuit  = "quit"
  uiGet   = "get"
  uiSend  = "send"
  uiMan   = "man"
)
```

- ls - List filders and files 
- cd - Change directories inside the server
- pwd - Echo current server folder structure
- get - request file download from server
- send - reques file upload to server
- man - list of the commands accepted
- quit - close application



