package server

import (
	"archive/tar"
	"io"
	"net/http"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"

	api "github.com/proximax-storage/go-xpx-dfms-api"
	drive "github.com/proximax-storage/go-xpx-dfms-drive"

	"github.com/valyala/fasthttp"

	"github.com/gabriel-vasile/mimetype"
)

type Handler interface {
	Serve(ctx *fasthttp.RequestCtx)
}

type gatewayHandler struct {
	api api.Client
}

func NewGatewayHandler(api api.Client) *gatewayHandler {
	return &gatewayHandler{
		api: api,
	}
}

func (g *gatewayHandler) Serve(ctx *fasthttp.RequestCtx) {
	switch {
	case ctx.IsGet():
		g.getFile(ctx)
	default:
		ctx.Error("Method "+string(ctx.Method())+" not allowed", fasthttp.StatusMethodNotAllowed)
	}
}

func (g *gatewayHandler) getFile(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()

	//TODO do we really need check the args length?
	if args.Len() != 2 {
		var msg string
		if args.Len() < 2 {
			msg = "Not enough arguments in the request."
		} else {
			msg = "Many arguments in the request."
		}

		ctx.Error(msg, fasthttp.StatusBadRequest)
		return
	}

	if !args.Has("drive") {
		ctx.Error("Drive parameter not found.", fasthttp.StatusBadRequest)
		return
	}
	if !args.Has("fileCID") {
		ctx.Error("File parameter not found.", fasthttp.StatusBadRequest)
		return
	}

	driveID, err := drive.IDFromBytes(args.Peek("drive"))
	if err != nil {
		ctx.Error("Cannot extract the  drive ID: "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	_, fileCID, err := cid.CidFromBytes(args.Peek("fileCID"))
	if err != nil {
		ctx.Error("Cannot extract the file CID: "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	node, err := g.api.FS().File(ctx, driveID, fileCID)
	if err != nil {
		ctx.Error("Cannot get the file: "+err.Error(), fasthttp.StatusNotFound)
		return
	}

	if f, ok := node.(files.File); ok {
		r := tar.NewReader(f)
		header, err := r.Next()
		if err != nil && err != io.EOF {
			ctx.Error("Error while creating reader: "+err.Error(), fasthttp.StatusBadGateway)
			return
		}
		if header == nil {
			ctx.Error("Returned nil reader for file", fasthttp.StatusBadGateway)
			return
		}

		g.serveFile(ctx, files.NewReaderFile(r))
		return
	} else {
		ctx.Error("Unsupported file type", fasthttp.StatusUnprocessableEntity)
		return
	}
}

func (g *gatewayHandler) serveFile(ctx *fasthttp.RequestCtx, file files.File) {
	m, err := mimetype.DetectReader(file)
	if err != nil {
		ctx.Error("cannot detect content-type: "+err.Error(), http.StatusBadGateway)
		return
	}

	ctx.Response.Header.Set("Content-Type", m.String())
}
