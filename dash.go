package dash

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	WriteCode  = 1 // write a value to the store
	RemoveCode = 2 // remove a value from the store
	CloseCode  = 3 // use to notify the queues to close cleanly
)

const (
	WriteAction  = "ADD" // add a value to the store
	RemoveAction = "DEL" // remove a value from the store
)

//represents the dash database interaction.
type Dash struct {
	store     map[string]string //the map used for the key/value opts.
	workQueue chan pair         //channel for queuing writes to the internal map
	fileQueue chan pair         //channel used for queuing database file writes
	dbFile    string            //the location of the database file. Default is ./dash.db
	DoSync    bool              //should it be fast or safe?
	maxDBSize float64           //the max size of the db file in MB. Default is 100MB. Be careful of making this size to small as it is possible to get caught in a infinite loop of shrinking the db.
}

//used to communicate actions across the channels
type pair struct {
	key   string
	value string
	code  int
}

//create a new Dash instance
func New() (*Dash, error) {
	d := Dash{store: make(map[string]string), workQueue: make(chan pair), fileQueue: make(chan pair), dbFile: "dash.db", maxDBSize: 100}
	d.DoSync = true
	d.loadData()
	go d.start()
	err := d.startDiskQueue()
	return &d, err
}

//add a value to the store
func (d *Dash) Set(key, value string) {
	d.workQueue <- pair{key: key, value: value, code: WriteCode}
}

//get a value from the store
func (d *Dash) Get(key string) string {
	return d.store[key]
}

//remove a value from the store
func (d *Dash) Remove(key string) {
	d.workQueue <- pair{key: key, value: "", code: RemoveCode}
}

//waits for all pending opts to complete and cleans up the file handles
func (d *Dash) Cleanup() {
	d.workQueue <- pair{key: "", value: "", code: CloseCode}
	d.fileQueue <- pair{key: "", value: "", code: CloseCode}
}

//private things related to disk persistence and key value processing

//start the main processing queue
func (d *Dash) start() {
	for {
		select {
		case pair := <-d.workQueue:
			if pair.code == CloseCode {
				return //the database is done being used
			}
			if pair.code == WriteCode {
				d.store[pair.key] = pair.value
			} else {
				delete(d.store, pair.key)
			}
			go d.writePair(pair)
		}
	}
}

//write an opt to the file queue
func (d *Dash) writePair(p pair) {
	d.fileQueue <- p
}

//opens the db file then starts the disk queue
func (d *Dash) startDiskQueue() error {
	f, err := d.openDBFile()
	if err != nil {
		return err
	}
	go d.runDiskQueue(f)
	return nil

}

//channel that waits and preforms file opts that are written to it.
func (d *Dash) runDiskQueue(f *os.File) {
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()
	for {
		select {
		case pair := <-d.fileQueue:
			if pair.code == CloseCode {
				return //the database is done being used
			}
			action := WriteAction
			if pair.code == RemoveCode {
				action = RemoveAction
			}
			b := buildDiskAction(action, pair.key, pair.value)
			f.Write(b)
			size += int64(len(b))
			//we can and should optimize this. e.g: Time interval (1 second) to do a sync, how many opts are pending (only sync every 100 opts when load is high),
			//buffer the written content, etc. allow for varying levels of safety vs speed.
			if d.DoSync {
				f.Sync()
			}
			//check if the db size hash been exceeded
			mbSize := float64(size) * 0.000001
			if mbSize > d.maxDBSize {
				//transfer the content to the temp.db
				tempPath := "temp.db"
				tempFile, err := os.Create(tempPath)
				if err != nil {
					//not sure what to do here but log
					fmt.Println("unable to create temp db file. That is a very serious error:", err)
				} else {
					for key, val := range d.store {
						b := buildDiskAction(WriteAction, key, val)
						tempFile.Write(b)
					}
					tempFile.Sync()
					os.Remove(d.dbFile)
					os.Rename(tempPath, d.dbFile)
					tempFile.Close()
					f.Close()
					f, _ = d.openDBFile()
					info, _ := f.Stat()
					size = info.Size()
				}
			}
		}
	}
}

//builds a disk action for writing
func buildDiskAction(action, key, value string) []byte {
	return []byte(fmt.Sprintf("%d\n%s%d\n%s%d\n%s", len(action), action, len(key), key, len(value), value))
}

func (d *Dash) openDBFile() (*os.File, error) {
	return os.OpenFile(d.dbFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

//loads the data stored in the db file into the store.
func (d *Dash) loadData() error {
	f, err := d.openDBFile()
	if err != nil {
		return err
	}
	defer f.Close()
	//should probably add some more validations to ensure the file hasn't be tampered with (although that would be very uncommon).
	offset := 0
	sliceStart := 0
	key := ""
	action := ""
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
				if action == "" {
					action = data
				} else if key == "" {
					key = data
				} else {
					// fmt.Println("action:", action)
					// fmt.Println("key:", key)
					// fmt.Println("value:", data)
					if action == RemoveAction {
						delete(d.store, key)
					} else {
						d.store[key] = data
					}
					action = ""
					key = ""
				}
				sliceStart = offset
			}
			offset++
		}
	}
	return nil
}

//helper method for loadData(). Reads data from the db file.
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
