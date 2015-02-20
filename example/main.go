package main

import (
	"fmt"
	"github.com/vluxe/DashDB"
	"time"
)

func main() {
	d, err := dash.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	loadVal := d.Get("name")
	fmt.Println("loaded value is:", loadVal)

	d.Set("name", "Dalton")
	val := d.Get("name")
	fmt.Println("value is:", val)

	//time.Sleep(time.Millisecond * 1000)
	d.Set("name", "Austin")
	v := d.Get("name")
	fmt.Println("value is:", v)

	d.Set("name", "Long\nName")
	t := d.Get("name")
	fmt.Println("value is:", t)
	time.Sleep(time.Millisecond * 1000)
}
