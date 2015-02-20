package dash

import (
	//"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
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
	d.loadData()
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
			f.Write([]byte(fmt.Sprintf("%d\n%s%d\n%s", len(pair.key), pair.key, len(pair.value), pair.value)))
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
	offset := 0
	sliceStart := 0
	key := ""
	for {
		buffer, err, size := readData(f)
		if err != nil {
			return err
		}
		if size <= 0 {
			return nil
		}
		for offset < size {
			if buffer[offset] == '\n' {
				count := buffer[sliceStart:offset]
				num, err := strconv.Atoi(count)
				if err != nil {
					return err
				}
				//fmt.Println("num is:", num)
				offset++
				//check to see if the buffer needs to be expanded
				for offset+num > size {
					b, err, s := readData(f)
					if err != nil {
						return err
					}
					size = s + (size - offset)
					buffer = buffer[offset:] + b
					offset = 0
				}
				data := buffer[offset:(offset + num)]
				//fmt.Println("data is:", data)
				offset += num
				if key == "" {
					key = data
				} else {
					// fmt.Println("key:", key)
					// fmt.Println("value:", data)
					d.store[key] = data
					key = ""
				}
				sliceStart = offset
			}
			offset++
		}
	}
	return nil
}

func readData(f *os.File) (string, error, int) {
	b := make([]byte, 2048)
	size, err := f.Read(b)
	if err != nil && err != io.EOF {
		return "", err, 0
	}
	if size <= 0 {
		return "", nil, size
	}
	return string(b), nil, size
}
