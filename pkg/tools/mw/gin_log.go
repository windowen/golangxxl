package mw

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"queueJob/pkg/constant"
	"queueJob/pkg/tools/apiresp"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/zlogger"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		language := strings.TrimSpace(c.GetHeader(constant.Language))
		platform := strings.TrimSpace(c.GetHeader(constant.OpUserPlatform))
		operationID := strings.TrimSpace(c.GetHeader(constant.OperationId))
		countryCode := strings.TrimSpace(c.GetHeader(constant.CountryCode))

		if language == "" || platform == "" || operationID == "" || countryCode == "" {
			c.Abort()
			apiresp.GinError(c, errs.ErrArgs.WithDetail("para_err"))
			return
		}

		setRequireParamsWithOpts(c,
			WithPlatform(platform),
			WithLanguage(language),
			WithOperationID(operationID),
			WithCountryCode(countryCode),
		)

		req, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		start := time.Now()
		c.Request.Body = io.NopCloser(bytes.NewReader(req))
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBuffer(nil),
		}
		c.Writer = writer
		c.Next()
		resp := writer.buf.Bytes()
		zlogger.Debugw("gin response", zap.Int64("time", int64(time.Since(start))), zap.Int("status", c.Writer.Status()), zap.String("resp", string(resp)))
	}
}
