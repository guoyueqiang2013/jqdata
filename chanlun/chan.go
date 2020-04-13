package chanlun

import (
	"container/list"
	"fmt"
)

type KCandle struct {
	Time uint32
	Open float32
	Close float32
	High float32
	Low float32
	isHandled bool //是否已经经过处理
	Parting int //0 未分型， -1 底分型  , 1 顶分型
}


type KLine struct {
	Level string   //1m , 5m , 30m , 1d ....
	Candles *list.List
	//FromTime uint32   //开始时间，格式举例：  2015年1月1日10点01分  -> 1501011001
	//ToTime uint32     //开始时间，格式举例：  2020年4月11日11点38分 -> 2004111138
}

func NewKLine(level string) *KLine {
	return &KLine{
		Level:level,
		//FromTime:ft,
		//ToTime:tt,
		Candles:list.New(),
	}
}




type Parting struct {

}

type ChanLun struct {
	//原始KLine
	SourceKLine *KLine

	//合并后的KLine
	MergedKLine *KLine

	//分型
	PartingKLine *KLine
}

func NewChanObject()*ChanLun  {
	ch := &ChanLun{
		SourceKLine:NewKLine("1m"),
		MergedKLine:NewKLine("1m"),
		PartingKLine:NewKLine("1m"),
	}
	return ch
}


func (c *ChanLun)AddKCandle(k KCandle)  {
	c.SourceKLine.Candles.PushBack(&k)
}


//处理数据
/*
①消除K线的包含（详细看处理原则）->
②分型（顶分型，底分型）->  顶分型区间，底分型区间 ->
③笔 ->
④线段（有缺口，无缺口） ->
⑤中枢 ->
⑥走势（盘整，趋势）->
⑦背驰 ->
⑧中枢的升级
 */
func (c *ChanLun)Handle()  {

	for e := c.SourceKLine.Candles.Front(); e != nil; e = e.Next() {
		cn, b := e.Value.(*KCandle)
		if b != true {
			fmt.Println("类型不匹配")
			continue
		}

		if cn.isHandled == true{
			continue
		}
		logCandle(cn,"[__Add__]")
		cn.isHandled = true
		handle_include(c.MergedKLine,cn)
		handle_parting(c.MergedKLine)
	}

	for e := c.MergedKLine.Candles.Front(); e != nil; e = e.Next() {
		cn, b := e.Value.(*KCandle)
		if b != true {
			fmt.Println("类型不匹配")
			continue
		}
		logCandle(cn,"")
	}

	fmt.Println("--------------")
}


//处理包含关系
func handle_include(mline *KLine, cn *KCandle)  {
	//var currCandle ,preCandle,nextCandle *KCandle
	if mline.Candles.Len()<=0 {
		mline.Candles.PushBack(cn)
		return
	}
	curr := mline.Candles.Back()
	currCandle , b := curr.Value.(*KCandle)
	if b != true {
		return
	}
	if  currCandle.High == cn.High && currCandle.Low == cn.Low{
		return
	}

	if !(currCandle.High >= cn.High && currCandle.Low <= cn.Low){
		mline.Candles.PushBack(cn)
		return
	}

	direct := 0
	pre := curr.Prev()
	if pre == nil{
		direct = 0
	}else{
		preCandle , b := pre.Value.(*KCandle)
		if b != true{
			direct = 0
		}else{
			if preCandle.High > currCandle.High{
				direct = 0
			}else{
				direct = 1
			}
		}
	}

	//如果方向向下，取两者高点中的低点，低点中的低点；如果方向向上，取两者高点中的高点，低点中的高点
	newCandle := &KCandle{}
	newCandle.isHandled = true
	if direct == 0{
		if currCandle.High >= cn.High{
			newCandle.High = cn.High
		}else{
			newCandle.High = currCandle.High
		}
		if currCandle.Low >= cn.Low{
			newCandle.Low = cn.Low
			newCandle.Time = cn.Time
		}else{
			newCandle.Low = currCandle.Low
			newCandle.Time = currCandle.Time
		}
	}else{
		if currCandle.High >= cn.High{
			newCandle.High = currCandle.High
			newCandle.Time = currCandle.Time
		}else{
			newCandle.High = cn.High
			newCandle.Time = cn.Time
		}
		if currCandle.Low >= cn.Low{
			newCandle.Low = currCandle.Low
		}else{
			newCandle.Low = cn.Low
		}
	}

	logCandle(currCandle,"[__Curr__]")
	mline.Candles.Remove(curr)
	mline.Candles.PushBack(newCandle)

}

func logCandle(candle *KCandle,tag string)  {
	if candle == nil{
		return
	}
	fmt.Printf("%sti:%v,open:%v,close:%v,high:%v,low:%v,isHandled:%v,parting:%d\n",
		tag,candle.Time,candle.Open,candle.Close,candle.High,candle.Low,candle.isHandled,candle.Parting)
}

/*分型
顶分型： 比两边 高点高，低点也高
低分型： 比两边 高点低，低点也低
 */
func handle_parting(mline *KLine)  {
	if nil == mline{
		return
	}

	i:=0
	for e := mline.Candles.Front(); e != nil; e = e.Next() {

		cobj, b := e.Value.(*KCandle)  //获得Candle对象
		if b != true {
			fmt.Println("类型不匹配")
			continue
		}

		cobj.Parting = 0
		i++
		if i <= 2{
			continue
		}

		right := e
		rightObj := cobj
		mid := right.Prev()
		midObj , b := mid.Value.(*KCandle)
		if b != true {
			fmt.Println("类型不匹配")
			continue
		}
		left := mid.Prev()
		leftObj , b := left.Value.(*KCandle)
		if b != true {
			fmt.Println("类型不匹配")
			continue
		}

		if midObj.Low < leftObj.Low  &&  midObj.Low < rightObj.Low && midObj.High < leftObj.High && midObj.High < rightObj.High{
			midObj.Parting = -1
		}
		if midObj.Low > leftObj.Low  &&  midObj.Low > rightObj.Low && midObj.High > leftObj.High && midObj.High > rightObj.High{
			if midObj.Parting == -1{
				panic("怎么可以同时两种分型？")
			}
			midObj.Parting = 1
		}

	}


}


//处理笔
func handle_bi(mline *KLine)  {
	
}




