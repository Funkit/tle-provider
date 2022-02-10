package data

import (
	"testing"
)

// TestGetDayOfYear Test parsing of date and display of the day of year fraction
func TestGetDayOfYear(t *testing.T) {

	date1 := "2022-02-09T12:04:21.971712"

	output, err := getDayOfYear(date1)
	if err != nil {
		t.Errorf("FAIL: parsing date string 1")
	}

	t.Log(output)
}

// TestGetLast2DigitsOfYear Test parsing of last 2 digits of date and display
func TestGetLast2DigitsOfYear(t *testing.T) {

	date1 := "2022-02-09T12:04:21.971712"

	output, err := getLast2DigitsOfYear(date1)
	if err != nil {
		t.Errorf("FAIL: parsing last 2 digits of date string 1")
	}

	t.Log(output)
}
