package tablestruct

import (
	"database/sql"
	"log"
	"os"
	"time"
)

//这是数据库对应的表结构

func CreatePanic() {
	if os.Args[4] == "" {
		log.Println("arg 4:", os.Args[4])
	}
}

// 库存表对应dbo.View_KC
type Inventory struct {
	//select 产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,城市名称,仓库名称,库存数量,库存成本,最近30天销售数量,消纳时间 from dbo.View_KC
	ProductID           string  `json:"ID"`                  //产品型号
	ProductName         string  `json:"Name"`                //产品名称
	ProductDescription  string  `json:"Descrption"`          //产品描述
	ProductSubclass     string  `json:"Subclass"`            //产品类型名称,二级分类
	ProductSuperClass   string  `json:"ProductSuperClass"`   //主类型名称,一级分类
	ProductUnitName     string  `json:"ProductUnitName"`     //产品单位名称
	CityName            string  `json:"CityName"`            //城市名称
	WarehouseName       string  `json:"WarehouseName"`       //仓库名称
	ResidualNum         float64 `json:"ResidualNum"`         //库存数量
	InventoryCost       float64 `json:"InventoryCost"`       //库存成本
	SalesQuantity30days int     `json:"SalesQuantity30days"` //最近30天销售数量
	UnsalableScale      float64 `json:"UnsalableScale"`      //消纳时间 = 库存数量/最近30天销售数量,表示库存将在多少个月内售罄,相当于滞销度
}

func GetInventory(db *sql.DB) [](*Inventory) {
	//编写查询语句
	//select 产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,城市名称,仓库名称,库存数量,库存成本,最近30天销售数量,消纳时间 from dbo.View_KC
	sqlStr := `select trim(产品型号),trim(产品名称),trim(产品描述),trim(产品类型名称),trim(主类型名称),trim(产品单位名称),` +
		`trim(城市名称),trim(仓库名称),库存数量,库存成本,最近30天销售数量,消纳时间 ` +
		`from dbo.View_KC order by 产品型号,仓库名称`
	//sqlStr := `select 产品型号,产品名称,最近30天销售数量,消纳时间 from dbo.View_KC`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*Inventory)
	for rows.Next() {
		data := new(Inventory)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.ProductID, &data.ProductName, &data.ProductDescription,
			&data.ProductSubclass, &data.ProductSuperClass, &data.ProductUnitName, &data.CityName,
			&data.WarehouseName, &data.ResidualNum, &data.InventoryCost, &data.SalesQuantity30days,
			&data.UnsalableScale)

		rowsData = append(rowsData, data)
	}
	return rowsData
}

// 整体库存表对应dbo.View_KC2
type InventorySummary struct {
	//select 产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,库存数量,库存成本,最近30天销售数量 from View_KC2
	ProductID           string  `json:"ID"`                  //产品型号
	ProductName         string  `json:"Name"`                //产品名称
	ProductDescription  string  `json:"Descrption"`          //产品描述
	ProductSubclass     string  `json:"Subclass"`            //产品类型名称,二级分类
	ProductSuperClass   string  `json:"ProductSuperClass"`   //主类型名称,一级分类
	ProductUnitName     string  `json:"ProductUnitName"`     //产品单位名称
	ResidualNum         float64 `json:"ResidualNum"`         //库存数量
	InventoryCost       float64 `json:"InventoryCost"`       //库存成本
	SalesQuantity30days float64 `json:"SalesQuantity30days"` //最近30天销售数量
}

func GetInventorySummary(db *sql.DB) [](*InventorySummary) {
	//编写查询语句
	//select 产品型号,产品名称,产品描述,产品类型名称,主类型名称,产品单位名称,库存数量,库存成本,最近30天销售数量 from View_KC2
	stmt, err := db.Prepare(`select trim(产品型号),trim(产品名称),trim(产品描述),trim(产品类型名称),trim(主类型名称),
		trim(产品单位名称),库存数量,库存成本,最近30天销售数量 
		from View_KC2`)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*InventorySummary)
	for rows.Next() {
		data := new(InventorySummary)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.ProductID, &data.ProductName, &data.ProductDescription,
			&data.ProductSubclass, &data.ProductSuperClass, &data.ProductUnitName,
			&data.ResidualNum, &data.InventoryCost, &data.SalesQuantity30days)
		rowsData = append(rowsData, data)
	}
	return rowsData
}

