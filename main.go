package main

import (
	_ "Hertz-Scaffold/biz/handler"
	"Hertz-Scaffold/biz/middleware"
	"Hertz-Scaffold/biz/repository"
	"Hertz-Scaffold/biz/utils/common"
	"Hertz-Scaffold/biz/utils/env"
	"Hertz-Scaffold/conf"
	"database/sql"
	"fmt"
	"net/http"
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
	//go cron_job.InitCronJob() // 如果需要定时任务 则开启定时任务
	go engine.Spin() // 开启Http服务
	//handler.JdbProxy()
	func() {
		err := http.ListenAndServe(":"+strconv.Itoa(conf.AppConf.GetBaseInfo().ServicePort+1), nil)
		if err != nil {

		}
	}()
	defer func(SqlDbPool *sql.DB) {
		err := SqlDbPool.Close()
		if err != nil {

		}
	}(repository.SqlDbPool)
}
