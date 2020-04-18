package hook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAsyncSettings(t *testing.T) {
	var data = []struct {
		settings map[string]string
		enable   bool
		size     int
		block    bool
	}{
		{map[string]string{}, AsyncEnable, DefaultAsyncBufferSize, DefaultAsyncBlock},
		{map[string]string{
			"async_enable":      "0",
			"async_buffer_size": "1000",
			"async_block":       "1",
		}, false, 1000, true},
		{map[string]string{
			"async_enable":      "1",
			"async_buffer_size": "2000",
			"async_block":       "0",
		}, true, 2000, false},
	}

	for _, d := range data {
		enable, size, block := getAsyncSettings(d.settings)
		assert.Equal(t, d.enable, enable)
		assert.Equal(t, d.size, size)
		assert.Equal(t, d.block, block)
	}
}
