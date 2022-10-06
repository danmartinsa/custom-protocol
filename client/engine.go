package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

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
