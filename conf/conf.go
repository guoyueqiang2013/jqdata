package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//定义配置文件解析后的结构
type JQConfig struct {
	DB_addr      string
	DB_port string
	DB_name string
	DB_username string
	DB_password        string
	TB_prefix       string
	Code string
	Unit string
	JQ_mob string
	JQ_pwd string
}

/*
"db_addr":"127.0.0.1", //db的地址
  "db_port":"3306", //db端口
  "db_username":"root",
  "db_password":"111111",
  "tb_prefix":"mk_", //数据表前缀
  "code":"600000.XSHG", //证券代码
  "unit":"1m",
  "jq_mob": "18538128837",
  "jq_pwd": "JQ111111aa"
 */

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}