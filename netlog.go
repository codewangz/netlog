package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)
var serverAddr,projectName,dir string

func Read(con net.Conn){
	data := make([]byte, 1000)
	for{
		n, err := con.Read(data)
		if err != nil{
			fmt.Println(err)
			break
		}
		go saveLog(string(data[0:n]))
	}
}

func saveLog(log string){
	date := time.Now().Format("2006-01-02")
	filePath := dir+"/"+projectName+"/"+date+".log"
	var file *os.File
	var err error
	if _,err = os.Stat(filePath);err == nil {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	}

	defer file.Close()
	if  err != nil{
		fmt.Println("创建日志文件失败",filePath)
	}
	write := bufio.NewWriter(file)
	write.WriteString(log)
	write.Flush()
}

func main(){
	flag.StringVar(&serverAddr, "server", "0.0.0.0:5921", "./netlog -server ip:port")
	flag.StringVar(&projectName, "project", "normal", "./netlog -project pojectName")
	flag.StringVar(&dir, "dir", "/tmp", "./netlog -path /tmp")
	flag.Parse()
	path := dir+"/"+projectName
	err := os.MkdirAll(path,0777)
	if err != nil{
		fmt.Println(err)
		return
	}
	listen, err := net.Listen("tcp", serverAddr)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("listen on "+serverAddr," from project "+projectName, " log in: "+dir+"/"+projectName)
	for{
		con, err := listen.Accept()
		if err != nil{
			fmt.Println(err)
			continue
		}
		go Read(con)
	}

}