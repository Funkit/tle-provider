package data

import (
	"testing"
)

func TestGetDayOfYear(t *testing.T) {

	correctDate := "2022-02-09T12:04:21.971712"

	output, err := getDayOfYear(correctDate)
	if err != nil {
		t.Errorf("FAIL: parsing correct date string")
	}

	t.Log(output)

	incorrectDate := "2022/02/09 12:04:21"
	output, err = getDayOfYear(incorrectDate)
	if err == nil {
		t.Errorf("FAIL: parsing incorrect date string")
	}

	t.Log(output)
}

func TestGetLast2DigitsOfYear(t *testing.T) {

	date1 := "2022-02-09T12:04:21.971712"

	output, err := getLast2DigitsOfYear(date1)
	if err != nil {
		t.Errorf("FAIL: parsing last 2 digits of date string 1")
	}

	t.Log(output)
}

func TestObjectIDToCOSPARID(t *testing.T) {
	objectID := "1964-063C"
	cosparID := "64063C  " //ID with trailing spaces

	output, err := objectIDToCOSPARID(objectID)
	if err != nil {
		t.Errorf("FAIL: parsing object ID")
	}

	t.Log(output)
	t.Log(len(output))

	if output != cosparID {
		t.Errorf("FAIL: COSPAR ID different than expected")
	}
}

func TestFormatWithoutDecimalPoint(t *testing.T) {
	bStar1 := 0.11693
	expectedOutput := " 11693+0" //value with leading spaces

	bStar1Formatted := formatWithoutDecimalPoint(bStar1)

	t.Log(bStar1Formatted)

	if bStar1Formatted != expectedOutput {
		t.Errorf("FAIL: formatting does not match expected output")
	}

	bStar2 := 19.352
	expectedOutput = " 19352+2" //value with leading spaces

	bStar2Formatted := formatWithoutDecimalPoint(bStar2)

	t.Log(bStar2Formatted)

	if bStar2Formatted != expectedOutput {
		t.Errorf("FAIL: formatting does not match expected output")
	}

}
