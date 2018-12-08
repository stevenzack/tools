package encodingToolkit

import "testing"

func Test_UrlEncode(t *testing.T) {
	t.Log(UrlDecode(`%E5%8E%8B%E7%BC%A9%E5%90%8E%E7%9A%84root%E7%9B%AE%E5%BD%95.zip`))
}
