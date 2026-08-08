package payment

import "testing"

func TestValidateTransition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		from    Status
		to      Status
		wantErr bool
	}{
		{name: "created to processing", from: StatusCreated, to: StatusProcessing},
		{name: "processing to success", from: StatusProcessing, to: StatusSuccess},
		{name: "success to refunded", from: StatusSuccess, to: StatusRefunded},
		{name: "success to processing invalid", from: StatusSuccess, to: StatusProcessing, wantErr: true},
		{name: "failed to success invalid", from: StatusFailed, to: StatusSuccess, wantErr: true},
		{name: "refunded to success invalid", from: StatusRefunded, to: StatusSuccess, wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateTransition(tc.from, tc.to)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
