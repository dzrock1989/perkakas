package perkakas

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrings(t *testing.T) {
	result := IsEmpty("")
	require.Equal(t, true, result, "Should return true")

	result = IsEmpty("Omama")
	require.Equal(t, false, result, "Should return true")

	result = IsEqual("a", "b")
	require.Equal(t, false, result, "Should return false")

	result = IsEqual("a", "a")
	require.Equal(t, true, result, "Should return true")

	result = IsEqual("a", "A")
	require.Equal(t, false, result, "Should return true")
}
