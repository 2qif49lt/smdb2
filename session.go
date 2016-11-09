package main

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const (
	validtime       = time.Minute * 10
	session_id_size = 32
)

type session struct {
	alive time.Time
	data  map[string]interface{}
}

type ssnMgr struct {
	rwlock *sync.RWMutex
	ssns   map[string]*session
}

func newssnMgr() *ssnMgr {
	return &ssnMgr{
		&sync.RWMutex{},
		make(map[string]*session),
	}
}

func (mgr *ssnMgr) NewId() string {
	id := getRandString(session_id_size)

	mgr.rwlock.Lock()
	defer mgr.rwlock.Unlock()

	mgr.ssns[id] = &session{
		time.Now(),
		make(map[string]interface{}),
	}

	return id
}

func getRandString(size int) string {
	p := make([]byte, size)
	rand.Read(p)

	return hex.EncodeToString(p)
}

func (mgr *ssnMgr) Set(id, key string, val interface{}) {
	mgr.rwlock.Lock()
	defer mgr.rwlock.Unlock()

	if s, exist := mgr.ssns[id]; exist {
		s.alive = time.Now()
		s.data[key] = val
	}
}

func (mgr *ssnMgr) IsExist(id string) bool {
	mgr.rwlock.RLock()
	defer mgr.rwlock.RUnlock()

	_, exist := mgr.ssns[id]
	return exist
}

func (mgr *ssnMgr) Del(id string) {
	mgr.rwlock.Lock()
	defer mgr.rwlock.Unlock()

	if _, exist := mgr.ssns[id]; exist {
		delete(mgr.ssns, id)
	}
}

func (mgr *ssnMgr) Get(id, key string) interface{} {
	mgr.rwlock.RLock()
	defer mgr.rwlock.RUnlock()

	if s, exist := mgr.ssns[id]; exist {
		s.alive = time.Now()
		if val, exist := s.data[key]; exist {
			return val
		}
	}
	return nil
}

func (mgr *ssnMgr) ClearTimeout() {
	mgr.rwlock.Lock()
	defer mgr.rwlock.Unlock()

	for k, v := range mgr.ssns {
		if time.Since(v.alive) > validtime {
			delete(mgr.ssns, k)
		}
	}
}

func sessionAliveRountion() {
	for {
		time.Sleep(time.Minute)
		ssnmgr.ClearTimeout()
	}
}
