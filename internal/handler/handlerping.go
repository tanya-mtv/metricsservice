package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerPing struct {
	db pingDB
}

func NewHandlerPing(db *sqlx.DB) *HandlerPing {
	return &HandlerPing{
		db: db,
	}
}

func (h *HandlerPing) Ping(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var ptr *sqlx.DB

	if ptr == h.db {
		newErrorResponse(c, http.StatusInternalServerError, "Can't connect to database")
		return
	}

	err := h.db.Ping()

	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, "Can't connect to database")
		return
	}
	c.JSON(http.StatusOK, "Success connection to database")
}
