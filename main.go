package main

import (
	_ "Hertz-Scaffold/biz/handler"
	"Hertz-Scaffold/biz/middleware"
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils/common"
	"Hertz-Scaffold/biz/utils/env"
	"Hertz-Scaffold/conf"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	// 全局参数需要加载到内存中的
	var globalModules []env.Module
	globalModules = append(globalModules, env.UserMapJwtToken)
	env.InitModules(globalModules)
	logger := fmt.Sprintf("######## init model cost %v ms", time.Since(start).Milliseconds())
	fmt.Println(logger)
	common.GlobalLogger.Infof(logger)

	engine := InitRouter(
		middleware.DefaultCorsMiddleware(), // 支持跨域
		middleware.RequestDoTracerId(),     // 全局链路中间件
		middleware.RecoveryMiddleware(),    // 最后捕获panic错误
	)
	proxy()

	//go cron_job.InitCronJob() // 如果需要定时任务 则开启定时任务
	engine.Spin() // 开启Http服务
	defer repository.SqlDbPool.Close()
}

var (
	JdbUsdProxy *httputil.ReverseProxy
)

func proxy() {
	//jdb
	jdbUsdUrl := conf.AppConf.GetGameInfo().JdbUsdUrl
	common.GlobalLogger.Infof("jdbUsdUrl %v", jdbUsdUrl)
	JdbUsdProxy, _ = newProxy(jdbUsdUrl)
	http.HandleFunc("/jdb/usd/apiRequest.do", func(w http.ResponseWriter, r *http.Request) {
		JdbUsdProxy.ServeHTTP(w, r)
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(conf.AppConf.GetBaseInfo().ServicePort+1), nil)
		if err != nil {

		}
	}()
}

func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
	}
	return proxy, nil
}
