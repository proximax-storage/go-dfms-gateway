package go_dfms_gateway

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	files "github.com/ipfs/go-ipfs-files"

	api "github.com/proximax-storage/go-xpx-dfms-api"
	drive "github.com/proximax-storage/go-xpx-dfms-drive"

	"github.com/valyala/fasthttp"

	"github.com/gabriel-vasile/mimetype"
)

const cidPattern = "[a-z1-9]{65}"

type DriveList struct {
	Drives []string `json:"drives"`
}

type DirList struct {
	Nodes []string `json:"list"`
}

type Handler interface {
	Handle(ctx *fasthttp.RequestCtx)
}

func newMiddleware(handler Handler, cors *cors) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Debugf("New %s Request %s, Client %s", ctx.Method(), ctx.URI(), ctx.RemoteAddr())

		if err := cors.check(ctx); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusNotAcceptable)
			return
		}

		if ctx.IsOptions() {
			return
		}

		handler.Handle(ctx)
		switch {
		case ctx.Response.StatusCode() == fasthttp.StatusNotFound:
			notFound(ctx)
		case ctx.Response.StatusCode() >= 500:
			serverError(ctx)
		}
	}
}

type gatewayHandler struct {
	api api.Client
}

func newGatewayHandler(api api.Client) *gatewayHandler {
	return &gatewayHandler{
		api: api,
	}
}

func (gh *gatewayHandler) Handle(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.Error("Method "+string(ctx.Method())+" not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	path := ctx.Path()
	cid := ""
	filePath := "/"

	switch {
	case string(path) == "/" || string(path) == "":
		// index page, list a node drives
		gh.getDrives(ctx)
		return
	case match("^/"+cidPattern+"/?$", path):
		// if only drive cid than list root
		cid = strings.Trim(string(path), "/")
	case match("^/"+cidPattern+"/.*", path):
		// cid + file path
		parsedPath := strings.SplitN(strings.Trim(string(path), "/"), "/", 2)
		cid = parsedPath[0]
		filePath += parsedPath[1]
	default:
		// path is not supported
		ctx.Error("Bad route", fasthttp.StatusNotFound)
		return
	}

	driveID, err := drive.IDFromString(cid)
	if err != nil {
		ctx.Error("Cannot extract drive ID: "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	gh.getFile(ctx, driveID, filePath)
}

func (gh *gatewayHandler) getDrives(ctx *fasthttp.RequestCtx) {
	ls, err := gh.api.Contract().List(ctx)
	if err != nil {
		ctx.Error("Cannot get the drives list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(ls) == 0 {
		ctx.Response.SetBody([]byte("No drives"))
		return
	}

	driveList := DriveList{}
	for _, v := range ls {
		driveList.Drives = append(driveList.Drives, v.String())
	}

	content, err := json.Marshal(driveList)
	if err != nil {
		ctx.Error("Cannot create JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetBody(content)
}

func (gh *gatewayHandler) getFile(ctx *fasthttp.RequestCtx, driveID drive.ID, filePath string) {
	node, err := gh.api.FS().Get(ctx, driveID, filePath)
	if err != nil {
		ctx.Error("Cannot get the file or directory: "+err.Error(), fasthttp.StatusNotFound)
		return
	}

	gh.serveNode(ctx, node)
}

func (gh *gatewayHandler) serveNode(ctx *fasthttp.RequestCtx, node files.Node) {
	if f, ok := node.(files.File); ok {
		gh.serveFile(ctx, f)
		return
	}

	if dir, ok := node.(files.Directory); ok {
		gh.serveDirectory(ctx, dir)
		return
	}

	ctx.Error("Unsupported file type", fasthttp.StatusInternalServerError)
}

func (gh *gatewayHandler) serveFile(ctx *fasthttp.RequestCtx, file files.File) {
	defer file.Close()

	size, err := file.Size()
	if err != nil {
		ctx.Error("Cannot get file size: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// is this a good idea?
	r := io.TeeReader(file, ctx.Response.BodyWriter())

	m, err := mimetype.DetectReader(r)
	if err != nil {
		ctx.Error("Cannot detect content-type: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType(m.String())
	ctx.Response.Header.SetContentLength(int(size))

	_, err = io.Copy(ctx.Response.BodyWriter(), file)
	if err != nil {
		ctx.Error("Cannot write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (gh *gatewayHandler) serveDirectory(ctx *fasthttp.RequestCtx, dir files.Directory) {
	defer dir.Close()

	dirList := DirList{}
	di := dir.Entries()
	for di.Next() {
		dirList.Nodes = append(dirList.Nodes, di.Name())
	}

	if len(dirList.Nodes) == 0 {
		ctx.Response.SetBody([]byte("Directory is empty"))
		return
	}

	content, err := json.Marshal(dirList)
	if err != nil {
		ctx.Error("Cannot create JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetBody(content)
}

//TODO handle 404 error
func notFound(ctx *fasthttp.RequestCtx) {
	log.Debugf("Client: %s, %s Request %s, URL: %s", ctx.RemoteAddr(), ctx.Method(), ctx.URI(), ctx.Response.Body())
}

//TODO handle server errors
func serverError(ctx *fasthttp.RequestCtx) {
	log.Warnf("Client: %s, %s Request %s, URL: %s", ctx.RemoteAddr(), ctx.Method(), ctx.URI(), ctx.Response.Body())
}
