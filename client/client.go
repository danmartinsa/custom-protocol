package client

import (
	"bufio"
	"fmt"
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

const CACHE = 512

func mainLoop() {
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
