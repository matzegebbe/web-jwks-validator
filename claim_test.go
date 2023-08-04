package main

import (
	"testing"
)

func TestCheckIfClaimContainsAllClaimContainsCheck(t *testing.T) {
	// Define some test cases
	testCases := []struct {
		claims             map[string]interface{}
		claimContainsCheck []string
		expectedResult     bool
		description        string
	}{
		{
			claims: map[string]interface{}{
				"name": "John",
				"tags": []interface{}{"student", "admin"},
			},
			claimContainsCheck: []string{"name=John", "tags=admin"},
			expectedResult:     true,
			description:        "Case 1: Expected to return true",
		},
		{
			claims: map[string]interface{}{
				"name": "John",
				"tags": []string{"student", "developer"},
			},
			claimContainsCheck: []string{"name=John", "tags=admin"},
			expectedResult:     false,
			description:        "Case 2: Expected to return false",
		},
		{
			claims: map[string]interface{}{
				"name": "Ralle",
				"tags": []string{"student", "developer"},
			},
			claimContainsCheck: []string{"name=John", "tags=admin"},
			expectedResult:     false,
			description:        "Case 2: Expected to return false",
		},
	}

	for _, tc := range testCases {
		result := checkIfClaimContainsAllClaimContainsCheck(tc.claims, tc.claimContainsCheck)
		if result != tc.expectedResult {
			t.Errorf("%s: checkIfClaimContainsAllClaimContainsCheck() = %v;"+
				"want %v", tc.description, result, tc.expectedResult)
		}
	}
}
