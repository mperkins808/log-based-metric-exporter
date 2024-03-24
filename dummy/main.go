package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0
	for {
		l := fmt.Sprintf("log entry %v", i)
		fmt.Println(l)

		if i > 10000 {
			i = 0
		}
		i += 1

		time.Sleep(time.Millisecond * 50)
	}
}
