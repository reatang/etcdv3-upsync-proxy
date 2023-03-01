package xbase64

import (
	"encoding/base64"
)

func Decode(src []byte) (string, error) {
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(src)))

	n, err := base64.StdEncoding.Decode(dbuf, src)
	if err != nil {
		return "", err
	}

	return string(dbuf[:n]), nil
}
