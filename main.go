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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	Database  DBConfig     `json:"database config"`
	Server    ServerConfig `json:"server config"`
	MysqlConn string       `json:"mysqlConn"`
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

func Recovermain() {
	if err := recover(); err != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		log.Printf("[panic] err: %v\nstack: %s\n", err, buf[:n])
	}
}

// 设置http路由,直接返回Html
func SetGinRouterByHtml(r *gin.Engine, db *sql.DB, projectPath string) {
	htmlPath := "HTML/*.html"
	if runtime.GOOS == "linux" {
		htmlPath = projectPath + htmlPath
	}
	log.Println(htmlPath)
	r.LoadHTMLGlob(htmlPath) //加载HTML文件

	//注册http路由
	r.GET("/KC", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/KC GET require")
		datas := tablemiddleware.GetInventory(db)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})
	r.GET("/KC2", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/KC2 GET require")
		datas := tablemiddleware.GetInventorySummary(db)
		c.HTML(http.StatusOK, "kc2.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})
	r.GET("/QK", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/QK GET require")
		datas := tablemiddleware.GetDebtDaily(db)
		c.HTML(http.StatusOK, "qk.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})
	r.GET("/XS1", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/XS1 GET require")
		datas := tablemiddleware.GetSalesManDailyRecord(db)
		c.HTML(http.StatusOK, "xs1.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})
	r.GET("/CYSP", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/CYSP GET require")
		datas := tablemiddleware.GetImportantCustomer(db)
		c.HTML(http.StatusOK, "cysp.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})
	r.GET("/LKC", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/QK GET require")
		datas := tablemiddleware.GetLostImportantCustomeromer(db)
		if datas == nil {
			c.JSON(http.StatusRequestTimeout, "数据库连接出错")
		}
		c.HTML(http.StatusOK, "lkc.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})
	r.GET("/NKC", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/QK GET require")
		datas := tablemiddleware.GetNewImportantCustomer(db)
		c.HTML(http.StatusOK, "nkc.html", gin.H{
			"data": datas,
		})
		//c.JSON(200, inventories)
	})

}

func SelectResponseJson[T any](c *gin.Context, datas []*T) {
	var partialDatas []*T
	vol := c.Query("volume") //获取volume参数
	if len(datas) > 200 {
		partialDatas = datas[:200]
	} else {
		partialDatas = datas
	}
	if vol == "all" {
		c.JSON(http.StatusOK, datas) //发送所有数据
	} else if vol == "partial" {
		c.JSON(http.StatusOK, partialDatas) //发送部分数据
	} else {
		c.JSON(http.StatusOK, datas) //发送所有数据
	}
}

// 返回内容为json格式的字符串
func SetGinRouterByJson(r *gin.Engine, mc *tablemiddleware.MemoryCache) {
	r.GET("/api/inventory_warehouse", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/inventory_warehouse GET require")
		var datas []*tablemiddleware.Inventory
		mc.GetMemoryCache(&datas)
		SelectResponseJson(c, datas)
	})
	r.GET("/api/inventory_summary", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/inventory_summary GET require")
		var datas []*tablemiddleware.InventorySummary
		mc.GetMemoryCache(&datas)
		SelectResponseJson(c, datas)
	})
	r.GET("/api/debt_daily", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/debt_daily GET require")
		var datas []*tablemiddleware.DebtDaily
		mc.GetMemoryCache(&datas)
		SelectResponseJson(c, datas)
	})
	r.GET("/api/debt_summary", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/debt_summary GET require")
		var datas []*tablemiddleware.DebtSummary
		mc.GetMemoryCache(&datas)
		c.JSON(http.StatusOK, datas)
	})
	r.GET("/api/sales_record", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/sales_detail GET require")
		var datas []*tablemiddleware.SalesmanDaily
		mc.GetMemoryCache(&datas)
		SelectResponseJson(c, datas)
	})
	r.GET("/api/sales_summary", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/sales_summary GET require")
		var datas []*tablemiddleware.SalesmanMonthly
		mc.GetMemoryCache(&datas)
		SelectResponseJson(c, datas)
	})
	r.GET("/api/key_customer", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/key_customer GET require")
	})
	r.GET("/api/lost_key_customer", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/lost_key_customer GET require")
	})
	r.GET("/api/new_key_customer", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/new_key_customer GET require")
	})

	//从mysql数据库里取数据
	r.GET("/wp1", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		log.Println("/QK GET require")
		//datas := tablemiddleware.GetWordpress(mc.)
		//c.JSON(200, datas)
	})
}

func main() {
	defer Recovermain() //退出前打印异常
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

	//mysql的连接字符串格式
	//connString := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8"

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		config.Database.IP, config.Database.Port, config.Database.DB, config.Database.UserId, config.Database.Password)

	//建立SQLSever数据库连接：db

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}

	//db, err := sql.Open("mysql", connString)
	defer db.Close()

	//release模式
	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.Default()) //使用cors，解决跨域问题
	mc := new(tablemiddleware.MemoryCache)
	mc.InitMemoryCache(db)
	//SetGinRouterByHtml(r, db, config.Server.Path)//直接返回Html网页，把前端后端放一起
	SetGinRouterByJson(r, mc) //返回json数据，前端后端分离，后端只返回数据，前端不管

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
