package signature

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type responseWriterWithHash struct {
	Writer     http.ResponseWriter
	hash       hash.Hash
	buf        *bytes.Buffer
	statusCode int
}

func (w *responseWriterWithHash) Header() http.Header {
	return w.Writer.Header()
}

func (w *responseWriterWithHash) WriteHeader(code int) {
	w.statusCode = code
}

func (w *responseWriterWithHash) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	l, err := w.buf.Write(b)
	if err != nil {
		return l, err
	}
	return w.hash.Write(b)
}

func (w *responseWriterWithHash) Close() error {
	w.Writer.Header().Set("HashSHA256", hex.EncodeToString(w.hash.Sum(nil)))
	_, err := w.Writer.Write(w.buf.Bytes())
	return err
}

func SignatureMiddleware(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			if len(body) != 0 {
				c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

				signature := c.Request().Header.Get("HashSHA256")

				h := sha256.Sum256([]byte(string(body) + key))
				calculatedSignature := fmt.Sprintf("%x", h)

				if signature != calculatedSignature {
					return c.NoContent(http.StatusBadRequest)
				}
			}

			writer := &responseWriterWithHash{
				Writer: c.Response().Writer,
				buf:    bytes.NewBuffer([]byte{}),
				hash:   sha256.New(),
			}
			c.Response().Writer = writer

			err = next(c)

			return err

		}
	}
}
