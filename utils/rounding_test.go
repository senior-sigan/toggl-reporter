package utils

import (
	"testing"
	"time"
)

func TestRoundTime(t *testing.T) {
	precision := time.Duration(5) * time.Minute
	{
		got := RoundTime(time.Duration(2)*time.Minute+time.Duration(3)*time.Second, precision)
		exp := time.Duration(5) * time.Minute
		if got != exp {
			t.Errorf("RoundTime exp=%v got=%v", exp, got)
		}
	}

	{
		got := RoundTime(time.Duration(10)*time.Minute, precision)
		exp := time.Duration(10) * time.Minute
		if got != exp {
			t.Errorf("RoundTime exp=%v got=%v", exp, got)
		}
	}

	{
		got := RoundTime(time.Duration(6)*time.Minute, precision)
		exp := time.Duration(10) * time.Minute
		if got != exp {
			t.Errorf("RoundTime exp=%v got=%v", exp, got)
		}
	}

	{
		got := RoundTime(time.Duration(0), precision)
		exp := time.Duration(0)
		if got != exp {
			t.Errorf("RoundTime exp=%v got=%v", exp, got)
		}
	}
}
