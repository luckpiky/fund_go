package main

import (
	"container/list"
	"testing"
)

func TestGetFundDataFromSina(t *testing.T) {
	for _, uint := range [] struct {
		code string
		page int
		expected int
	}{
		{"519712", 1, 21},
		{"519712", 99999, 0},
	}{
		data := list.New()
		ret := GetFundDataFromSina(uint.code, uint.page, data)
		if ret != uint.expected {
			t.Errorf("GetFundDataFromSina: [%v], actually: [%v]", 0, ret)
		}		
	}
}


func BenchmarkGetFundDataFromSina(b *testing.B) {
    // b.N会根据函数的运行时间取一个合适的值
    for i := 0; i < b.N; i++ {
		data := list.New()
		GetFundDataFromSina("519712", 1, data)
    }
}