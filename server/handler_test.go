package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/valyala/fasthttp"
)

func TestGatewayHandler_Handle(t *testing.T) {
	s := NewGateway()
	go s.Start()

	addr := s.address
	if match("^[:][0-9]{4}$", []byte(addr)) {
		addr = "http://localhost" + addr
	}
	time.Sleep(5 * time.Second)

	client := fasthttp.Client{}

	t.Run("Index", func(t *testing.T) {
		var dst []byte
		statusCode, b, err := client.Get(dst, addr)
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)

		dl := DriveList{}
		err = json.Unmarshal(b, &dl)
		assert.Nil(t, err, err)
		if len(dl.Drives) == 0 {
			assert.Equal(t, "No drives", string(b))
		}

		statusCode, b, err = client.Get(dst, addr+"/")
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)
	})

	t.Run("Drive", func(t *testing.T) {
		var dst []byte
		statusCode, b, err := client.Get(dst, addr)
		assert.Nil(t, err, err)

		dl := DriveList{}
		json.Unmarshal(b, &dl)

		if len(dl.Drives) == 0 {
			t.Skip()
		}

		statusCode, b, err = client.Get(dst, addr+"/"+dl.Drives[0])
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)

		driveList := DirList{}
		err = json.Unmarshal(b, &dl)
		assert.Nil(t, err, err)
		if len(driveList.Nodes) == 0 {
			assert.Equal(t, "Directory is empty", string(b))
		}

		statusCode, b, err = client.Get(dst, addr+"/"+dl.Drives[0]+"/")
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		var dst []byte
		statusCode, b, err := client.Get(dst, addr+"/bad")
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusNotFound, statusCode)
		assert.Contains(t, string(b), "Bad route")
	})

	t.Run("InternalError", func(t *testing.T) {
		var dst []byte
		statusCode, _, err := client.Get(dst, addr+"/aaaaaaaaaaaaaafkyhloktcagx6sctriu7ehx5supn2ryejpdhjqc4flvgzeui4d4")
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusInternalServerError, statusCode)

		statusCode, _, err = client.Get(dst, addr+"/aaaaaaaaaaaaaafkyhloktcagx6sctriu7ehx5supn2ryejpdhjqc4flvgzeui4d4/")
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusInternalServerError, statusCode)
	})

	s.Stop()
}
