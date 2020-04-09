package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	MOB = "18538128837"
	PWD = "JQ111111aa"
)

type Client struct {
	mob string
	pwd string
	token string
}

func request(body map[string]interface{}) string {
	url := "https://dataapi.joinquant.com/apis"
	bodyStr, err := json.Marshal(body)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyStr)))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(res)
}


func NewClint() *Client {
	return &Client{MOB, PWD,""}
}

func (cli *Client) GetToken() (token string , err error) {
	method := "get_token"
	body := map[string]interface{}{
		"method": method,
		"mob": cli.mob,
		"pwd": cli.pwd,
	}
	token = request(body)
	cli.token = token
	return
}

func (cli *Client)GetCurrentToken()(token string , err error)   {
	method := "get_current_token"
	body := map[string]interface{}{
		"method": method,
		"mob": cli.mob,
		"pwd": cli.pwd,
	}
	token = request(body)
	cli.token = token
	return
}

func (cli *Client)GetAllSecurities (code string,date string) string {
	method := "get_all_securities"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"code": code,
		"date": date,
	}
	ret := request(body)
	return ret
}

/*
获取一个指数给定日期在平台可交易的成分股列表
*/

func (cli *Client)GetIndexStocks (code string,date string) string {
	method := "get_index_stocks"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"code": code,
		"date": date,
	}
	ret := request(body)
	return ret
}

/*
 获取融资标的列表
*/

func (cli *Client)GetMargincashStocks  (date string) string {
	method := "get_margincash_stocks"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"date": date,
	}
	ret := request(body)
	return ret
}

/*

 - 获取融券标的列表
*/


func (cli *Client)GetMarginsecStocks  (date string) string {
	method := "get_marginsec_stocks"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"date": date,
	}
	ret := request(body)
	return ret
}

/*
get_current_tick - 获取最新的 tick 数据
688333.XSHG
688357.XSHG
688358.XSHG
688363.XSHG
688366.XSHG
688368.XSHG
*/

func (cli *Client)GetCurrentTick(code string) string {
	method := "get_current_tick"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"code": code,
	}
	ret := request(body)
	return ret
}

/*

get_price_period / get_bars_period- 获取指定时间段的行情数据
code: 证券代码
unit: bar的时间单位, 支持如下周期：1m, 5m, 15m, 30m, 60m, 120m, 1d, 1w, 1M。其中m表示分钟，d表示天，w表示周，M表示月
date : 开始时间，不能为空，格式2018-07-03或2018-07-03 10:40:00，如果是2018-07-03则默认为2018-07-03 00:00:00
end_date：结束时间，不能为空，格式2018-07-03或2018-07-03 10:40:00，如果是2018-07-03则默认为2018-07-03 23:59:00
fq_ref_date：复权基准日期，该参数为空时返回不复权数据
*/

func (cli *Client)GetPricePeriod(code , unit,date,end_date,fq_ref_date string) string {
	method := "get_price_period"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"code": code,
		"unit":unit,
		"date":date,
		"end_date":end_date,
		"fq_ref_date":fq_ref_date,
	}
	ret := request(body)
	return ret
}


/*
get_price / get_bars - 获取指定时间周期的行情数据
code: 证券代码
count: 大于0的整数，表示获取bar的条数，不能超过5000
unit: bar的时间单位, 支持如下周期：1m, 5m, 15m, 30m, 60m, 120m, 1d, 1w, 1M。其中m表示分钟，d表示天，w表示周，M表示月
end_date：查询的截止时间，默认是今天
fq_ref_date：复权基准日期，该参数为空时返回不复权数据
*/

func (cli *Client)GetPrice(code ,count,unit,end_date,fq_ref_date string) string {
	method := "get_price"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"code": code,
		"count":count,
		"unit":unit,
		"end_date":end_date,
		"fq_ref_date":fq_ref_date,
	}
	ret := request(body)
	return ret
}

//get_trade_days
/*
获取指定日期范围内的所有交易日

参数：

date: 开始日期
end_date: 结束日期

*/

func (cli *Client)GetTradeDays(date,end_date string)[]string  {
	method := "get_trade_days"
	body := map[string]interface{}{
		"method": method,
		"token":cli.token,
		"date": date,
		"end_date":end_date,
	}
	ret := request(body)
	return strings.Split(ret,"\n")
}


//将数据库中的日期数据转化成 字符串类型  例如： 510271010 => 2005-10-27 10:10:00
func ChangeDateFormat_I_S(old int32 )( new string)  {
	i64 := int64(old) + 200000000000
	new = strconv.FormatInt(int64(i64),10)
	//2005-10271010
	new = new[:4]+"-"+new[4:]
	//2005-10-271010
	new = new[:7]+"-"+new[7:]
	//2005-10-27 1010
	new = new[:10]+" "+new[10:]
	//2005-10-27 10:10
	new = new[:13]+":"+new[13:]
	//2005-10-27 10:10:00
	new = new+":00"
	return
}

//将 字符串类型 的日期转化成 数据库中的日期数据  2005-10-27 10:10 => 510271010
func ChangeDateFormat_S_I(old string)(new int32)  {
	old = strings.Replace(old,"-","",-1)
	old = strings.Replace(old," ","",-1)
	old = strings.Replace(old,":","",-1)
	if 12 != len(old){
		return 0
	}

	i64 ,err := strconv.Atoi(old)
	if err != nil{
		return 0
	}

	new = int32(i64 - 200000000000)

	return
}

//将 字符串类型 的日期转化成 数据库中的日期数据  2005-10-27 => 510270000
func ChangeDateFormat_S_I_2(old string)(new int32)  {
	old = strings.Replace(old,"-","",-1)
	if 8 != len(old){
		return 0
	}
	i64 ,err := strconv.Atoi(old)
	if err != nil{
		return 0
	}
	i64 = i64 * 10000
	new = int32(i64 - 200000000000)
	return
}

//将一个Time 对象 修改成日期字符串  例如：  time => 2005-01-01
func ChangeDateFormt_T_S(t time.Time) ( new string) {

	var month,day string

	if t.Month() >=10{
		month = fmt.Sprintf("%d",t.Month())
	}else{
		month = fmt.Sprintf("0%d",t.Month())
	}

	if t.Day() >=10{
		day = fmt.Sprintf("%d",t.Day())
	}else{
		day = fmt.Sprintf("0%d",t.Day())
	}

	new = fmt.Sprintf("%d-%s-%s",t.Year(),month,day)

	return
}