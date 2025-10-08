package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"m-db-ui/internal/config"
	"m-db-ui/internal/database"
	"m-db-ui/internal/handlers"
	"html/template"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func printHelp() {
	fmt.Println("MongoDB管理工具")
	fmt.Println("用法:")
	fmt.Println("  m-db-ui [选项]")
	fmt.Println("")
	fmt.Println("选项:")
	fmt.Println("  -host string")
	fmt.Println("        指定运行的主机IP (默认: 127.0.0.1)")
	fmt.Println("  -port string")
	fmt.Println("        指定运行的端口 (默认: 8082)")
	fmt.Println("  -h, --help")
	fmt.Println("        显示帮助信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  m-db-ui                          # 使用默认配置启动")
	fmt.Println("  m-db-ui -host 0.0.0.0 -port 8080 # 指定主机和端口启动")
	fmt.Println("  m-db-ui -port 9000               # 指定端口启动")
}

func main() {
	// 定义命令行参数
	host := flag.String("host", "127.0.0.1", "指定运行的主机IP")
	port := flag.String("port", "8082", "指定运行的端口")
	help := flag.Bool("help", false, "显示帮助信息")

	// 自定义-h标志
	flag.BoolVar(help, "h", false, "显示帮助信息")

	flag.Parse()

	// 如果请求帮助，显示帮助信息并退出
	if *help {
		printHelp()
		os.Exit(0)
	}

	// 设置环境变量，覆盖默认配置
	os.Setenv("HOST", *host)
	os.Setenv("PORT", *port)

	// 加载配置
	cfg := config.Load()

	// 初始化连接管理器
	connectionManager := config.NewConnectionManager("connections.json")
	if err := connectionManager.LoadConnections(); err != nil {
		log.Printf("Failed to load connections: %v", err)
	}
	config.GlobalConnectionManager = connectionManager

	// 获取当前连接配置
	currentConn := connectionManager.GetCurrentConnection()
	if currentConn == nil {
		log.Fatal("No database connection configured")
	}

	// 初始化MongoDB连接
	mongoClient, err := database.Connect(currentConn.GetURI())
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// 初始化服务
	dbService := database.NewService(mongoClient)

	// 初始化处理器
	h := handlers.New(dbService, connectionManager)

	// 设置Gin路由
	r := gin.Default()

	// 添加模板函数
	r.SetFuncMap(template.FuncMap{
		"sub": func(a, b int64) int64 {
			return a - b
		},
		"add": func(a, b int64) int64 {
			return a + b
		},
	})

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// 静态文件服务
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// 路由设置
	api := r.Group("/api/v1")
	{
		// 连接管理
		api.GET("/connections", h.GetConnections)
		api.GET("/connections/:id", h.GetConnection)
		api.POST("/connections", h.AddConnection)
		api.PUT("/connections/:id", h.UpdateConnection)
		api.DELETE("/connections/:id", h.DeleteConnection)
		api.POST("/connections/:id/current", h.SetCurrentConnection)
		api.GET("/connections/current", h.GetCurrentConnection)
		api.POST("/connections/test", h.TestConnection)

		// 数据库相关
		api.GET("/databases", h.GetDatabases)
		api.POST("/databases", h.CreateDatabase)
		api.GET("/databases/:name", h.GetDatabase)
		api.DELETE("/databases/:name", h.DeleteDatabase)

		// 统计信息
		api.GET("/stats", h.GetStats)

		// 集合相关 - 使用不同的路径避免冲突
		api.GET("/db/:db/collections", h.GetCollections)
		api.POST("/db/:db/collections", h.CreateCollection)
		api.DELETE("/db/:db/collections/:collection", h.DeleteCollection)

		// 文档相关
		api.GET("/db/:db/collections/:collection/documents", h.GetDocuments)
		api.GET("/db/:db/collections/:collection/documents/:id", h.GetDocument)
		api.POST("/db/:db/collections/:collection/documents", h.CreateDocument)
		api.PUT("/db/:db/collections/:collection/documents/:id", h.UpdateDocument)
		api.DELETE("/db/:db/collections/:collection/documents/:id", h.DeleteDocument)
		api.POST("/db/:db/collections/:collection/query", h.QueryDocuments)
	}

	// Web界面路由
	r.GET("/", h.Index)
	r.GET("/connections", h.ConnectionsPage)
	r.GET("/database/:db", h.DatabasePage)
	r.GET("/database/:db/collection/:collection", h.CollectionPage)

	// 启动服务器
	address := cfg.Host + ":" + cfg.Port
	log.Printf("Server starting on %s", address)
	log.Fatal(r.Run(address))
}