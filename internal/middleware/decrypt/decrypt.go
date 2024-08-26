package decrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func RSADecryptMiddleware(key []byte) (echo.MiddlewareFunc, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("invalid block type")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			encryptedBody, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.NoContent(http.StatusBadRequest)
			}

			decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBody)
			fmt.Println(string(decryptedBody))
			if err != nil {
				return c.NoContent(http.StatusBadRequest)
			}

			c.Request().Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

			return next(c)
		}
	}, nil
}
