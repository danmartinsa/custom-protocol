package ctp

import "net"


type Conn struct {
  conn      net.Conn
  dataType  dataType
}
