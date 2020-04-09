package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guoyueqiang2013/jqdata/conf"
	"github.com/guoyueqiang2013/jqdata/sdk"
	"strings"
	"time"
)

func main() {

	/*
	date := "2020-04-01"

	cli := sdk.NewClint()
	token ,err := cli.GetCurrentToken()
	if err != nil{
		println(err)
	}
	println(token)

	//println(cli.GetCurrentTick("000002.XSHE"))

	//code ,count,unit,end_date,fq_ref_date
	println(cli.GetPrice("000002.XSHE","120","1m",date,""))

	*/

	//读取配置文件
	JsonParse := conf.NewJsonStruct()
	v:=conf.JQConfig{}
	JsonParse.Load("config.json", &v)
	fmt.Println(v)

	//连接数据库
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",v.DB_username, v.DB_password, "tcp", v.DB_addr, v.DB_port, v.DB_name)
	fmt.Println(conn)
	DB, err := sql.Open("mysql", conn)
	defer DB.Close()
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}
	DB.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超时的连接就close
	DB.SetMaxOpenConns(100)                //设置最大连接数

	tn:=fmt.Sprintf("%s%s_%s",v.TB_prefix,v.Unit,v.Code)
	tn = strings.Replace(tn,".","_",-1)
	fmt.Println(tn)
	CreateTable(DB,tn)

	//读取 聚宽 的数据

	maxTime := QueryMaxDate(DB,tn)
	fmt.Printf("max time => %d , %v\n",maxTime,sdk.ChangeDateFormat_I_S(maxTime))


	ReadJQData(v,maxTime)

	//date,open,close,high,low,volume,money



	//println(sdk.ChangeDateFormat_I_S(510271010))
	//println(sdk.ChangeDateFormat_I_S(2010271010))
	//println(sdk.ChangeDateFormat_S_I("2005-10-27 10:10"))
	//println(sdk.ChangeDateFormat_S_I("2010-10-27 10:10"))


}

func ReadJQData(v conf.JQConfig,maxTime int32)(back []string)  {

	fromDateTime := time.Date(2005, time.January, 1, 0, 0, 0, 0 , time.Local)
	fromDate := sdk.ChangeDateFormt_T_S(fromDateTime)
	yTime := time.Now().AddDate(0, 0, -8)
	yesterdayDate := sdk.ChangeDateFormt_T_S(yTime)
	todayDate := sdk.ChangeDateFormt_T_S(time.Now())
	fmt.Printf("fromDate=> %s ,yesterdayDate=>%s ,todayDate => %s\n",fromDate,yesterdayDate,todayDate)

	//startDateTime := sdk.ChangeDateFormat_S_I_2(fromDate)


/*
	cli := sdk.NewClint()
	token ,err := cli.GetCurrentToken()
	if err != nil{
		println(err)
	}
	println(token)

	days :=cli.GetTradeDays(fromDate,yesterdayDate)
	for _,k := range days{
		println(sdk.ChangeDateFormat_S_I_2(k))
	}
*/

	/*
	backstr := cli.GetPricePeriod(v.Code,v.Unit,fromDate,toDate,"2005-01-01")
	back = strings.Split(backstr,"\n")

	for i := range back{
		println(back[i])
	}
	*/
	return
}


type TB struct {
	Time int32       	`json:"time" form:"time"`
	Money int        	`json:"money" form:"money"`
	Open float32    	`json:"open" form:"open"`
	Close float32       `json:"close" form:"close"`
	High float32        `json:"high" form:"high"`
	Low float32         `json:"low" form:"low"`
	Volume float32      `json:"volume" form:"volume"`
	High_limit float32  `json:"high_limit" form:"high_limit"`
	Low_limit float32   `json:"low_limit" form:"low_limit"`
}


func CreateTable(DB *sql.DB,tablename string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(" +
		"time int(11) unsigned zerofill NOT NULL," +
		"open float unsigned zerofill NOT NULL," +
		"close float unsigned zerofill NOT NULL," +
		"high float unsigned zerofill NOT NULL,"+
		"low float unsigned zerofill NOT NULL,"+
		"volume int(10) unsigned zerofill NOT NULL,"+
		"money int(11) unsigned zerofill NOT NULL,"+
		"high_limit float unsigned zerofill NOT NULL,"+
		"low_limit float unsigned zerofill NOT NULL,"+
		"PRIMARY KEY (time)"+
		");",tablename)

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}

//插入数据
func InsertData(DB *sql.DB,tn string , tb *TB) {
	fmt.Sprintf("INSERT INTO `%s` VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v');",
		tn,tb.Time,tb.Open,tb.Close,tb.High,tb.Low,tb.Volume,tb.Money,tb.High_limit,tb.Low_limit)
	_,err := DB.Exec("insert INTO users(username,password) values(?,?)","test","123456")
	if err != nil{
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}
}

//选择最大日期的数据
func QueryMaxDate(DB *sql.DB,tn string) (maxTime int32) {
	//SELECT * FROM tablename ORDER BY alexa DESC;
	//tb := new(TB)
	maxTime = 0
	sql := fmt.Sprintf("SELECT MAX(time) FROM %s;",tn)
	println(sql)
	row,err := DB.Query(sql)
	if err != nil{
		return maxTime
	}
	for row.Next(){
		er := row.Scan(&maxTime)
		if er != nil{
			return maxTime
		}
		println(maxTime)
	}
	return
}

/*
//查询单行
func QueryOne(DB *sql.DB,tn string) {
	tb := new(TB)   //用new()函数初始化一个结构体对象
	sql := fmt.Sprintf("select time,open,close,high,low,volume,money,high_limit,low_limit from %s ORDER BY time DESC",tn)
	row := DB.QueryRow(, 1)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	if err := row.Scan(&user.Id,&user.Username,&user.Password); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Println("Single row data:", *user)
}
*/