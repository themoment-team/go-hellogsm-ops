package main

import (
	"fmt"
	"math/rand"
)

func generateFloatPointer(min, max float64) *float64 {
	floatValue := min + (max-min)*rand.Float64()
	return &floatValue
}
func generateIntPointer(min, max int) *int {
	intValue := rand.Intn(max-min+1) + min
	return &intValue
}
func formatNullableFloat(value *float64, scale int) string {
	if value == nil {
		return "NULL"
	}
	format := fmt.Sprintf("%%.%df", scale)
	return fmt.Sprintf(format, *value)
}
func formatNullableString(value *string) string {
	if value == nil {
		return "NULL"
	}
	return fmt.Sprintf("'%s'", *value)
}
