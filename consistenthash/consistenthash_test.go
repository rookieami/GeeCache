package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6", "3", "2") //3个真实节点
	testCaces := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "3",
		"27": "2",
	}
	for k, v := range testCaces {
		if hash.Get(k) != v {
			t.Errorf("要求%s，应该产生%s", k, v)
		}
	}
	hash.Add("8")
	testCaces["27"] = "8"
	for k, v := range testCaces {
		if hash.Get(k) != v {
			t.Errorf("要求%s，应该产生%s", k, v)
		}
	}
}
