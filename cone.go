package cone

import (
	"sync"
)

// cone, basically the struct to do our work
type Cone struct {
	dbPath     string //file path to our DB file
	wg         *sync.WaitGroup
	shouldWait bool
	mutex      *sync.Mutex
}

//creates a new Cone object
func New(dbPath string) *Cone {
	return Cone{dbPath: dbPath, wg: &sync.WaitGroup, mutex: &sync.Mutex}
}

//Get a value from the store
func (c *Cone) Get(key string) []byte {
	c.checkWait()
	c.wg.Add(1)
	defer c.wg.Done()
	//do work
	return nil
}

//add a value to the store
func (c *Cone) Set(key string, value []byte) error {
	c.updateWait(true)
	c.checkWait()
	c.wg.Add(1)
	defer c.wg.Done()
	//do work
	c.updateWait(false)
	return nil
}

func (c *Cone) updateWait(status bool) {
	c.mutex.Lock()
	c.shouldWait = status
	c.mutex.Unlock()
}

//check if the method should wait
func (c *Cone) checkWait() {
	if c.shouldWait {
		c.wg.Wait() //waits on any pending opts to finish.
	}
}
