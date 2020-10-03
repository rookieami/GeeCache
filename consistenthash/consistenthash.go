package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//字节映射到uint32
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //hash函数
	replicas int            //虚拟节点倍数
	keys     []int          //哈希环
	hashMap  map[int]string //虚拟节点与真实节点的映射表
}

//构造函数 ，允许自定义虚拟节点倍数与hash函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//添加真实节点/机器的add()方法

func (m *Map) Add(keys ...string) {
	//传入多个真实节点的名称
	for _, key := range keys {
		//每个真实节点创建m.replicas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			//虚拟节点的名称为strconv.Itoa(i)+key
			//m.hash计算虚拟节点的hash值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//添加到环上
			m.keys = append(m.keys, hash)
			//在hashMap中添加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	//环上的hash值排序
	sort.Ints(m.keys)
}

//获取最接近键值的项
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		//环为空
		return ""
	}
	//计算键值对应hash
	hash := int(m.hash([]byte(key)))
	//二分搜索
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//通过hashMap映射到真实的节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
