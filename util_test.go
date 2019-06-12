package main

import "testing"

func TestTruncateWithLine(t *testing.T) {
	type args struct {
		str      string
		maxLines int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "maxLines < 1",
			args: args{
				str:      "a\nb\nc\n",
				maxLines: 0,
			},
			want: "a\nb\nc\n",
		},
		{
			name: "lines <= maxLines",
			args: args{
				str:      "a\nb\n",
				maxLines: 3,
			},
			want: "a\nb\n",
		},
		{
			name: "lines > maxLines",
			args: args{
				str:      "a\nb\nc\nd\n",
				maxLines: 3,
			},
			want: "a\nb\nc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateWithLine(tt.args.str, tt.args.maxLines); got != tt.want {
				t.Errorf("TruncateWithLine() = %v, want %v", got, tt.want)
			}
		})
	}
}