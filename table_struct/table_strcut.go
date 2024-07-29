package tablestruct

//这是数据库对应的表结构

//库存表对应dbo.View_KC
type Inventory struct {
	ProductID   string  //产品类型
	ProductName string  //产品名字
	ResidualNum float64 //库存数量
}
