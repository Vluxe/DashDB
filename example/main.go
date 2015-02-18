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
	d.Set("name", "Dalton")
	val := d.Get("name")
	fmt.Println("value is: ", val)

	//time.Sleep(time.Millisecond * 1000)
	d.Set("name", "Austin")
	v := d.Get("name")
	fmt.Println("value is: ", v)
	time.Sleep(time.Millisecond * 1000)
}
