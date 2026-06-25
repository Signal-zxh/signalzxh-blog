package main

import (
	"log"

	"github.com/Signal-zxh/signal-zxh/db"
	"github.com/Signal-zxh/signal-zxh/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// 加载配置
	if err := godotenv.Load(); err != nil {
		log.Fatal("load env failed:", err)
	}
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatal("db connect failed:", err)
	}
	defer db.DB.Close()
	// 初始化Redis
	if err := db.InitRedis(); err != nil {
		log.Fatal("redis connect failed:", err)
	}
	defer db.RDB.Close()
	// 初始化路由
	r := router.SetupRouter()
	// 启动服务器
	r.Run(":8080") // 监听 8080 端口
}