// 欠款表对应dbo.View_QK
type Debt struct {
	//select 欠款客户名称,欠款金额,金额单位,欠款订单时间,欠款时长,时间单位,销售员,欠款客户编号,销售员编号 from dbo.View_QK
	DebtorName    string  `json:"DebtorName"`    //欠款客户名称
	DebtAmount    float64 `json:"DebtAmount"`    //欠款金额
	CurrencyUnit  string  `json:"CurrencyUnit"`  //货币单位，写死为'Bs'，可以不用select
	OrderFormDate string  `json:"OrderFormDate"` //欠款订单时间,数据库里这是个datetime类型对应golang的time.Time类型
	DebtDuration  int     `json:"DebtDuration"`  //欠款时长，DATEDIFF(DAY, dbo.cp022.va_fec_doc, GETDATE()) AS 欠款时长
	DateUnit      string  `json:"DateUnit"`      //货币单位，写死为'días'，可以不用select
	Salesman      string  `json:"Salesman"`      //销售员
	DebtorID      int     `json:"DebtorID"`      //欠款客户编号
	SalesmanID    int     `json:"SalesmanID"`    //销售员编号
}

func GetDebt(db *sql.DB) [](*Debt) {
	//编写查询语句
	//select 欠款客户名称,欠款金额,金额单位,欠款订单时间,欠款时长,时间单位,销售员,欠款客户编号,销售员编号 from dbo.View_QK
	stmt, err := db.Prepare(`select trim(欠款客户名称),欠款金额,trim(金额单位),欠款订单时间,欠款时长,
		trim(时间单位),trim(销售员),欠款客户编号,销售员编号 from dbo.View_QK`)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
	}
	//将数据读取到实体中
	var rowsData [](*Debt)
	var OrderFormDate time.Time
	for rows.Next() {
		data := new(Debt)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.DebtorName, &data.DebtAmount, &data.CurrencyUnit, &OrderFormDate,
			&data.DebtDuration, &data.DateUnit, &data.Salesman,
			&data.DebtorID, &data.SalesmanID)
		data.OrderFormDate = OrderFormDate.Format(time.DateOnly)
		rowsData = append(rowsData, data)
	}
	return rowsData
}

// 销售员每日销售额表对应dbo.View_XS1
type Salesman struct {
	//select 销售日期,销售员姓名,销售总金额,订单数量 from dbo.View_XS1
	SalesDate    string  `json:"SalesDate"`    //销售日期
	Name         string  `json:"Name"`         //销售员姓名
	SalesAmount  float64 `json:"SalesAmount"`  //销售总额
	OrderFormNum int     `json:"OrderFormNum"` //订单数量
}

func GetSalesman(db *sql.DB) [](*Salesman) {
	//编写查询语句
	//select 销售日期,销售员姓名,销售总金额,订单数量 from dbo.View_XS1
	stmt, err := db.Prepare(`select 销售日期,trim(销售员姓名),销售总金额,订单数量 from dbo.View_XS1`)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*Salesman)
	var SalesDate time.Time
	for rows.Next() {
		data := new(Salesman)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&SalesDate, &data.Name, &data.SalesAmount, &data.OrderFormNum)
		data.SalesDate = SalesDate.Format(time.DateOnly)
		rowsData = append(rowsData, data)
	}
	return rowsData
}

// 重点客户表对应dbo.CustomerYearlySalesReport
type ImportantCustomer struct {
	//select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from dbo.CustomerYearlySalesReport
	CustomerID        int     `json:"CustomerID"`        //客户ID
	CustomerName      string  `json:"CustomerName"`      //客户姓名
	Month             string  `json:"Month"`             //月开始日期
	ConsumptionAmount float64 `json:"ConsumptionAmount"` //月购买总金额
	ConsumptionTimes  int     `json:"ConsumptionTimes"`  //月购买次数
}

func GetImportantCustomer(db *sql.DB) [](*ImportantCustomer) {
	//编写查询语句
	//select 客户ID,客户姓名,月开始时间,月购买总金额,月购买次数 from dbo.CustomerYearlySalesReport
	stmt, err := db.Prepare(`select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from dbo.CustomerYearlySalesReport`)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*ImportantCustomer)
	var Month time.Time
	for rows.Next() {
		data := new(ImportantCustomer)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.CustomerID, &data.CustomerName, &Month, &data.ConsumptionAmount,
			&data.ConsumptionTimes)
		data.Month = Month.Format(time.DateOnly)
		rowsData = append(rowsData, data)
	}
	return rowsData
}

// 丢失的关键客户表对应dbo.LostKeyCustomers
type LostImportantCustomer struct {
	//select 客户ID,客户姓名,当前月,月购买金额,月购买次数,距离上次购买月数 from dbo.LostKeyCustomers
	CustomerID         int     `json:"CustomerID"`         //客户ID
	CustomerName       string  `json:"CustomerName"`       //客户姓名
	CurrentMonth       string  `json:"CurrentMonth"`       //当前月
	ConsumptionAmount  float64 `json:"ConsumptionAmount"`  //月购买总金额
	ConsumptionTimes   int     `json:"ConsumptionTimes"`   //月购买次数
	MonthSinceLastTime int     `json:"MonthSinceLastTime"` //距离上次购买月数
}

