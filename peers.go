package GeeCache

import pb "GeeCache/geecachepb"

type PeerPicker interface {
	//根据传入key选择相应节点PeerGetter
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//从对应group查找缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
