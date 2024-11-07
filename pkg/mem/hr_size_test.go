package mem_test

import (
	"testing"

	"github.com/yanodincov/skyeng-ics/pkg/mem"
)

func TestGetHumanReadableSize(t *testing.T) {
	t.Parallel()

	type args struct {
		bytes int
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "1 byte",
			args: args{bytes: 1},
			want: "1 B",
		},
		{
			name: "1 kilobyte",
			args: args{bytes: 1024},
			want: "1.0 KB",
		},
		{
			name: "1 megabyte",
			args: args{bytes: 1024 * 1024},
			want: "1.0 MB",
		},
		{
			name: "1 gigabyte",
			args: args{bytes: 1024 * 1024 * 1024},
			want: "1.0 GB",
		},
		{
			name: "1.1 gigabyte",
			args: args{bytes: 1024*1024*1024 + 1024*1024*1024/10},
			want: "1.1 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := mem.GetHumanReadableSize(tt.args.bytes); got != tt.want {
				t.Errorf("GetHumanReadableSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
