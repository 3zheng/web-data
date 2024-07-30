package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	tablestruct "github.com/3zheng/webdata/table_struct"
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
	Name      string `json:"name"`
	ForceIPv4 int    `json:"force ipv4"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
}

type Config struct {
	Database DBConfig     `json:"database config"`
	Server   ServerConfig `json:"server config"`
}

func main() {

	content, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

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
	r.LoadHTMLGlob("HTML/*") //加载HTML文件

	r.GET("/KC", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		inventories := GetInventory(db)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"data": inventories,
		})
		//c.JSON(200, inventories)
	})

	addr := fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port)
	//ln := net.Listener
	if config.Server.ForceIPv4 == 1 {
		// 强制使用IPv4
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
		r.Run(addr)
	}
	//r.Run(addr)
}

func GetInventory(db *sql.DB) [](*tablestruct.Inventory) {
	//编写查询语句
	stmt, err := db.Prepare(`select 产品型号, 产品名称, 库存数量 from dbo.View_KC`)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	//将数据读取到实体中
	var rowsData [](*tablestruct.Inventory)
	for rows.Next() {
		data := new(tablestruct.Inventory)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.ProductID, &data.ProductName, &data.ResidualNum)
		rowsData = append(rowsData, data)
	}
	return rowsData
}
