package main

import "testing"

func TestNegativeSearchType(t *testing.T) {
	result, err := ParseAndValidateInput("-1")
	if err == nil {
		t.Error("Expected an error, got", result)
	}
}

func TestTooLargeSearchType(t *testing.T) {
	result, err := ParseAndValidateInput("4")
	if err == nil {
		t.Error("Expected an error, got", result)
	}
}

func TestStringSearchType(t *testing.T) {
	result, err := ParseAndValidateInput("A")
	if err == nil {
		t.Error("Expected an error, got", result)
	}
}

func TestValidSearchTypes(t *testing.T) {

	tables := []struct {
		searchType string
		result int
	}{
		{"1", 1}, {"2",2}, {"3",3},
	}

	for _, table := range tables {
		result, _ := ParseAndValidateInput(table.searchType)
		if result != table.result {
			t.Errorf("Expected %s got %d", table.searchType, result)
		}
	}
}