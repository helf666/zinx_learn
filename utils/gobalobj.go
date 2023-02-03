package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

//存储一切zinx有关的的全局参数
//一些参数可以通过zinx.json由用户进行配置

type GlobalObj struct {
	TcpServer ziface.IServer //当前zinx全局的server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int
	Name      string //服务器名称

	Version        string //zinx的版本号
	MaxConn        int    //最大的连接数
	MaxPackageSize uint32 //数据报最大数量

}

// 定义一个全局的global对象
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 初始化globalobject对象
func init() {
	//如果配置没有加载，则执行默认初始化
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "127.0.0.1",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}
