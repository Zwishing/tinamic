package storage

import (
	"context"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNewMinioClient(t *testing.T) {
	c := NewMinioClient()
	ctx := context.Background()
	buckets, err := c.cli.ListBuckets(ctx)
	if err != nil {
		return
	}
	assert.Equal(t, buckets[0].Name, "raster", "通过！")
}
