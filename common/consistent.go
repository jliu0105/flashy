package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// define the type of new slice
type units []uint32

// return the length of the slice
func (x units) Len() int {
	return len(x)
}

// compare two numbers
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

// swap 2 values of the slice
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("No data in Hash")

type Consistent struct {
	// hash circle, key is hash value, storing the info
	circle map[uint32]string
	//Already sorted node hash slice
	sortedHashes units
	// virtual node to add the balance of hash
	VirtualNode int
	//map read write lock
	sync.RWMutex
}

// set default node number
func NewConsistent() *Consistent {
	return &Consistent{
		// initialize variable
		circle: make(map[uint32]string),
		// set virtual node
		VirtualNode: 20,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

// get the location of hash key
func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {
		var srcatch [64]byte
		// copy the data into the array
		copy(srcatch[:], key)
		//Returns the CRC-32 checksum of the data using an IEEE polynomial
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// update sorting so easy to check
func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	// Determine whether the slice capacity is too large, and reset if it is too large
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	//add hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	// sort all hashes so binary search can be easier
	sort.Sort(hashes)
	// rename
	c.sortedHashes = hashes

}

// add node in hash circle
func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

//add node
func (c *Consistent) add(element string) {
	//Loop virtual nodes, set replicas
	for i := 0; i < c.VirtualNode; i++ {
		//Add to the hash ring according to the generated nodes
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	// update sort
	c.updateSortedHashes()
}

//remove node
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

// remove one node
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// check the closest node clock wise
func (c *Consistent) search(key uint32) int {
	// check algo
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	// use binary search to get the min val
	i := sort.Search(len(c.sortedHashes), f)
	// if exceed the range then set i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

// Obtain the nearest server node information according to the data label
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	// if 0 then return error
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	// calculate hash value
	key := c.hashkey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil

}
