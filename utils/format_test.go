package utils

import (
	"testing"
	"time"
)

func TestHours(t *testing.T) {
	d := time.Duration(32)*time.Second + time.Duration(78)*time.Minute + time.Duration(25)*time.Hour
	h := Hours(d)
	if h != 26 {
		t.Errorf("Expected 26 but got %v", h)
	}
}

func TestMinutes(t *testing.T) {
	d := time.Duration(32)*time.Second + time.Duration(78)*time.Minute + time.Duration(25)*time.Hour
	m := Minutes(d)
	if m != 18 {
		t.Errorf("Expected 18 but got %v", m)
	}
}

func TestSeconds(t *testing.T) {
	d := time.Duration(32)*time.Second + time.Duration(78)*time.Minute + time.Duration(25)*time.Hour
	s := Seconds(d)
	if s != 32 {
		t.Errorf("Expected 32 but got %v", s)
	}
}

func TestFormatDuration_1(t *testing.T) {
	d := time.Duration(32)*time.Second + time.Duration(78)*time.Minute + time.Duration(25)*time.Hour
	str := FormatDuration(d)
	if str != "26:18:32" {
		t.Errorf("Expected '26:18:32', got '%v'", str)
	}
}

func TestFormatDuration_2(t *testing.T) {
	d := time.Duration(2)*time.Second + time.Duration(15)*time.Minute + time.Duration(1)*time.Hour
	str := FormatDuration(d)
	if str != "01:15:02" {
		t.Errorf("Expected '01:15:02', got '%v'", str)
	}
}
