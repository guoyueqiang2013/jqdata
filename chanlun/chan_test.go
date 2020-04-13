package chanlun

import (
	"math/rand"
	"testing"
)

func TestChan(t *testing.T) {

	var ti uint32 = 1501011001
	var open,close,high,low , yesterday float32

	ti = ti + uint32(rand.Intn(9)+1)
	open = 100.0
	close = 95.5
	high = 101.0
	low = 92.6
	yesterday = close

	ch := NewChanObject()
	for i:=0; i<50;i++{
		ti = ti + uint32(rand.Intn(10)+1)
		yesterday = close
		open = yesterday + changePrice(yesterday) * float32(updown())
		close = yesterday + changePrice(yesterday) * float32(updown())
		if open >= close{
			high = open + changePrice(yesterday)
			low = close - (yesterday * rand.Float32()*0.1)
		}else{
			high = close + (yesterday * rand.Float32()*0.1)
			low = open - (yesterday * rand.Float32()*0.1)
		}

		//fmt.Printf("ti:%v,open:%v,close:%v,high:%v,low:%v\n",ti,open,close,high,low)

		ch.AddKCandle(KCandle{
			Time:ti,
			Open:open,
			Close:close,
			High:high,
			Low:low,
		})
		ch.Handle()
	}

}

func changePrice(yesterday float32)float32  {
	ret := (yesterday *  float32(rand.Intn(11))/100.0 )
	return ret
}

func updown()int  {
	if rand.Intn(2) <= 0{
		return -1
	}else{
		return 1
	}
}