package taillog

import (
	"github.com/hpcloud/tail"
	"fmt"
)

var (
	tailObj *tail.Tail
)

func Init(file string) (err error) {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开文件
		Follow:    true,                                 //跟随文件
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	tailObj, err = tail.TailFile(file, config)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	return
}

func ReadChanel()<-chan *tail.Line{
	return tailObj.Lines
}
