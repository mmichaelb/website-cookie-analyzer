package util

import "testing"

func TestMedianEmptyArray(t *testing.T) {
	expected := float32(-1)
	var input []int
	result := GetMedian(input)
	if result != expected {
		t.Errorf("Expected %f, got: %f", expected, result)
	}
}

func TestMedianSingleElemArray(t *testing.T) {
	expected := float32(4)
	input := []int{
		4,
	}
	result := GetMedian(input)
	if result != expected {
		t.Errorf("Expected %f, got: %f", expected, result)
	}
}

func TestMedianTwoElemArray(t *testing.T) {
	expected := float32(3)
	input := []int{
		2, 4,
	}
	result := GetMedian(input)
	if result != expected {
		t.Errorf("Expected %f, got: %f", expected, result)
	}
}

func TestMedianThreeElemArray(t *testing.T) {
	expected := float32(54)
	input := []int{
		2, 54, 100,
	}
	result := GetMedian(input)
	if result != expected {
		t.Errorf("Expected %f, got: %f", expected, result)
	}
}

func TestMedianMiddlePresent(t *testing.T) {
	expected := float32(5)
	input := []int{
		0, 2, 3, 3,
		5,
		6, 7, 7, 8,
	}
	result := GetMedian(input)
	if result != expected {
		t.Errorf("Expected %f, got: %f", expected, result)
	}
}

func TestMedianMiddleNotPresent(t *testing.T) {
	expected := float32(4.5)
	input := []int{
		0, 2, 3, 3,
		6, 7, 7, 8,
	}
	result := GetMedian(input)
	if result != expected {
		t.Errorf("Expected %f, got: %f", expected, result)
	}
}
