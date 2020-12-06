package main

import (
	"testing"
	"time"
)

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

func TestCpm(t *testing.T) {
	type args struct {
		correctCharacters int
		duration          time.Duration
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "ShouldReturnCPMGivenDurationInSeconds",
			args: args{
				correctCharacters: 20,
				duration:          time.Second * 15,
			},
			want: 80,
		},
		{
			name: "ShouldReturnCPMGivenDurationInMinutes",
			args: args{
				correctCharacters: 80,
				duration:          time.Minute,
			},
			want: 80,
		},
		{
			name: "ShouldReturnCPMGivenDurationInMixedUnits",
			args: args{
				correctCharacters: 80,
				duration:          time.Minute + time.Second * 20,
			},
			want: 60,
		},

		{
			name: "ShouldReturnCPMWithoutRoundingOff",
			args: args{
				correctCharacters: 78,
				duration:          time.Second * 35,
			},
			want: 133.71428571428572,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cpm(tt.args.correctCharacters, tt.args.duration); got != tt.want {
				t.Errorf("Cpm() = %v, want %v", got, tt.want)
			}
		})
	}
}
