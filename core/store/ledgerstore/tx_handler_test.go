

package ledgerstore

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestSyncMapRange(t *testing.T) {
	m := sync.Map{}

	for i := 0; i < 10; i++ {
		m.Store("k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
	}
	cnt := 0

	m.Range(func(key, value interface{}) bool {
		fmt.Printf("key :%s, val :%s\n", key, value)

		if key == "k5" {
			return false
		}
		cnt += 1
		return true
	})

	fmt.Println(cnt)

}


func add(m map[string]int, va int) {
	m["key"] = va
}

func TestSyncMapRW(t *testing.T) {
	m := &sync.Map{}
	m.Store("key", 10)
	for i := 0; i < 100000; i++ {
		go addsync(m, i)
	}

}

func addsync(m *sync.Map, va int) {
	m.Store("key", va)
}
