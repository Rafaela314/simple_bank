package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedCurrency(t *testing.T) {
	tests := []struct {
		currency string
		want     bool
	}{
		{USD, true},
		{EUR, true},
		{CAD, false},
		{"", false},
		{"GBP", false},
	}
	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			got := IsSupportedCurrency(tt.currency)
			require.Equal(t, tt.want, got)
		})
	}
}
