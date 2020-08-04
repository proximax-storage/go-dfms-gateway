package server

import (
	"bytes"
	"encoding/json"
	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/valyala/fasthttp"

	files "github.com/ipfs/go-ipfs-files"
)

//TODO come up with better tests. Maybe also add mocks
func TestGatewayHandler(t *testing.T) {
	api := apihttp.NewClientAPI("http://localhost:6366")
	gateway := NewGateway(api)
	go gateway.Start()

	addr := gateway.address
	if match("^[:][0-9]{4}$", []byte(addr)) {
		addr = "http://localhost" + addr
	}
	time.Sleep(5 * time.Second)

	client := fasthttp.Client{}

	t.Run("Index", func(t *testing.T) {
		var dst []byte
		statusCode, b, err := client.Get(dst, addr)
		assert.Nil(t, err, err)

		if statusCode == fasthttp.StatusInternalServerError {
			t.Skip()
		}
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
		if fasthttp.StatusOK != statusCode {
			t.Skip()
		}

		driveList := DirList{}
		err = json.Unmarshal(b, &driveList)
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

	gateway.Stop()
}

type (
	file struct {
		path    []byte
		content []byte
	}
	folder struct {
		path    []byte
		content []string
	}
	unknown []byte

	mockHandler struct {
		gatewayHandler
		file    file
		folder  folder
		unknown unknown
	}
)

func initMock() mockHandler {
	return mockHandler{
		gatewayHandler: gatewayHandler{},
		file: file{
			path:    []byte("/file"),
			content: []byte("test_content"),
		},
		folder: folder{
			path: []byte("/folder"),
			content: []string{
				"first",
				"second",
				"third",
			},
		},
		unknown: []byte("/unknown"),
	}
}

func (m *mockHandler) Handle(ctx *fasthttp.RequestCtx) {
	switch {
	case bytes.Equal(ctx.Path(), m.file.path):
		m.mockFile(ctx)
	case bytes.Equal(ctx.Path(), m.folder.path):
		m.mockFolder(ctx)
	}
}

func (m *mockHandler) mockFile(ctx *fasthttp.RequestCtx) {
	file := files.NewBytesFile(m.file.content)
	m.serveNode(ctx, file)
}

func (m *mockHandler) mockFolder(ctx *fasthttp.RequestCtx) {
	m.serveNode(ctx, newDir(m.folder.content))
}

func TestGatewayHandler_Mock(t *testing.T) {
	client := fasthttp.Client{}
	mh := initMock()

	s := gateway{
		server: fasthttp.Server{
			Handler: mh.Handle,
		},
		address: "localhost:5000",
	}
	go s.Start()
	defer s.Stop()

	serverAddrs := "http://" + s.address

	time.Sleep(5 * time.Second)
	t.Run("File", func(t *testing.T) {
		var dst []byte
		statusCode, b, err := client.Get(dst, serverAddrs+"/"+string(mh.file.path))
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)
		assert.Equal(t, mh.file.content, b)
	})

	t.Run("Folder", func(t *testing.T) {
		var dst []byte
		statusCode, b, err := client.Get(dst, serverAddrs+"/"+string(mh.folder.path))
		assert.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusOK, statusCode)

		var list DirList
		err = json.Unmarshal(b, &list)
		assert.Nil(t, err, err)
		assert.Equal(t, mh.folder.content, list.Nodes)
	})
}

func newDir(content []string) files.Directory {
	dir := make(map[string]files.Node, len(content))
	for _, v := range content {
		dir[v] = files.NewBytesFile([]byte(""))
	}

	return files.NewMapDirectory(dir)
}
