package main

import (
	"github.com/scg130/tools"
	"kanshu/router"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/web"

	"log"
)

// @title kanshu
// @version 1.0
// @description micro

// @contact.name scg130
// @contact.url
// @contact.email scg130@163.com

// @schemes http
// @host http://www.scg130.cn
// @base http://www.scg130.cn
// @description	 micro

func main() {
	go func() {
		h := hystrix.NewStreamHandler()
		h.Start()
		err := http.ListenAndServe(net.JoinHostPort("", "81"), h)
		if err != nil {
			panic(err)
		}
	}()

	srv := web.NewService(
		web.Name("go.kanshu.web"),
		web.Registry(tools.Reg()),
		web.Address(":2222"),
		web.Handler(router.HttpRouter()),
	)

	if err := srv.Init(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
