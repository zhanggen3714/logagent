package seek_local_ip
import (
	"fmt"
	"net"
	"strings"
)

//获取本机 ip地址 作为etcd的key的一部分，起到划分etcd配置的作用。
func SeekLocalIP()(ip string,err error){
	conn,err:=net.Dial("udp","8.8.8.8:80")
	if err!=nil{
		return
	}
	defer conn.Close()
	localIP:=conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("服务器IP",localIP.String())
	ip=strings.Split(localIP.IP.String(),":")[0]
	return
}

