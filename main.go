package main

import (
	"time"

	"github.com/unix2dos/zj-business/business"
)

func main() {

	bs := business.New()
	bs.Run()

	for {
		time.Sleep(time.Hour * 10)
	}
}
