package main

import (
	"fmt"
	"github.com/vluxe/DashDB"
	"math/rand"
	//"os"
	"time"
)

func main() {
	d, err := dash.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	start := time.Now()
	defer done(d, start) //d.Cleanup()
	runTask(d, "with Sync")
	//fmt.Println("doing some random work for 4 seconds")
	//time.Sleep(time.Second * 4)
	//fmt.Println("done with random work. Waiting on cleanup")

	// os.Remove("dash.db")
	// d.DoSync = false //insanely fast
	// runTask(d, "without sync")
	// os.Remove("dash.db")
	// loadVal := d.Get("name")
	// fmt.Println("loaded value is:", loadVal)

	// loaded := d.Get("name")
	// fmt.Println("disk value is:", loaded)

	// d.Set("name", "Dalton")
	// val := d.Get("name")
	// fmt.Println("value is:", val)

	// d.Remove("name")
	// fmt.Println("remove value")

	// d.Set("name", "Austin")
	// v := d.Get("name")
	// fmt.Println("value is:", v)

	//d.Set("name", "Long\nName")
	//t := d.Get("name")
	//fmt.Println("value is:", t)
	//time.Sleep(time.Millisecond * 1000)
}

func done(d *dash.Dash, start time.Time) {
	d.Cleanup()
	timeTrack(start, "overall")
}

func runTask(d *dash.Dash, name string) {
	defer timeTrack(time.Now(), name)
	count := 0
	seed := 100000
	keySize := 10
	valSize := 100
	for count < seed {
		d.Set(randSeq(keySize), randSeq(valSize))
		count++
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
