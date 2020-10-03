package lru

import (
	"container/list"
)

//lru 缓存淘汰策略

type Cache struct {
	maxBytes  int64                         //允许使用最大内存
	nbytes    int64                         //当前已使用内存
	ll        *list.List                    //双向链表
	cache     map[string]*list.Element      //字典,key是字符串,value是链表指针
	onEvicted func(key string, value Value) //记录移除回调函数
}

//键值对,双向链表的数据类型
type entry struct {
	key   string
	value Value //内存大小
}

//返回值所占用内存
type Value interface {
	Len() int
}

//构造函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

//查找键的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //移动到队尾
		//返回查找值
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//删除,移除最近最少访问的节点,队首
func (c *Cache) RemoveOldest() {
	//取得队首节点
	ele := c.ll.Back()
	if ele != nil {
		//移除
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//从字典中删除节点映射关系
		delete(c.cache, kv.key)
		//更新当前使用内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}

}

//新增/修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//键存在,更新键,节点移动到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		//新申请内存-已有内存=需加内存
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//新增
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		//使用内存超过了最大值,移除最少访问的节点
		c.RemoveOldest()
	}
}

//缓存条目数
func (c *Cache) Len() int {
	return c.ll.Len()
}
