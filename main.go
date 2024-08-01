package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	tablemiddleware "github.com/3zheng/web-data/table_middleware"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type DBConfig struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DB       string `json:"database"`
	UserId   string `json:"user id"`
	Password string `json:"password"`
}

// IP   string `json:"ip"`
type ServerConfig struct {
	Path      string `json:"path"`
	ForceIPv4 int    `json:"force ipv4"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
}

type Config struct {
	Database DBConfig     `json:"database config"`
	Server   ServerConfig `json:"server config"`
}

func InitLog() {
	var filepath string
	fmt.Println(runtime.GOOS)
	if runtime.GOOS == "windows" {
		filepath = "./logfile"
	} else if runtime.GOOS == "linux" {
		filepath = "/home/wuzhibin86/workspace/web-data/logfile"
	} else {
		fmt.Println("系统不明")
		os.Exit(0)
	}
	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open log file failed.")
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
	log.Println("log file opened.")

}

func recovermain() {
	if err := recover(); err != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		log.Fatalf("[panic] err: %v\nstack: %s\n", err, buf[:n])
	}
}

func main() {
	defer recovermain() //退出前打印异常
	InitLog()
	//读取配置文件
	args := os.Args //main命令行参数
	log.Println("main args = ", args)

	var content []byte
	var err error
	//这时候如果不输入的命令行参数就会导致panic
	if args[1] == "-config" {
		log.Println("读取配置文件：", args[2])
		content, err = os.ReadFile(args[2])
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
	} else { //不带参数直接返回
		return
	}

	//tablemiddleware.CreatePanic()
	// Now let's unmarshall the data into `payload`
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	//connString := "server=127.0.0.1;port=38336;database=user;user id=admin;password=123456"
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		config.Database.IP, config.Database.Port, config.Database.DB, config.Database.UserId, config.Database.Password)

	//建立数据库连接：db
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}
	defer db.Close()

	//release模式
	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	htmlPath := "HTML/*"
	if runtime.GOOS == "linux" {
		htmlPath = config.Server.Path + htmlPath
	}
	log.Println(htmlPath)
	r.LoadHTMLGlob(htmlPath) //加载HTML文件

	r.GET("/KC", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		fmt.Println("/KC GET require")
		inventories := tablemiddleware.GetInventory(db)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"data": inventories,
		})
		//c.JSON(200, inventories)
	})
	log.Println("开始启动web服务")
	addr := fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port)
	//ln := net.Listener
	if config.Server.ForceIPv4 == 1 {
		// 强制使用IPv4
		log.Println("强制使用IPv4")
		server := &http.Server{Addr: addr, Handler: r}
		ln, err := net.Listen("tcp4", addr)
		if err != nil {
			panic(err)
		}
		type tcpKeepAliveListener struct {
			*net.TCPListener
		}

		server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	} else {
		log.Println("http服务地址：", addr)
		r.Run(addr)
	}
	//r.Run(addr)
}
