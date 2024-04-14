package main

import (
	// "Wechat-project/config"
    // "log"
	"Wechat-project/router"
)

func main(){
	// 插入数据，插入完后注释掉了
	// if err := utils.InsertData(); err != nil {
	// 	log.Fatalf("Failed to insert data: %v", err)
	// }
	r := router.Router()
	r.Run(":9999")
}