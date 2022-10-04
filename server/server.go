package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
  DIR = "DIR"
  CD  = "CD"
  PWD = "PWD"
  GET = "GET"
  UPL = "UPL"
  CACHE = 4096
)



func main() {
  service := "0.0.0.0:1202"
  tcpAddr, err := net.ResolveTCPAddr("tcp", service)
  checkError(err)

  listener, err := net.ListenTCP("tcp", tcpAddr)
  checkError(err)

  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }
    go HandleClient(conn)
  }
}

func HandleClient(conn net.Conn) {
  defer conn.Close()

  buffer := make([]byte, CACHE)
  for {
    n, err := conn.Read(buffer)
    if err != nil {
      conn.Close()
      return
    }
    s := string(buffer[:n])
    strs := strings.Split(s, " ")
    fmt.Println(s)
    if s[0:2] == CD {
      chdir(conn, s[3:])
    } else if s[0:3] == DIR {
      dirList(conn)
    } else if s[0:3] == PWD {
      pwd(conn)
    } else if s[0:3] == GET {
      dwn(conn, strings.Trim(strs[1], "\r\t\f\v "))
    } else if s[0:3] == UPL {
      n, _ := conn.Read(buffer)
      fileSize, err := strconv.ParseInt(string(buffer[:n]), 10, 64)
      log.Println("filesize =", fileSize)
      if err != nil || fileSize == -1{
        fmt.Println("FILEERROR")
        continue
      }
      upl(conn, strings.Trim(strs[1], "\r\t\f\v "), fileSize)
    }
  }
}

func chdir(conn net.Conn, s string) {
  if os.Chdir(s) == nil {
    conn.Write([]byte("OK"))
  } else {
    conn.Write([]byte("ERROR"))
  }
}

func pwd(conn net.Conn) {
  s, err := os.Getwd()
  if err != nil {
    conn.Write([]byte(""))
  }
  conn.Write([]byte(s))
}

func dwn(conn net.Conn, fileName string) {
  selectedFile, err := os.Open("./" + fileName)
  defer selectedFile.Close()

  if err != nil {
    conn.Write([]byte(err.Error()))
    return
  } else {
    stat,_ :=selectedFile.Stat()
    conn.Write([]byte(strconv.FormatInt(stat.Size(), 10)))
  }
  buffer := make([]byte, CACHE)
  for {
    _, err := selectedFile.Read(buffer)
    if err == io.EOF{
      break
    }
    conn.Write(buffer)
  }
  fmt.Println("File Sent")
}

func upl(conn net.Conn, fileName string, fileSize int64){
  selectedFile, err := os.Create("./" + fileName)
  var fileSizeDownload int64

  if err != nil {
    fmt.Println("Error ", err)
  }

  defer selectedFile.Close()

  for {
    if (fileSize - fileSizeDownload) < CACHE {
      io.CopyN(selectedFile, conn, (fileSize - fileSizeDownload)) 
      conn.Read(make([]byte, (fileSizeDownload + CACHE) - fileSize))
      break
    }
    io.CopyN(selectedFile, conn, CACHE)
    fileSizeDownload += CACHE
  }
  fmt.Println("File Received Succesfully")
}

func dirList(conn net.Conn) {
  defer conn.Write([]byte("\r\n"))

  dir, err := os.Open(".")
  if err != nil {
    return
  }

  names, err := dir.Readdirnames(-1)
  if err != nil {
    return
  }
  for _, nm := range names {
    conn.Write([]byte(nm + "\r\n"))
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("Fatal error", err.Error())
    os.Exit(1)
  }
}
