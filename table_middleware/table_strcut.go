package tablestruct

import (
	"database/sql"
	"log"
	"os"
)

//这是数据库对应的表结构

// 库存表对应dbo.View_KC
type Inventory struct {
	//产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,城市名称,仓库名称,库存数量,库存成本,最近30天销售数量,消纳时间
	ProductID          string  `json:"ID"`                //产品型号
	ProductName        string  `json:"Name"`              //产品名称
	ProductDescription string  `json:"Descrption"`        //产品描述
	ProductSubclass    string  `json:"Subclass"`          //产品类型名称,二级分类
	ProductSuperClass  string  `json:"ProductSuperClass"` //主类型名称,一级分类
	ProductUnitName    string  `json:"ProductUnitName"`   //产品单位名称
	CityName           string  `json:"CityName"`          //城市名称
	ResidualNum        float64 `json:"ResidualNum"`       //库存数量
}

func CreatePanic() {
	if os.Args[4] == "" {
		log.Println("arg 4:", os.Args[4])
	}
}

func GetInventory(db *sql.DB) [](*Inventory) {
	//编写查询语句
	//select 产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,城市名称,仓库名称,库存数量,库存成本,最近30天销售数量,消纳时间 from dbo.View_KC
	stmt, err := db.Prepare(`select 产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,
		城市名称,仓库名称,库存数量,库存成本,最近30天销售数量,消纳时间 
		from dbo.View_KC`)
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
	var rowsData [](*Inventory)
	for rows.Next() {
		data := new(Inventory)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.ProductID, &data.ProductName, &data.ResidualNum)
		rowsData = append(rowsData, data)
	}
	return rowsData
}
