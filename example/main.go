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
	{
		d, err := dayjs.Parse("2021-01-02")
		if err != nil {
			return err
		}
		result, err := d.ToTime()
		if err != nil {
			return err
		}
		fmt.Println(result)
	}

	{
		d, err := dayjs.ParseFormat("02-01-2020", "DD-MM-YYYY")
		if err != nil {
			return err
		}
		result, err := d.ToTime()
		if err != nil {
			return err
		}
		fmt.Println(result)
	}

	{
		now := time.Now()
		d, err := dayjs.FromTime(now)
		if err != nil {
			return err
		}

		if err := d.TimeZone("Asia/Tokyo"); err != nil {
			return err
		}
		result, err := d.Format("YYYY-MM-DD HH:mm:ss")
		if err != nil {
			return err
		}
		fmt.Println(result)

		if err := d.TimeZone("Europe/London"); err != nil {
			return err
		}
		result, err = d.Format("YYYY-MM-DD HH:mm:ss")
		if err != nil {
			return err
		}
		fmt.Println(result)
	}

	{
		now := time.Now()
		d, err := dayjs.FromTime(now)
		if err != nil {
			return err
		}
		result, err := d.Format("X")
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}
