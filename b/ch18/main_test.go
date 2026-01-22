package main

import (
	"math"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{"Freezing point of water", 0.0, 32.0},
		{"Boiling point of water", 100.0, 212.0},
		{"Body temperature", 37.0, 98.6},
		{"Room temperature", 20.0, 68.0},
		{"Negative temperature", -40.0, -40.0}, // -40°C is -40°F
		{"Extreme cold", -273.15, -459.67},     // Absolute zero
		{"Decimal value", 25.5, 77.9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelsiusToFahrenheit(tt.celsius)
			// Allow for a small margin of error due to floating point arithmetic
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("CelsiusToFahrenheit(%v) = %v, expected %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		name       string
		fahrenheit float64
		expected   float64
	}{
		{"Freezing point of water", 32.0, 0.0},
		{"Boiling point of water", 212.0, 100.0},
		{"Body temperature", 98.6, 37.0},
		{"Room temperature", 68.0, 20.0},
		{"Negative temperature", -40.0, -40.0}, // -40°F is -40°C
		{"Extreme cold", -459.67, -273.15},     // Absolute zero
		{"Decimal value", 77.9, 25.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FahrenheitToCelsius(tt.fahrenheit)
			// Allow for a small margin of error due to floating point arithmetic
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("FahrenheitToCelsius(%v) = %v, expected %v", tt.fahrenheit, result, tt.expected)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		decimals int
		expected float64
	}{
		{"Round to 2 decimals", 12.345, 2, 12.35},
		{"Round to 0 decimals", 12.345, 0, 12.0},
		{"Round to 1 decimal", 12.345, 1, 12.3},
		{"Round down", 12.344, 2, 12.34},
		{"Round negative number", -12.345, 2, -12.35},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.value, tt.decimals)
			if result != tt.expected {
				t.Errorf("Round(%v, %v) = %v, expected %v", tt.value, tt.decimals, result, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that converting from Celsius to Fahrenheit and back gives the original value
	originalCelsius := 25.0
	roundTrip := FahrenheitToCelsius(CelsiusToFahrenheit(originalCelsius))
	if math.Abs(roundTrip-originalCelsius) > 0.01 {
		t.Errorf("Round trip conversion failed: original=%v, after conversion=%v", originalCelsius, roundTrip)
	}

	// Test that converting from Fahrenheit to Celsius and back gives the original value
	originalFahrenheit := 77.0
	roundTrip = CelsiusToFahrenheit(FahrenheitToCelsius(originalFahrenheit))
	if math.Abs(roundTrip-originalFahrenheit) > 0.01 {
		t.Errorf("Round trip conversion failed: original=%v, after conversion=%v", originalFahrenheit, roundTrip)
	}
}
