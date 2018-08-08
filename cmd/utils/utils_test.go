package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatMbg(t *testing.T) {
	assert.Equal(t, "1", FormatMbg(1000000000))
	assert.Equal(t, "1.1", FormatMbg(1100000000))
	assert.Equal(t, "1.123456789", FormatMbg(1123456789))
	assert.Equal(t, "1000000000.123456789", FormatMbg(1000000000123456789))
	assert.Equal(t, "1000000000.000001", FormatMbg(1000000000000001000))
	assert.Equal(t, "1000000000.000000001", FormatMbg(1000000000000000001))
}

func TestParseMbg(t *testing.T) {
	assert.Equal(t, uint64(1000000000), ParseMbg("1"))
	assert.Equal(t, uint64(1000000000000000000), ParseMbg("1000000000"))
	assert.Equal(t, uint64(1000000000123456789), ParseMbg("1000000000.123456789"))
	assert.Equal(t, uint64(1000000000000000100), ParseMbg("1000000000.0000001"))
	assert.Equal(t, uint64(1000000000000000001), ParseMbg("1000000000.000000001"))
	assert.Equal(t, uint64(1000000000000000001), ParseMbg("1000000000.000000001123"))
}

func TestFormatMbc(t *testing.T) {
	assert.Equal(t, "0", FormatMbc(0))
	assert.Equal(t, "1", FormatMbc(1))
	assert.Equal(t, "100", FormatMbc(100))
	assert.Equal(t, "1000000000", FormatMbc(1000000000))
}

func TestParseMbc(t *testing.T) {
	assert.Equal(t, uint64(0), ParseMbc("0"))
	assert.Equal(t, uint64(1), ParseMbc("1"))
	assert.Equal(t, uint64(1000), ParseMbc("1000"))
	assert.Equal(t, uint64(1000000000), ParseMbc("1000000000"))
	assert.Equal(t, uint64(1000000), ParseMbc("1000000.123"))
}
