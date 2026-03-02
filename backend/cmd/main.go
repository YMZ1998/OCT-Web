package main

import (
    "oct-backend/pkg/mongo"
    "oct-backend/internal/repository"
    "oct-backend/internal/service"
    "oct-backend/internal/controller"
    "oct-backend/internal/router"
    "oct-backend/internal/middleware"
    "os"
    "log"
    "github.com/joho/godotenv"
)

func main() {
    // 自动加载 .env 文件
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    mongo.InitMongo()
    middleware.InitLogger()
    db := mongo.Client.Database(os.Getenv("MONGO_DB"))
    userRepo := repository.NewUserRepository(db)
    userService := &service.UserService{Repo: userRepo}
    userCtl := &controller.UserController{Service: userService}
    r := router.SetupRouter(userCtl)
    r.Run(":" + os.Getenv("PORT"))
}
