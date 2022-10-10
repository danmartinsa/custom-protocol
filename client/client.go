package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
  uiDir  = "ls"
  uiPwd  = "pwd"
  uiCd   = "cd"
  uiQuit = "quit"
  uiGet  = "get"
  uiSend = "send"
  uiMan  = "man"
)

const (
  DIR = "DIR"
  CD  = "CD"
  PWD = "PWD"
  GET = "GET"
  UPL = "UPL"
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

  
  handle(conn)
}

func handle(conn net.Conn) {

  reader := bufio.NewReader(os.Stdin)
  buffer := make([]byte, CACHE)
  for {
    fmt.Printf("ctp > ")
    command, err := reader.ReadString('\n')
    line := strings.TrimRight(command, "\t\r\n")
    if err != nil {
      break
    }
    strs := strings.SplitN(line, " ", 2)
    fmt.Println(line)

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
      _, err := conn.Write([]byte(GET + " " + strs[1]))
      checkWriteError(err)
      n, _ := conn.Read(buffer)
      fileSize, err := strconv.ParseInt(string(buffer[:n]), 10, 64)
      if err != nil {
        fmt.Println("ERROR: ", string(buffer[:n]))
        continue
      }
      GetFileRequest(conn, strs[1], fileSize)

    case uiSend:
      if len(strs) == 1 {
        fmt.Println("Please Provide a filename")
        continue
      }
      _, err := conn.Write([]byte(UPL + " " + strs[1]))
      checkWriteError(err)
      fmt.Println("###", strs)
      SendFileRequest(conn, strs[1])
    case uiQuit:
      conn.Close()
      os.Exit(0)
    case uiMan:
      fmt.Println("ls          list file and folder")
      fmt.Println("cd          change current folder")
      fmt.Println("pwd         print current folder")
      fmt.Println("get  <file> file download request")
      fmt.Println("send <file> file upload request")
      fmt.Println("quit        quit the application")
    default:
      fmt.Println("Unknow command, type \"man\" to see accepeted words and arguments")
    }
  }
}

func dirRequests(conn net.Conn) {
  _, err := conn.Write([]byte(DIR + " "))
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  buffer := make([]byte, CACHE)
  n, _ := conn.Read(buffer[0:])
  s := string(buffer[0:n])
  // fmt.Println(s)
  fileList := strings.Split(s, "#")
  // fmt.Println(len(fileList), fileList)
  for _, file := range fileList[:len(fileList)-1] {
    fmt.Println(file)
  }
}

func cdRequest(conn net.Conn, dir string) {
  _, err := conn.Write([]byte(CD + " " + dir))
  checkWriteError(err)
  var response [CACHE]byte
  n, _ := conn.Read(response[0:])
  s := string(response[0:n])
  if s != "OK" {
    fmt.Println("Failed to change dir")
  }
}

func pwdRequest(conn net.Conn) {
  _, err := conn.Write([]byte(PWD))
  checkWriteError(err)
  var response [CACHE]byte
  n, _ := conn.Read(response[0:])
  s := string(response[0:n])
  fmt.Println("Current dir \"" + s + "\"")
}

func GetFileRequest(conn net.Conn, fileName string, fileSize int64) {
  newFile, err := os.Create("./" + fileName)
  var fileSizeDownload int64

  if err != nil {
    fmt.Println(err)
  }
  defer newFile.Close()

  for {
    if (fileSize - fileSizeDownload) < CACHE {
      _, err := io.CopyN(newFile, conn, (fileSize - fileSizeDownload))
      checkWriteError(err)
      break
    }
    _, err := io.CopyN(newFile, conn, CACHE)
    checkWriteError(err)
    fileSizeDownload += CACHE
  }
  fmt.Println("File Received Succesfully")
}

func SendFileRequest(conn net.Conn, fileName string) {
  newFile, err := os.Open("./" + fileName)
  defer newFile.Close()

  if err != nil {
    fmt.Println("Error ", err.Error())
    _, err := conn.Write([]byte("-1"))
    if err != nil {
      fmt.Println(err.Error())
    }
    return
  } else {
    stat, _ := newFile.Stat()
    _, err := conn.Write([]byte(strconv.FormatInt(stat.Size(), 10)))
    if err != nil {
      fmt.Println(err.Error())
      return
    }
  }

  buffer := make([]byte, CACHE)
  for {
    _, err := newFile.Read(buffer)
    if err == io.EOF {
      break
    }
    _, error := conn.Write(buffer)
    if error != nil {
      fmt.Println(err.Error())
      return
    }
  }
  fmt.Println("File Sent")
}

func checkError(err error) {
  if err != nil {
    fmt.Println("Fatal error", err.Error())
    os.Exit(1)
  }
}

func checkWriteError(err error) {
  if err != nil {
    fmt.Println("Write Error", err.Error())
  }
}
