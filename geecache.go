package GeeCache

import (
	pb "GeeCache/geecachepb"
	"GeeCache/siinglefight"
	"fmt"
	"log"
	"sync"
)

//负责与外部交互,控制缓存的存储与获取流程

//缓存名称空间,分布着相关加载数据
type Group struct {
	name      string
	getter    Getter     //缓存未命中时获取数据源的回调
	mainCache cache      //并发缓存
	peers     PeerPicker //一致性哈希算法的map,根据key找到节点
	loader    *siinglefight.Group
}

//回调函数,缓存不存在时,调用该函数得到源数据
type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(key string) ([]byte, error)

//回调函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

//实例化group,将group存储到全局变量groups
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &siinglefight.Group{},
	}
	groups[name] = g
	return g
}

//获得指定名称的group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

//从缓存中获得键值对
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		//缓存存在,从mainCache 中查找缓存
		log.Println("hit")
		return v, nil
	}
	//不存在,调用load方法,load调用getLocally调用回调函数g.getter.Get()获取源数据
	return g.load(key)
}

//将实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}
func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil { //节点不为空
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		//获取源数据
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

//添加数据到缓存中
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	//添加到缓存中
	g.populateCache(key, value)
	return value, nil
}

//从访问远程节点，获取缓存值
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
