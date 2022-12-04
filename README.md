# go-dayjs

* `go-dayjs` is a wrapper of [day.js](https://day.js.org/) running on QuickJS.
  - This package is using [goja](https://github.com/dop251/goja).

## Installation

```
go get github.com/syumai/go-dayjs
```

## Usage

```go
djs, _ := dayjs.New()
defer djs.Free()

{
  result, _ := djs.Parse("2022-01-25")
  fmt.Println(result) // 2021-01-02 00:00:00 +0900 JST
}

{
  result, _ = djs.ParseFormat("03-01-2020", "DD-MM-YYYY")
  fmt.Println(result) // 2020-01-02 00:00:00 +0900 JST
}

{
  now := time.Now()
  result, _ = djs.Format(now, "YYYY-MM-DD HH:mm:ss")
  fmt.Println(result) // 2022-11-26 23:57:21
}
```

## License

MIT

## Author

* syumai
* [iamkun](https://github.com/iamkun) (original author of day.js)
