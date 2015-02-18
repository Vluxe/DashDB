package dash

import (
	"fmt"
	"os"
	"sync"
)

type Dash struct {
	store  map[string]string
	mutex  *sync.Mutex
	wQueue chan pair
	dbFile string
}

type pair struct {
	key   string
	value string
}

//create a new Dash instance
func New() (*Dash, error) {
	d := Dash{store: make(map[string]string), mutex: new(sync.Mutex), wQueue: make(chan pair), dbFile: "dash.db"}
	err := d.startQueue()
	return &d, err
}

//add a value to the store
func (d *Dash) Set(key, value string) {
	d.mutex.Lock()
	d.store[key] = value
	d.mutex.Unlock()
	d.wQueue <- pair{key: key, value: value}
}

//get a value from the store
func (d *Dash) Get(key string) string {
	return d.store[key]
}

//private things related to disk persistence

//start the disk writing queue
func (d *Dash) startQueue() error {
	f, err := os.OpenFile(d.dbFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	go d.runQueue(f)
	return nil
}

//the running queue
func (d *Dash) runQueue(f *os.File) {
	defer f.Close()
	for {
		select {
		case pair := <-d.wQueue:
			f.Write([]byte(fmt.Sprintf("%d%s%d%s", len(pair.key), pair.key, len(pair.value), pair.value)))
		}
	}
}

//loads the data stored in the db file
func (d *Dash) loadData() error {
	f, err := os.OpenFile(d.dbFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	// b := make([]byte, 1)
	// f.Read(b)
	// count := Int(b)
	for {
		//read each key & value out in this loop
	}
	return nil
}
