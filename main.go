package main

import (
	"TestProject/crawler"
	"fmt"
	"time"
)

func main() {
    c := &crawler.Fetcher{}
	c.Init()

    coo, _ := c.GetCookie()

	for  {
		c1 := &crawler.Fetcher{}
		c1.Init()
		time.Sleep(30 * time.Second )
		fmt.Println(1111)
		c1.GetXq(coo)
	}

}
