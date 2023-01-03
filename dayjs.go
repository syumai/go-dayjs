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

//go:embed assets/datetime-polyfill.min.js
var dateTimeFormatPolyfill string

var (
	ErrInitializeDayJSVM = fmt.Errorf("failed to initialize DayJS VM")
	ErrParseDate         = fmt.Errorf("failed to parse date")
	ErrUnexpectedValue   = fmt.Errorf("unexpected type of value returned")
)

type DayJS struct {
	vm        *goja.Runtime
	dayjsFunc goja.Callable
	dayjsVal  goja.Value
	mu        sync.Mutex
}

func newInstance() (*DayJS, error) {
	vm := goja.New()
	if _, err := vm.RunString("globalThis.window = globalThis;"); err != nil {
		return nil, err
	}
	if _, err := vm.RunString(dateTimeFormatPolyfill); err != nil {
		return nil, err
	}
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
	d.dayjsVal = v
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
	d.dayjsVal = v
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
	d.dayjsVal = v
	return d, nil
}

// TimeZone sets time zone (e.g. "Asia/Tokyo") to current DayJS instance.
// See: https://day.js.org/docs/en/timezone/timezone
func (d *DayJS) TimeZone(timeZone string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	v := d.dayjsVal.ToObject(d.vm).Get("tz")
	tzFunc, ok := goja.AssertFunction(v)
	if !ok {
		return ErrUnexpectedValue
	}
	dayjsVal, err := tzFunc(d.dayjsVal, d.vm.ToValue(timeZone))
	if err != nil {
		return err
	}
	d.dayjsVal = dayjsVal
	return nil
}

func (d *DayJS) Format(format string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	v := d.dayjsVal.ToObject(d.vm).Get("format")
	formatFunc, ok := goja.AssertFunction(v)
	if !ok {
		return "", ErrUnexpectedValue
	}
	strValue, err := formatFunc(d.dayjsVal, d.vm.ToValue(format))
	if err != nil {
		return "", err
	}
	resultStr := strValue.String()
	return resultStr, nil
}

func (d *DayJS) ToTime() (time.Time, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	v := d.dayjsVal.ToObject(d.vm).Get("toDate")
	toDateFunc, ok := goja.AssertFunction(v)
	if !ok {
		return time.Time{}, ErrUnexpectedValue
	}
	dateValue, err := toDateFunc(d.dayjsVal)
	if err != nil {
		return time.Time{}, err
	}
	resultDate := dateValue.ToFloat()
	if math.IsNaN(resultDate) {
		return time.Time{}, ErrParseDate
	}
	return time.UnixMilli(int64(resultDate)), nil
}
