// Package signature
// Мидлвэр для подписи запроса
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

// responseWriterWithHash
// Реализация http.ResponseWriter с поддержкой подсчета хеша
type responseWriterWithHash struct {
	Writer     http.ResponseWriter
	hash       hash.Hash
	buf        *bytes.Buffer
	statusCode int
	n          int
}

func (w *responseWriterWithHash) Header() http.Header {
	return w.Writer.Header()
}

func (w *responseWriterWithHash) WriteHeader(code int) {
	w.statusCode = code
	w.Writer.WriteHeader(code)
}

func (w *responseWriterWithHash) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	n, err := w.buf.Write(b)
	if err != nil {
		return n, err
	}
	w.n += n
	return w.hash.Write(b)
}

func (w *responseWriterWithHash) Close() error {
	if w.n != 0 {
		w.Writer.Header().Set("HashSHA256", hex.EncodeToString(w.hash.Sum(nil)))
		_, err := w.Writer.Write(w.buf.Bytes())
		return err
	}
	return nil
}

// SignatureMiddleware
// Реализует мидлвэр, который проверяет подпись запроса, а так же добавляет подпись в заголовок ответа на основе тела ответа
func SignatureMiddleware(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			signature := c.Request().Header.Get("HashSHA256")

			if signature != "" {
				body, err := io.ReadAll(c.Request().Body)
				if err != nil {
					return c.NoContent(http.StatusInternalServerError)
				}

				c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

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
			defer writer.Close()

			c.Response().Writer = writer

			err := next(c)

			return err

		}
	}
}
