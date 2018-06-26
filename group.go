package foog

import (
	"errors"
	"sync"
)

//定义Group
type Group struct {
	mutex   sync.RWMutex
	once    sync.Once
	members map[int64]*Session
}

//新建
var (
	errMemberNotFound = errors.New("member not found")
)

//new group
func NewGroup() *Group {
	group := &Group{
		members: make(map[int64]*Session),
	}
	return group
}

//加入session
func (this *Group) Join(sess *Session) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.members[sess.Id] = sess
}

//删除session
func (this *Group) Leave(sess *Session) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.members, sess.Id)
}

//清除所有的session
func (this *Group) Clean() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.members = make(map[int64]*Session)
}

//是否存在
func (this *Group) Member(sid int64) (*Session, error) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	mem, ok := this.members[sid]
	if ok {
		return mem, nil
	}

	return nil, errMemberNotFound
}

//返回所有的session
func (this *Group) Members() []*Session {
	members := []*Session{}
	for _, v := range this.members {
		members = append(members, v)
	}
	return members
}

//广播所有的session
func (this *Group) Broadcast(msg interface{}) {
	for _, v := range this.members {
		v.WriteMessage(msg)
	}
}

//通过合适的session写入数据
func (this *Group) BroadcastWithoutSession(msg interface{}, filters map[int64]bool) {
	for k, v := range this.members {
		if _, ok := filters[k]; !ok {
			v.WriteMessage(msg)
		}
	}
}
