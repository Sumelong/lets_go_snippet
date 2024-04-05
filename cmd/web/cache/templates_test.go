package cache

import (
	"testing"
	"time"
)

func Test_humanDate(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Local",
			args: args{t: time.Date(2022, 12, 2, 10, 0, 0, 0, time.Local)},
			want: "02 Dec 2022 at 10:00",
		},
		{
			name: "Empty",
			args: args{t: time.Time{}},
			want: "",
		},
		{
			name: "UTC",
			args: args{t: time.Date(2022, 12, 2, 10, 0, 0, 0, time.UTC)},
			want: "02 Dec 2022 at 10:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanDate(tt.args.t); got != tt.want {
				t.Errorf("humanDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
