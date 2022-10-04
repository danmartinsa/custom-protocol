package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
  uiDir   = "dir"
  uiCd    = "cd"
  uiPwd   = "pwd"
  uiQuit  = "quit"
  uiGet   = "get"
)

const (
  DIR = "DIR"
  CD  = "CD"
  PWD = "PWD"
  GET = "GET"
)

const CACHE = 512 

func ConnHandler() {
  if len(os.Args) != 2 {
    fmt.Println("Usage: ", os.Args[0], "host")
    os.Exit(1)
  }

  host := os.Args[1]

  conn, err := net.Dial("tcp", host+":1202")
  checkError(err)

  reader := bufio.NewReader(os.Stdin)
  buffer := make([]byte, CACHE)
  for {
    line, err := reader.ReadString('\n')
    line = strings.TrimRight(line, "\t\r\n") 
    if err != nil {
      break
    }
    strs := strings.SplitN(line, " ", 2)

    switch strs[0] {
    case uiDir:
      dirRequests(conn)

    case uiCd:
      if len(strs) != 2 {
        fmt.Println("cd <dir>")
        continue
      }
      fmt.Println("CD \"", strs[1], "\"")
      cdRequest(conn, strings.Trim(strs[1], "\t\r\n\f\v "))

    case uiPwd:
      pwdRequest(conn)

    case uiGet:
      if len(strs) == 1 {
        fmt.Println("Please Provide a filename")
        continue
      }
      conn.Write([]byte(GET + " " + strs[1]))
      n, _ := conn.Read(buffer)
      fileSize, err := strconv.ParseInt(string(buffer[:n]), 10, 64)
      if err != nil {
        fmt.Println("ERROR: ", string(buffer[:n]))
        continue
      }
      GetFileRequest(conn, strs[1], fileSize)

    case uiQuit:
      conn.Close()
      os.Exit(0)
    default:
      fmt.Println("Unknow command")
  }
  }
}

func dirRequests(conn net.Conn) {
  conn.Write([]byte(DIR + " "))

  var buf [CACHE]byte
  result := bytes.NewBuffer(nil)
  for {
    n, _ := conn.Read(buf[0:])
    result.Write(buf[0:n])
    length := result.Len()
    contents := result.Bytes()
    if string(contents[length-4:]) == "\r\n\r\n" {
      fmt.Println(string(contents[0:length-4]))
      return
    }
  }
}

func cdRequest(conn net.Conn, dir string) {
  conn.Write([]byte(CD + " " + dir))
  var response [CACHE]byte
  n, _ := conn.Read(response[0:])
  s := string(response[0:n])
  if s != "OK" {
    fmt.Println("Failed to change dir")
  }
}

func pwdRequest(conn net.Conn) {
  conn.Write([]byte(PWD))
  var response [CACHE]byte
  n, _ := conn.Read(response[0:])
  s := string(response[0:n])
  fmt.Println("Current dir \"" + s + "\"")
}

func GetFileRequest(conn net.Conn, fileName string, fileSize int64) {
  newFile, err := os.Create("." + "/" + fileName)
  var fileSizeDownload int64

  if err != nil {
    fmt.Println(err)
  } 
  defer newFile.Close()

  for {
    if (fileSize - fileSizeDownload) < CACHE {
      io.CopyN(newFile, conn, (fileSize - fileSizeDownload)) 
      break
    }
    io.CopyN(newFile, conn, CACHE)
    fileSizeDownload += CACHE
  }
  fmt.Println("OK")
}

func checkError(err error) {
  if err != nil {
    fmt.Println("Fatal error", err.Error())
    os.Exit(1)
  }
}ge main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
  uiDir   = "dir"
  uiCd    = "cd"
  uiPwd   = "pwd"
  uiQuit  = "quit"
  uiGet   = "get"
)

const (
  DIR = "DIR"
  CD  = "CD"
  PWD = "PWD"
  GET = "GET"
)

const CACHE = 4096

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Usage: ", os.Args[0], "host")
    os.Exit(1)
  }

  host := os.Args[1]

  conn, err := net.Dial("tcp", host+":1202")
  checkError(err)

  reader := bufio.NewReader(os.Stdin)
  buffer := make([]byte, CACHE)
  for {
    line, err := reader.ReadString('\n')
    line = strings.TrimRight(line, "\t\r\n") 
    if err != nil {
      break
    }
    strs := strings.SplitN(line, " ", 2)

    switch strs[0] {
    case uiDir:
      dirRequests(conn)

    case uiCd:
      if len(strs) != 2 {
        fmt.Println("cd <dir>")
        continue
      }
      fmt.Println("CD \"", strs[1], "\"")
      cdRequest(conn, strings.Trim(strs[1], "\t\r\n\f\v "))

    case uiPwd:
      pwdRequest(conn)

    case uiGet:
      if len(strs) == 1 {
        fmt.Println("Please Provide a filename")
        continue
      }
      conn.Write([]byte(GET + " " + strs[1]))
      n, _ := conn.Read(buffer)
      fileSize, err := strconv.ParseInt(string(buffer[:n]), 10, 64)
      if err != nil {
        fmt.Println("ERROR: ", string(buffer[:n]))
        continue
      }
      GetFileRequest(conn, strs[1], fileSize)

    case uiQuit:
      conn.Close()
      os.Exit(0)
    default:
      fmt.Println("Unknow command")
  }
  }
}

func dirRequests(conn net.Conn) {
  conn.Write([]byte(DIR + " "))

  var buf [CACHE]byte
  result := bytes.NewBuffer(nil)
  for {
    n, _ := conn.Read(buf[0:])
    result.Write(buf[0:n])
    length := result.Len()
    contents := result.Bytes()
    if string(contents[length-4:]) == "\r\n\r\n" {
      fmt.Println(string(contents[0:length-4]))
      return
    }
  }
}

func cdRequest(conn net.Conn, dir string) {
  conn.Write([]byte(CD + " " + dir))
  var response [CACHE]byte
  n, _ := conn.Read(response[0:])
  s := string(response[0:n])
  if s != "OK" {
    fmt.Println("Failed to change dir")
  }
}

func pwdRequest(conn net.Conn) {
  conn.Write([]byte(PWD))
  var response [CACHE]byte
  n, _ := conn.Read(response[0:])
  s := string(response[0:n])
  fmt.Println("Current dir \"" + s + "\"")
}

func GetFileRequest(conn net.Conn, fileName string, fileSize int64) {
  newFile, err := os.Create("." + "/" + fileName)
  var fileSizeDownload int64

  if err != nil {
    fmt.Println(err)
  } 
  defer newFile.Close()

  for {
    if (fileSize - fileSizeDownload) < CACHE {
      io.CopyN(newFile, conn, (fileSize - fileSizeDownload)) 
      break
    }
    io.CopyN(newFile, conn, CACHE)
    fileSizeDownload += CACHE
  }
  fmt.Println("OK")
}

func checkError(err error) {
  if err != nil {
    fmt.Println("Fatal error", err.Error())
    os.Exit(1)
  }
}

