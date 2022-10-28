package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println("Location : ", t.Location(), " Time minute : ", t.Minute(), " Time Hour : ", t.Hour(), " Time Seconds : ", t.Second()) // local time
	fmt.Println("Location : ", t.Location(), t.Hour(), t.Minute(), t.Second())                                                         // local time

	// location, err := time.LoadLocation("America/New_York")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Location : ", location, " Time : ", t.In(location)) // America/New_York

	// loc, _ := time.LoadLocation("Asia/Kolkata")
	// now := time.Now().In(loc)
	// fmt.Println("Location : ", loc, " Time : ", now) // Asia/Shanghai
}
