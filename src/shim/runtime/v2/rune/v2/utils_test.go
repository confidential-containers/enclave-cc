package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainerIDSuccess(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		image string
		cid   string
	}{
		{
			image: "docker.io/library/busybox",
			cid:   "busybox",
		},
		{
			image: "docker.io/library/busybox:latest",
			cid:   "busybox_latest",
		},
	}

	for _, tt := range tests {
		cid, err := getContainerID(tt.image)
		if cid != tt.cid {
			t.Errorf("Image: \"%v\" has cid \"%v\", wants \"%v\"", tt.image, cid, tt.cid)
		}
		assert.NoError(err)
	}

	_, err := getContainerID("")
	assert.Error(err)
}

func TestGetContainerIDFail(t *testing.T) {
	assert := assert.New(t)

	_, err := getContainerID("")

	assert.Error(err)
}
