package util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	// min == max
	require.Equal(t, int64(5), RandomInt(5, 5))
	// max < min
	require.Equal(t, int64(10), RandomInt(10, 3))
	// normal range
	for i := 0; i < 20; i++ {
		n := RandomInt(1, 10)
		require.True(t, n >= 1 && n <= 10, "RandomInt(1,10) = %d", n)
	}
}

func TestRandomString(t *testing.T) {
	require.Len(t, RandomString(0), 0)
	s := RandomString(10)
	require.Len(t, s, 10)
	for _, c := range s {
		require.True(t, (c >= 'a' && c <= 'z'), "character %q not in [a-z]", c)
	}
}

func TestRandomOwner(t *testing.T) {
	require.Len(t, RandomOwner(), 6)
}

func TestRandomMoney(t *testing.T) {
	for i := 0; i < 50; i++ {
		m := RandomMoney()
		require.True(t, m >= 0 && m <= 100, "RandomMoney() = %d", m)
	}
}

func TestRandomCurrency(t *testing.T) {
	allowed := map[string]bool{USD: true, EUR: true, CAD: true}
	for i := 0; i < 20; i++ {
		c := RandomCurrency()
		require.True(t, allowed[c], "RandomCurrency() = %q", c)
	}
}

var emailRegex = regexp.MustCompile(`^[a-z]+@email\.com$`)

func TestRandomEmail(t *testing.T) {
	for i := 0; i < 10; i++ {
		e := RandomEmail()
		require.Regexp(t, emailRegex, e)
	}
}
