package fast

import (
	"github.com/StevenZack/tools/fast/views"
	"github.com/valyala/fasthttp"
)

func GetHeader(cx *fasthttp.RequestCtx, s string) string {
	return string(cx.Request.Header.Peek(s))
}
func GetURI(cx *fasthttp.RequestCtx) string {
	return string(cx.URI().Path())
}
func NotFound(cx *fasthttp.RequestCtx) {
	cx.Response.SetStatusCode(404)
	SetHeaderHTML(cx)
	cx.WriteString(views.Str_404)
}
func SetHeaderHTML(cx *fasthttp.RequestCtx) {
	cx.Response.Header.Set("Content-Type", "text/html")
}
