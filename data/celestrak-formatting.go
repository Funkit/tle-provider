package data

import (
	"fmt"
	"math"
)

func formatMeanMotionDOT(meanDOT float64) string {
	meanDOTString := fmt.Sprintf("%.8f", math.Abs(meanDOT))
	if meanDOT < 0 {
		return "-" + meanDOTString[1:]
	}

	return " " + meanDOTString[1:]
}

func formatWithoutDecimalPoint(value float64) string {

	// Compute exponent for scientific notation
	exponent := int(math.Floor(math.Log10(math.Abs(value))))

	// If |exponent| is > 9, it is negligible so it should return 00000-0
	if exponent > 9 || exponent < -9 {
		return " 00000-0"
	}

	exponentAbs := int(math.Abs(float64(exponent)))
	// Remove exponent (transform 1.2e-5 into 1.2) then multiply by 10000 to have 5 figures to display (transform 1.2 into 12000)
	varwithoutexp := value * math.Pow10(exponentAbs) * math.Pow10(4)

	// If varwithoutexp < 1, it means the result is negligible hence the 0 output
	if varwithoutexp < 1 {
		return " 00000-0"
	}

	var2string := fmt.Sprintf("%5d", int(varwithoutexp))

	// Substract 1 as it is not proper scientific notion like 1.2e-4 but 0.12e-3
	output := " " + var2string + fmt.Sprintf("%+d", exponent+1)

	return output
}

func formatAngles(val float64) string {
	leadingSpaces := ""

	if val < 100 {
		leadingSpaces += " "
		if val < 10 {
			leadingSpaces += " "
		}

	}

	return leadingSpaces + fmt.Sprintf("%.4f", val)
}

func formatWithLeadingSpaces(element int) string {
	leadingSpaces := ""

	if element < 1000 {
		leadingSpaces += " "
		if element < 100 {
			leadingSpaces += " "
			if element < 10 {
				leadingSpaces += " "
			}
		}
	}

	return leadingSpaces + fmt.Sprintf("%d", element)
}

func formatEccentricity(ecc float64) string {
	eccString := fmt.Sprintf("%.7f", ecc)
	return eccString[2:]
}

func formatMeanMotion(meanMotion float64) string {
	leadingSpaces := ""

	if meanMotion < 10 {
		leadingSpaces += " "
	}
	return leadingSpaces + fmt.Sprintf("%.8f", meanMotion)
}

func formatRevNumber(revNumber int) string {
	leadingSpaces := ""

	if revNumber < 10000 {
		leadingSpaces += " "
		if revNumber < 1000 {
			leadingSpaces += " "
			if revNumber < 100 {
				leadingSpaces += " "
				if revNumber < 10 {
					leadingSpaces += " "
				}
			}
		}
	}
	return leadingSpaces + fmt.Sprintf("%d", revNumber)
}
