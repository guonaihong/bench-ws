package core

type BaseCmd struct {
	ReadBufferSize int  `clop:"long" usage:"read buffer size" default:"1024"`
	Reuse          bool `clop:"short;long" usage:"reuse port"`
}
