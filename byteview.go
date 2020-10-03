package GeeCache

//缓存值的抽象与封装
type ByteView struct {
	b []byte //缓存值
}

//返回视图长度
func (v ByteView) Len() int {
	return len(v.b)
}

//返回数据副本作为切片
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//以字符串形式返回数据,必要时返回复制
func (v ByteView) String() string {
	return string(v.b)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
