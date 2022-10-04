package main

import "flag"

var port int
var rootDir string

fun init() {
  flag.Intvar(&port, "port", 8080, "port number")
  flag.StringVar()
}
