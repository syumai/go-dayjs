package dayjs

import (
	_ "embed"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dop251/goja"
)

//go:embed assets/dayjs.min.js
var dayjsCode string

var (
	ErrInitializeDayJSVM = fmt.Errorf("failed to initialize DayJS VM")
	ErrParseDate         = fmt.Errorf("failed to parse date")
	ErrUnexpectedValue   = fmt.Errorf("unexpected type of value returned")
)

type DayJS struct {
	vm        *goja.Runtime
	dayjsFunc goja.Callable
	dayjsObj  *goja.Object
	mu        sync.Mutex
}

func newInstance() (*DayJS, error) {
	vm := goja.New()
	if _, err := vm.RunString(dayjsCode); err != nil {
		return nil, err
	}
	dayjsFuncValue := vm.Get("dayjs")
	dayjsFunc, ok := goja.AssertFunction(dayjsFuncValue)
	if !ok {
		return nil, ErrInitializeDayJSVM
	}
	return &DayJS{
		vm:        vm,
		dayjsFunc: dayjsFunc,
	}, nil
}

func Parse(date string) (*DayJS, error) {
	d, err := newInstance()
	if err != nil {
		return nil, err
	}
	v, err := d.dayjsFunc(nil, d.vm.ToValue(date)) // dayjs(date)
	if err != nil {
		return nil, err
	}
	d.dayjsObj = v.ToObject(d.vm)
	return d, nil
}

func ParseFormat(date, format string) (*DayJS, error) {
	d, err := newInstance()
	if err != nil {
		return nil, err
	}
	v, err := d.dayjsFunc(nil, d.vm.ToValue(date), d.vm.ToValue(format))
	if err != nil {
		return nil, err
	}
	d.dayjsObj = v.ToObject(d.vm)
	return d, nil
}

func FromTime(t time.Time) (*DayJS, error) {
	d, err := newInstance()
	if err != nil {
		return nil, err
	}
	milliNum := t.UnixMilli()
	v, err := d.dayjsFunc(nil, d.vm.ToValue(milliNum))
	if err != nil {
		return nil, err
	}
	d.dayjsObj = v.ToObject(d.vm)
	return d, nil
}

func (d *DayJS) Format(format string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	v := d.dayjsObj.Get("format")
	formatFunc, ok := goja.AssertFunction(v)
	if !ok {
		return "", ErrUnexpectedValue
	}
	strValue, err := formatFunc(d.dayjsObj, d.vm.ToValue(format))
	if err != nil {
		return "", err
	}
	resultStr := strValue.String()
	return resultStr, nil
}

func (d *DayJS) ToTime() (time.Time, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	v := d.dayjsObj.Get("toDate")
	toDateFunc, ok := goja.AssertFunction(v)
	if !ok {
		return time.Time{}, ErrUnexpectedValue
	}
	dateValue, err := toDateFunc(d.dayjsObj)
	if err != nil {
		return time.Time{}, err
	}
	resultDate := dateValue.ToFloat()
	if math.IsNaN(resultDate) {
		return time.Time{}, ErrParseDate
	}
	return time.UnixMilli(int64(resultDate)), nil
}
