package dayjs

import (
	"reflect"
	"testing"
	"time"
)

func TestDayJS_Parse(t *testing.T) {
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
			d, err := Parse(tt.date)
			if err != nil {
				t.Fatal(err)
			}
			got, err := d.ToTime()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayJS_ParseFormat(t *testing.T) {
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
			d, err := ParseFormat(tt.date, tt.format)
			if err != nil {
				t.Fatal(err)
			}
			got, err := d.ToTime()
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayJS_FromTime_Format(t *testing.T) {
	date := time.Date(2022, 1, 2, 3, 4, 5, 0, time.Local)
	format := "YYYY-MM-DD HH:mm:ss"
	want := "2022-01-02 03:04:05"
	d, err := FromTime(date)
	if err != nil {
		t.Fatal(err)
	}
	err = d.TimeZone("America/Toronto")
	if err != nil {
		t.Fatal(err)
	}
	got, err := d.Format(format)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
