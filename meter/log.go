package meter

import (
	"io"
	"log"
)

func InitLogger(f io.Writer) {

	// 设置日志输出到文件
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
