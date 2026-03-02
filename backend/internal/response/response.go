package response

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type Response struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, code int, msg string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code: code,
        Msg:  msg,
        Data: data,
    })
}
