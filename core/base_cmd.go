package core

type BaseCmd struct {
	ReadBufferSize int `clop:"long" usage:"read buffer size" default:"1024"`
	// 使用限制端口范围, 默认1， -1表示不限制
	LimitPortRange int  `clop:"short;long" usage:"limit port range" default:"1"`
	Reuse          bool `clop:"short;long" usage:"reuse port"`
}
