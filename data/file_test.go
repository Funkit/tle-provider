package data

import (
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestFileSource_extractSatelliteData(t *testing.T) {
	type fields struct {
		filePath           string
		TwoLineElements    []Satellite
		TwoLineElementsMap map[string]Satellite
		LastCelestrakPull  time.Time
		UpdatePeriod       float64
		mu                 sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Satellite
		wantErr bool
	}{
		{
			name: "basic case",
			fields: fields{
				filePath:           "../samples/tle.txt",
				TwoLineElements:    []Satellite{},
				TwoLineElementsMap: map[string]Satellite{},
				LastCelestrakPull:  time.Time{},
				UpdatePeriod:       0,
				mu:                 sync.RWMutex{},
			},
			want: []Satellite{
				{
					SatelliteName: "OPS 5712 (P/L 153)",
					NORADID:       2874,
					TLELine1:      "1 02874U 67053H   22206.60472723 -.00000017  00000-0  26447-4 0  9991",
					TLELine2:      "2 02874  69.9738 283.4261 0009834 250.7192 109.2850 13.96410943808158",
				},
				{
					SatelliteName: "CALSPHERE 1",
					NORADID:       900,
					TLELine1:      "1 00900U 64063C   22206.83199285  .00000371  00000-0  38562-3 0  9993",
					TLELine2:      "2 00900  90.1732  41.6116 0024844 266.8448 104.5887 13.73849434875933",
				},
				{
					SatelliteName: "LAGEOS 1",
					NORADID:       8820,
					TLELine1:      "1 08820U 76039A   22206.68532073  .00000028  00000-0  00000-0 0  9999",
					TLELine2:      "2 08820 109.8533  52.0899 0045094 246.5947 308.4924  6.38664901822297",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FileSource{
				filePath:           tt.fields.filePath,
				TwoLineElements:    tt.fields.TwoLineElements,
				TwoLineElementsMap: tt.fields.TwoLineElementsMap,
				UpdatePeriod:       tt.fields.UpdatePeriod,
			}
			got, err := fs.extractSatelliteData()
			if (err != nil) != tt.wantErr {
				t.Errorf("extractSatelliteData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].NORADID > got[j].NORADID
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].NORADID > tt.want[j].NORADID
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractSatelliteData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeTrailingSpaces(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default case",
			args: args{
				source: "CALSPHERE 1             ",
			},
			want: "CALSPHERE 1",
		},
		{
			name: "empty string",
			args: args{
				source: "",
			},
			want: "",
		},
		{
			name: "only spaces",
			args: args{
				source: "         ",
			},
			want: "",
		},
		{
			name: "one character",
			args: args{
				source: "a",
			},
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeTrailingSpaces(tt.args.source); got != tt.want {
				t.Errorf("removeTrailingSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
