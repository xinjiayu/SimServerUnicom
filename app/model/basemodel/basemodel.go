package basemodel


type BaseModel struct {
	Createtime int     `json:"createtime"`
	Updatetime int     `json:"updatetime"`
}
// 定义model interface
type IModel interface {
	// 获取表名
	TableName() string
	// 获取主键值
	PkVal() int
}
