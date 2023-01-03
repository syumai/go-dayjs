# go-dayjs

* `go-dayjs` is a wrapper of [day.js](https://day.js.org/).
  - This package is using [goja](https://github.com/dop251/goja).

## Installation

```
go get github.com/syumai/go-dayjs
```

## Usage

```go
{
  d, _ := dayjs.Parse("2021-01-02")
  result, _ := d.ToTime()
  fmt.Println(result) // 2021-01-02 00:00:00 +0900 JST
}

{
  d, _ := dayjs.ParseFormat("02-01-2020", "DD-MM-YYYY")
  result, _ := d.ToTime()
  fmt.Println(result) // 2021-01-02 00:00:00 +0900 JST
}

{
  now := time.Now()
  d, _ := dayjs.FromTime(now)
  result, _ := d.Format("YYYY-MM-DD HH:mm:ss")
  fmt.Println(result) // 2022-12-04 21:08:55
}

{
  now := time.Now()
  d, _ := dayjs.FromTime(now)
  result, _ := d.Format("X")
  fmt.Println(result) // 1670155735
}
```

* There is no timezone support.
* If you want to use specific time zone, please overwrite your environment's `time.Local`.

```go
{
  l, _ := time.LoadLocation("Asia/Tokyo")
  time.Local = l
  now := time.Now()
  d, _ := dayjs.FromTime(now)
  ...
}
```

## License

MIT

## Author

* syumai
* [iamkun](https://github.com/iamkun) (original author of day.js)
