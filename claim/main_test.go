package claim

import (
	"reflect"
	"testing"
)

func TestPreparePayload(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		claims    []string
		operation string
		want      claimPayload
	}{
		{
			name:      "Single claim",
			claims:    []string{"claim1:.*"},
			operation: "insert",
			want: claimPayload{
				Operation: "insert",
				Claims:    map[string]string{"claim1": ".*"},
			},
		},
		{
			name:      "Multiple claims",
			claims:    []string{"claim1:.*", "claim2:.*"},
			operation: "delete",
			want: claimPayload{
				Operation: "delete",
				Claims:    map[string]string{"claim1": ".*", "claim2": ".*"},
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := preparePayload(tt.claims, tt.operation)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preparePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