func GetLostImportantCustomeromer(db *sql.DB) [](*LostImportantCustomer) {
	//编写查询语句
	//select 客户ID,客户姓名,当前月,月购买总金额,月购买次数,距离上次购买月数 from dbo.LostKeyCustomers
	stmt, err := db.Prepare(`select 客户ID,客户姓名,当前月,月购买总金额,月购买次数,距离上次购买的月数 from dbo.LostKeyCustomers`)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*LostImportantCustomer)
	var CurrentMonth time.Time
	for rows.Next() {
		data := new(LostImportantCustomer)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.CustomerID, &data.CustomerName, &CurrentMonth, &data.ConsumptionAmount,
			&data.ConsumptionTimes, &data.MonthSinceLastTime)
		data.CurrentMonth = CurrentMonth.Format(time.DateOnly)
		rowsData = append(rowsData, data)
	}
	return rowsData
}

// 新增的关键客户表对应dbo.LostKeyCustomers
type NewImportantCustomer struct {
	//select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from NewKeyCustomers
	CustomerID        int     `json:"CustomerID"`        //客户ID
	CustomerName      string  `json:"CustomerName"`      //客户姓名
	Month             string  `json:"Month"`             //月开始日期
	ConsumptionAmount float64 `json:"ConsumptionAmount"` //月购买总金额
	ConsumptionTimes  int     `json:"ConsumptionTimes"`  //月购买次数
}

func GetNewImportantCustomer(db *sql.DB) [](*NewImportantCustomer) {
	//编写查询语句
	//select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from dbo.CustomerYearlySalesReport

	stmt, err := db.Prepare(`select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from dbo.CustomerYearlySalesReport`)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*NewImportantCustomer)
	var Month time.Time
	for rows.Next() {
		data := new(NewImportantCustomer)
		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.CustomerID, &data.CustomerName, &Month, &data.ConsumptionAmount,
			&data.ConsumptionTimes)
		data.Month = Month.Format(time.DateOnly)
		rowsData = append(rowsData, data)
	}
	return rowsData
}

// wordpress的mysql测试数据
type Wordpress struct {
	//select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from NewKeyCustomers
	ID             int       `json:"ID"`             //客户ID
	Post_author    string    `json:"Post_author"`    //客户姓名
	Post_date      time.Time `json:"Post_date"`      //月开始日期
	Post_date_gmt  time.Time `json:"Post_date_gmt"`  //月开始日期
	Post_content   string    `json:"Post_content"`   //客户姓名
	Post_title     string    `json:"Post_title"`     //客户姓名
	Post_status    string    `json:"Post_status"`    //客户姓名
	Comment_status string    `json:"Comment_status"` //客户姓名
	Ping_status    string    `json:"Ping_status"`    //客户姓名
	Post_password  string    `json:"Post_password"`  //客户姓名
	Post_name      string    `json:"Post_name"`      //客户姓名
	Post_parent    int       `json:"Post_parent"`    //客户ID
	Guid           string    `json:"Guid"`           //客户姓名
	Menu_order     int       `json:"Menu_order"`     //客户ID
	Post_type      string    `json:"Post_type"`      //客户姓名
	Comment_count  int       `json:"Comment_count"`  //客户ID
}

func GetWordpress(db *sql.DB) [](*Wordpress) {
	//编写查询语句
	//select 客户ID,客户姓名,月开始日期,月购买总金额,月购买次数 from dbo.CustomerYearlySalesReport
	strsql := "SELECT  `ID`,  `post_author`,  `post_date`,  `post_date_gmt`,  " +
		"LEFT(`post_content`, 256), LEFT(`post_title`, 256),  " +
		"`post_status`,  `comment_status`,  `ping_status`,  `post_password`,  `post_name`,  " +
		"`post_parent`,  `guid`,  `menu_order`, `post_type`, `comment_count`" +
		" FROM `wordpress`.`wp_posts`;"
	stmt, err := db.Prepare(strsql)
	if err != nil {
		log.Println("Prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return nil
	}
	//将数据读取到实体中
	var rowsData [](*Wordpress)
	for rows.Next() {
		data := new(Wordpress)

		//其中一个字段的信息 ， 如果要获取更多，就在后面增加：rows.Scan(&row.Name,&row.Id)
		rows.Scan(&data.ID, &data.Post_author, &data.Post_date, &data.Post_date_gmt,
			&data.Post_content, &data.Post_title, &data.Post_status, &data.Comment_status,
			&data.Ping_status, &data.Post_password, &data.Post_name,
			&data.Post_parent, &data.Guid, &data.Menu_order,
			&data.Post_type, &data.Comment_count)

		rowsData = append(rowsData, data)
	}
	return rowsData
}
