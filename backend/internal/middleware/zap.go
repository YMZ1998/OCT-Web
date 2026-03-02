package middleware

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "os"
)

var Logger *zap.Logger

func InitLogger() {
    logPath := os.Getenv("LOG_PATH")
    cfg := zap.NewProductionConfig()
    cfg.OutputPaths = []string{logPath, "stdout"}
    logger, _ := cfg.Build()
    Logger = logger
}

func ZapLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        Logger.Info("request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
        )
        c.Next()
    }
}
