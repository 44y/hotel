package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	check_prd int = 5
)

var (
	DIR        string //txt file directory
	s_port     int
	WAIT_GROUP sync.WaitGroup
)

type conf_st struct {
	Dir  string `json:"dir"`  //用户日志目录
	Port int    `json:"port"` //服务端口
}

func srvRun() {

	/*parse config*/
	tmp_conf, _ := parseConfig()
	DIR = tmp_conf.Dir
	s_port = tmp_conf.Port
	srv := &http_svc{}
	srv.addr = fmt.Sprintf(":%d", s_port)
	srv.svc_t = &http.Server{Addr: srv.addr}
	http.HandleFunc("/", srv.handle)
	srv.out_ch = make(chan int, 1)
	go srv.StartHttpSvc()

	ch := make(chan int, 1)
	go UpdateDir(ch)

	for {
		time.Sleep(1 * time.Second)
	}
}

/*
初始化服务：
    1. 查找用户日志路径
    2. 确认可用端口
*/
func srvInit() bool {
	dft_conf, conf_path := parseConfig()

	changed := false
	//默认用户日志路径是否存在
	if isDirExist(dft_conf.Dir) == false {
		//搜索用户日志路径
		fmt.Println("搜索用户日志路径...")
		is_find, new_dir := searchUserLog()
		if is_find {
			dft_conf.Dir = new_dir
			changed = true
		} else {
			fmt.Println("未找到用户日志路径")
			return false
		}
	}
	//默认端口是否可用
	if isPortAvl(dft_conf.Port) == false {
		new_port := getAvlPort()
		if new_port == 0 {
			fmt.Println("获取端口失败")
			return false
		} else {
			dft_conf.Port = new_port
			changed = true
		}
	}

	if changed {
		updateConfig(dft_conf, conf_path)
	}
	fmt.Println("当前使用端口：", dft_conf.Port)
	return true
}
