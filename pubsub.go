package main

import (
	"sync"
	"time"
)

var wgPool = sync.Pool{New: func() interface{} { return new(sync.WaitGroup) }}

func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

type subscriber chan interface{}
type topicFunc func(v interface{}) bool

type Publisher struct {
	m           sync.RWMutex
	buffer      int
	timeout     time.Duration
	subscribers map[subscriber]topicFunc
}

func (p *Publisher) Len() int {
	p.m.RLock()
	i := len(p.subscribers)
	p.m.RUnlock()
	return i
}

func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
}

func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	delete(p.subscribers, sub)
	close(sub)
	p.m.Unlock()
}

func (p *Publisher) Publish(v interface{}) {
	p.m.RLock()
	if len(p.subscribers) == 0 {
		p.m.RUnlock()
		return
	}

	wg := wgPool.Get().(*sync.WaitGroup)
	for sub, topic := range p.subscribers {
		if topic != nil && !topic(v) {
			continue
		}
		wg.Add(1)
		go p.sendTopic(sub, v, wg)
	}
	wg.Wait()
	wgPool.Put(wg)
	p.m.RUnlock()
}

func (p *Publisher) Close() {
	p.m.Lock()
	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
	p.m.Unlock()
}

func (p *Publisher) sendTopic(sub subscriber, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	if p.timeout > 0 {
		select {
		case sub <- v:
		case <-time.After(p.timeout):
		}
		return
	}

	select {
	case sub <- v:
	default:
	}
}
