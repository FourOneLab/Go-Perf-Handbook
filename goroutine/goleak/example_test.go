package goleak

import (
	"testing"

	"go.uber.org/goleak"
)

func Test_leak(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "leak goroutine",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := leak(); (err != nil) != tt.wantErr {
				t.Errorf("leak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
