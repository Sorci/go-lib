package main

import (
	"fmt"
	"go-lib/container/pool"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr string = "127.0.0.1:8098"

func main()  {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	go server()
	//等待tcp server启动
	time.Sleep(2 * time.Second)
	client()
	fmt.Println("使用: ctrl+c 退出服务")
	<-c
	fmt.Println("服务退出")
}


func client() {

	//factory 创建连接的方法
	factory := func() (net.Conn, error) { return net.Dial("tcp", addr) }

	//创建一个连接池： 初始化2，最大连接5，空闲连接数是4
	poolConfig := &pool.Config{
		InitialCap: 2,
		MaxCap:     5,
		Factory:    factory,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 5 * time.Second,
	}
	p, err := pool.New(poolConfig)
	if err != nil {
		fmt.Println("pool new failed: ", err)
	}

	conn, err := p.Get()
	if err != nil {
		fmt.Println("pool get conn failed: ", err)
	}
	// do something
	conn.Close()

	fmt.Println("len=", p.Len())
}

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}

		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
	}
}
