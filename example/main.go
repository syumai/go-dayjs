package main

import (
	"fmt"
	"log"
	"time"

	"github.com/syumai/go-dayjs"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	djs, err := dayjs.New()
	if err != nil {
		return err
	}
	defer djs.Free()

	{
		result, err := djs.Parse("2021-01-02")
		if err != nil {
			return err
		}
		fmt.Println(result)
	}

	{
		result, err := djs.ParseFormat("02-01-2020", "DD-MM-YYYY")
		if err != nil {
			return err
		}
		fmt.Println(result)
	}

	{
		now := time.Now()
		result, err := djs.Format(now, "YYYY-MM-DD HH:mm:ss")
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}
