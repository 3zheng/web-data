package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	tablestruct "github.com/3zheng/webdata/table_struct"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type DBConfig struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	DB       string `json:"database"`
	UserId   string `json:"user id"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Name string `json:"name"`
	Port int    `json:"port"`
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

	connString := "server=127.0.0.1;port=38336;database=user;user id=admin;password=123456"
	//建立数据库连接：db
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/KC", func(c *gin.Context) {
		//var inventories [](*tablestruct.Inventory)
		inventories := GetInventory(db)
		c.JSON(200, inventories)
	})

	r.Run()
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
