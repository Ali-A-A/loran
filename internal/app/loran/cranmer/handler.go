package cranmer

import (
	"encoding/json"
	"net/http"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Handler represents cranmer handler.
// It should have *cmq.Conn to publish new messages on nats server.
type Handler struct {
	nc      *cmq.Conn
	subject string
}

// NewHandler creates new cranmer Handler.
func NewHandler(nc *cmq.Conn, cfg config.NATS) *Handler {
	return &Handler{
		nc:      nc,
		subject: cfg.JetStream.Consumer.Subject,
	}
}

// Add responsible to publish new requests on to the nats server.
// It may failed if it cannot parse request.
// In this case, it returns http.StatusBadRequest status.
// In the interval error cases, like failure in publishing, it returns http.StatusInternalServerError.
// Otherwise, it returns http.StatusOK.
func (h *Handler) Add(c echo.Context) error {
	req := &AddRequest{}

	if err := c.Bind(req); err != nil {
		logrus.Warnf("failed to bind request: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": "request's body is invalid"})
	}

	if req.UserID == 0 || req.EntityID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "user_id or entity_id is invalid"})
	}

	b, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("failed to marshal request: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	if _, err = h.nc.JS.Publish(h.subject, b); err != nil {
		logrus.Errorf("failed to publish message: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
