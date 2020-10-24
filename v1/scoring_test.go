package main

import "testing"

func TestAccuracy(t *testing.T) {
	type args struct {
		correctCharacters int
		typedCharacters   int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "ShouldReturn0.5",
			args: args{correctCharacters: 20, typedCharacters: 40},
			want: 0.5,
		},
		{
			name: "ShouldReturn0.7428",
			args: args{correctCharacters: 52, typedCharacters: 70},
			want: 0.7428571428571429,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Accuracy(tt.args.correctCharacters, tt.args.typedCharacters); got != tt.want {
				t.Errorf("Accuracy() = %v, want %v", got, tt.want)
			}
		})
	}
}
