package cache

import (
	"errors"
	"strings"
	"sync"
	"time"
)

const (
	string_t = iota
	map_t
	list_t
)

type ListNode struct {
	val  string
	next *ListNode
}

type List struct {
	head   *ListNode
	last   *ListNode
	length int
}

func (v *List) printList() string {
	var accum strings.Builder
	for next := v.head; next != nil; next = next.next {
		accum.Write([]byte(next.val + "\n"))
	}
	return accum.String()
}

func (v *List) pushFirst(val string) {
	head := &ListNode{val, v.head}
	v.head = head
	v.length++
}

func (v *List) pushLast(val string) {
	last := &ListNode{val, nil}
	v.last.next = last
	v.last = last
	v.length++
}

func (v *List) get(ind int) string {
	if v.length < ind {
		return "Index out bounds"
	}
	i, n := 0, v.head
	for i < ind {
		i, n = i+1, n.next
	}
	return n.val
}

func (v *List) set(ind int, value string) error {
	if v.length < ind {
		return errors.New("Index out bounds")
	}
	i, n := 0, v.head
	for i < ind {
		i, n = i+1, n.next
	}
	n.val = value
	return nil
}

type CacheValue struct {
	ExpireTime int64
	Created    time.Time
	Type       int
	Value      interface{}
}

func (v *CacheValue) toString() string {
	switch v.Value.(type) {
	case List:
		head := v.Value.(List)
		return head.printList()
	case map[string]string:
		asMap := v.Value.(map[string]string)
		var accum strings.Builder
		for k, v := range asMap {
			accum.Write([]byte(k + ":" + v))
		}
		return accum.String()
	default:
		return v.Value.(string)
	}
}

func (v CacheValue) timeLeft() int {
	endDur := time.Duration(v.ExpireTime) * time.Second
	return int(v.Created.Add(endDur).Sub(time.Now()).Seconds())
}

type Cache struct {
	sync.Mutex
	values    map[string]CacheValue
	cancelTTL map[string]chan bool
}

func (c *Cache) Get(key string) (string, error) {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		return "", errors.New("No such key!")
	}
	if val.Type != string_t {
		return "", errors.New("Wrong key-method.")
	}
	return val.Value.(string), nil
}

func (c *Cache) Expire(key string, expire int64) string {
	val, ok := c.values[key]
	if ok {
		c.values[key] = CacheValue{expire, time.Now(), val.Type, val.Value}
		go c.KeyCleaner(key, expire)
		c.cancelTTL[key] = make(chan bool)
		return "Ok!"
	}
	return "No such key!"
}

func (c *Cache) Set(key string, expTime int64, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.cancelTTL[key]; ok {
		c.cancelTTL[key] <- true
	}

	if expTime > 0 {
		go c.KeyCleaner(key, expTime)
		c.cancelTTL[key] = make(chan bool)
	}
	stored, ok := c.values[key]
	if ok && stored.Type != string_t {
		return errors.New("Wrong set-method.")
	}
	c.values[key] = CacheValue{expTime, time.Now(), string_t, value}
	return nil
}

func (c *Cache) HGet(key, field string) (string, error) {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		return "", errors.New("No such key!")
	}
	if val.Type != map_t {
		return "", errors.New("Wrong get-method.")
	}
	ret, ok := val.Value.(map[string]string)[field]
	if !ok {
		return "", errors.New("No such field!")
	}
	return ret, nil
}

func (c *Cache) HGetAll(key string) ([]string, error) {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		return nil, errors.New("No such key!")
	}
	if val.Type != map_t {
		return nil, errors.New("Wrong get-method.")
	}
	hmap := val.Value.(map[string]string)
	var accum []string
	for k, v := range hmap {
		accum = append(accum, k, v)
	}
	return accum, nil
}

func (c *Cache) HSet(key, field, value string) error {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		c.values[key] = CacheValue{0, time.Now(), map_t, map[string]string{field: value}}
		return nil
	} else {
		if val.Type != map_t {
			return errors.New("Wrong set-method.")
		}
		c.values[key].Value.(map[string]string)[field] = value
		return nil
	}
}

func (c *Cache) RPush(key, elem string) (int, error) {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		root := &ListNode{elem, nil}
		c.values[key] = CacheValue{
			0,
			time.Now(),
			list_t,
			List{root, root, 1},
		}
		return 1, nil
	} else {
		if val.Type != list_t {
			return 0, errors.New("Key is not for list-value.")
		}
		listSnap := c.values[key].Value.(List)
		listSnap.pushLast(elem)
		c.values[key] = CacheValue{0, time.Now(), list_t, listSnap}
		return listSnap.length, nil
	}
}

func (c *Cache) TTL(key string) int64 {
	v, ok := c.values[key]
	if !ok {
		return -2
	}
	if v.ExpireTime == 0 {
		return -1
	}
	return v.ExpireTime
}

func (c *Cache) LPush(key, elem string) (int, error) {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		root := &ListNode{elem, nil}
		c.values[key] = CacheValue{
			0,
			time.Now(),
			list_t,
			List{root, root, 1},
		}
		return 1, nil
	} else {
		if val.Type != list_t {
			return 0, errors.New("Key is not for list-value.")
		}
		listSnap := c.values[key].Value.(List)
		listSnap.pushFirst(elem)
		c.values[key] = CacheValue{0, time.Now(), list_t, listSnap}
		return listSnap.length, nil
	}
}

func (c *Cache) LGet(key string, ind int) (string, error) {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		return "", errors.New("No such key")
	}
	if val.Type != list_t {
		return "", errors.New("Key is not for list-value.")
	}
	listSnap := val.Value.(List)
	return listSnap.get(ind), nil
}

func (c *Cache) LSet(key, value string, ind int) error {
	c.Lock()
	defer c.Unlock()
	val, ok := c.values[key]
	if !ok {
		return errors.New("No such key!")
	}
	if val.Type != list_t {
		return errors.New("Key is not for list-value.")
	}
	listSnap := val.Value.(List)
	return listSnap.set(ind, value)
}

func (c *Cache) LRange(key string, start, end int) ([]string, error) {
	val, ok := c.values[key]
	if !ok {
		return nil, errors.New("No such key!")
	}
	if val.Type != list_t {
		return nil, errors.New("Key is not for list-value.")
	}
	listSnap := val.Value.(List)
	if end < 0 {
		end = listSnap.length + end
	}
	if end < 0 || start > end {
		return nil, errors.New("Wrong interval!")
	}
	var accum []string
	for i := start; i <= end; i++ {
		accum = append(accum, listSnap.get(i))
	}
	return accum, nil
}

func (c *Cache) KeyCleaner(key string, expireTime int64) {
	select {
	case <-time.After(time.Duration(expireTime) * time.Second):
		c.Delete(key)
		delete(c.cancelTTL, key)
		return
	case <-c.cancelTTL[key]:
		delete(c.cancelTTL, key)
		return
	}
}

func (c *Cache) Delete(key string) {
	delete(c.values, key)
}

func (c *Cache) Keys() []string {
	c.Lock()
	defer c.Unlock()
	keys := make([]string, 0, len(c.values))
	for k := range c.values {
		keys = append(keys, k)
	}
	return keys
}

func InitCache() *Cache {
	return &Cache{values: make(map[string]CacheValue), cancelTTL: make(map[string]chan bool)}
}
