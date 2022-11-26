package dayjs

import (
	"reflect"
	"testing"
	"time"
)

func TestDayJS_Parse(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Free()
	tests := map[string]struct {
		date    string
		want    time.Time
		wantErr bool
	}{
		"valid": {
			date:    "2022-01-02 03:04:05",
			want:    time.Date(2022, 1, 2, 3, 4, 5, 0, time.Local),
			wantErr: false,
		},
		"valid with only year-month-day": {
			date:    "2022-01-02",
			want:    time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
		"invalid": {
			date:    "abcde",
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := d.Parse(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayJS_ParseFormat(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Free()
	tests := map[string]struct {
		date    string
		format  string
		want    time.Time
		wantErr bool
	}{
		"valid": {
			date:    "03:04:05 2022-01-02",
			format:  "HH:mm:ss YYYY-MM-DD",
			want:    time.Date(2022, 1, 2, 3, 4, 5, 0, time.Local),
			wantErr: false,
		},
		"valid with only year-month-day": {
			date:    "02-01-2022",
			format:  "DD-MM-YYYY",
			want:    time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
		"invalid": {
			date:    "abcde",
			format:  "YYYY-MM-DD",
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := d.ParseFormat(tt.date, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayJS_Format(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Free()
	date := time.Date(2022, 1, 2, 3, 4, 5, 0, time.Local)
	format := "YYYY-MM-DD HH:mm:ss"
	want := "2022-01-02 03:04:05"
	got, err := d.Format(date, format)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
