package fast

import "github.com/valyala/fasthttp"

func GetHeader(cx *fasthttp.RequestCtx, s string) string {
	return string(cx.Request.Header.Peek(s))
}
