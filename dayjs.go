package dayjs

import (
	_ "embed"
	"fmt"
	"sync"
	"time"

	"github.com/elsaland/quickjs"
)

//go:embed assets/dayjs.min.js
var dayjsCode string

type DayJS struct {
	qjsRuntime *quickjs.Runtime
	qjsCtx     *quickjs.Context
	global     *quickjs.Value
	free       func()
	mu         sync.Mutex
}

func New() (*DayJS, error) {
	qjsRuntime := quickjs.NewRuntime()
	qjsCtx := qjsRuntime.NewContext()
	global := qjsCtx.Globals()
	free := func() {
		qjsCtx.Free()
		qjsRuntime.Free()
	}
	dayjsResult, err := qjsCtx.Eval(dayjsCode)
	if err != nil {
		free()
		return nil, err
	}
	dayjsResult.Free()
	var freeOnce sync.Once
	return &DayJS{
		qjsRuntime: &qjsRuntime,
		qjsCtx:     qjsCtx,
		global:     &global,
		free: func() {
			freeOnce.Do(free)
		},
	}, nil
}

func (d *DayJS) Free() {
	d.free()
}

// this method must be called in locked method.
func (d *DayJS) clearGlobal(name string) error {
	result, err := d.qjsCtx.Eval(fmt.Sprintf("delete globalThis.%s", name))
	if err != nil {
		return err
	}
	result.Free()
	return nil
}

func (d *DayJS) Parse(date string) (time.Time, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.global.Set("date", d.qjsCtx.String(date))
	defer d.clearGlobal("date")
	result, err := d.qjsCtx.Eval("Number(dayjs(date).toDate())")
	if err != nil {
		return time.Time{}, err
	}
	resultDate := result.Int64()
	result.Free()
	return time.UnixMilli(resultDate), nil
}

func (d *DayJS) ParseFormat(date, format string) (time.Time, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.global.Set("date", d.qjsCtx.String(date))
	defer d.clearGlobal("date")
	d.global.Set("format", d.qjsCtx.String(format))
	defer d.clearGlobal("format")
	result, err := d.qjsCtx.Eval("Number(dayjs(date, format).toDate())")
	if err != nil {
		return time.Time{}, err
	}
	resultDate := result.Int64()
	result.Free()
	return time.UnixMilli(resultDate), nil
}

func (d *DayJS) Format(t time.Time, format string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.global.Set("date", d.qjsCtx.Int64(t.UnixMilli()))
	defer d.clearGlobal("date")
	d.global.Set("format", d.qjsCtx.String(format))
	defer d.clearGlobal("format")
	result, err := d.qjsCtx.Eval("dayjs(date).format(format)")
	if err != nil {
		return "", err
	}
	resultStr := result.String()
	result.Free()
	return resultStr, nil
}