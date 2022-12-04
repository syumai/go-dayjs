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

type DayJS struct {
	vm     *goja.Runtime
	global *goja.Object
	free   func()
	mu     sync.Mutex
}

func New() (*DayJS, error) {
	vm := goja.New()
	_, err := vm.RunString(dayjsCode)
	if err != nil {
		return nil, err
	}
	return &DayJS{
		vm:     vm,
		global: vm.GlobalObject(),
	}, nil
}

// this method must be called in locked method.
func (d *DayJS) clearGlobal(name string) error {
	_, err := d.vm.RunString(fmt.Sprintf("delete globalThis.%s", name))
	if err != nil {
		return err
	}
	return nil
}

var ErrParseDate = fmt.Errorf("failed to parse date")

func (d *DayJS) Parse(date string) (time.Time, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	err := d.global.Set("date", date)
	if err != nil {
		return time.Time{}, err
	}
	defer d.clearGlobal("date")
	result, err := d.vm.RunString("Number(dayjs(date).toDate())")
	if err != nil {
		return time.Time{}, err
	}
	resultDate := result.ToFloat()
	if math.IsNaN(resultDate) {
		return time.Time{}, ErrParseDate
	}
	return time.UnixMilli(int64(resultDate)), nil
}

func (d *DayJS) ParseFormat(date, format string) (time.Time, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.global.Set("date", date); err != nil {
		return time.Time{}, err
	}
	defer d.clearGlobal("date")
	if err := d.global.Set("format", format); err != nil {
		return time.Time{}, err
	}
	defer d.clearGlobal("format")
	result, err := d.vm.RunString("Number(dayjs(date, format).toDate())")
	if err != nil {
		return time.Time{}, err
	}
	resultDate := result.ToFloat()
	if math.IsNaN(resultDate) {
		return time.Time{}, ErrParseDate
	}
	return time.UnixMilli(int64(resultDate)), nil
}

func (d *DayJS) Format(t time.Time, format string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.global.Set("date", t.UnixMilli()); err != nil {
		return "", err
	}
	defer d.clearGlobal("date")
	if err := d.global.Set("format", format); err != nil {
		return "", err
	}
	defer d.clearGlobal("format")
	result, err := d.vm.RunString("dayjs(date).format(format)")
	if err != nil {
		return "", err
	}
	resultStr := result.String()
	return resultStr, nil
}
