package tablestruct

import (
	"database/sql"
	"log"
)

//这是数据库对应的表结构

// 库存表对应dbo.View_KC
type Inventory struct {
	ProductID   string  `json:"ID"`          //产品类型
	ProductName string  `json:"Name"`        //产品名字
	ResidualNum float64 `json:"ResidualNum"` //库存数量
}

func GetInventory(db *sql.DB) [](*Inventory) {
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
