//go:generate goversioninfo -icon=ikuai.ico

package main

import (
	"fmt"
	"github.com/kardianos/service"
	"os"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go srvRun()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}
func main() {

	//Log_init(DEBUG)

	svcConfig := &service.Config{
		Name:        "IK认证服务器", //服务显示名称
		DisplayName: "IK认证服务器", //服务名称
		Description: "IK认证服务器", //服务描述
	}
	//依赖服务
	svcConfig.Dependencies = append(svcConfig.Dependencies, "Tcpip")

	prg := &program{}
	s, err := service.New(prg, svcConfig)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err = s.Install()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			//fmt.Println("服务安装成功")
			os.Exit(0)
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Println("服务卸载成功")
			os.Exit(0)
		}

		if os.Args[1] == "test" {
			if srvInit() {
				fmt.Println("初始化成功")
				os.Exit(0)
			} else {
				fmt.Println("初始化失败")
				os.Exit(1)
			}
		}
		//debug模式，日志输出到前台
		if os.Args[1] == "-D" {
			Log_init(DEBUG)
		} else {
			fmt.Println("Parameter error!")
			os.Exit(1)
		}
	} else {
		Log_init(NORMAL)
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
