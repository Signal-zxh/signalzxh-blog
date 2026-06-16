package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

    // 静态页面（知识图谱前端放这里）
    r.Static("/static", "./static")
    
    // API 示例
    r.GET("/api/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello, 个人网站",
        })
    })
    
    // 首页
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "<h1>个人网站</h1><p>服务运行中...</p>")
    })

    r.Run(":8080")  // 监听 8080 端口
}