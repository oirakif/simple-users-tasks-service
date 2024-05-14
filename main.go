package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-users-tasks-service/handler"
	"simple-users-tasks-service/usecase"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/go-session/session"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	usecase := usecase.NewUsecase(db)
	handler := handler.NewHandler(srv, usecase)

	r := gin.Default()
	r.Use(ginsession.New())

	iniitateUsersEndpoints(r, handler)
	iniitateTasksEndpoints(r, handler)

	initiateOAuthEndpoints(r, handler)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func iniitateUsersEndpoints(r *gin.Engine, h handler.Handler) {
	r.POST("/users", h.HandleRegisterUser)

	r.GET("/users", h.HandleGetAllUsers)

	r.GET("/users/:id", h.HandleGetUserDetails)
	r.PUT("/users/:id", h.HandleUpdateUserData)
	r.DELETE("/users/:id", h.HandleDeleteUser)

}

func iniitateTasksEndpoints(r *gin.Engine, h handler.Handler) {
	r.POST("/tasks", h.HandleCreateTask)

	r.GET("/tasks", h.HandleGetAllTasks)

	r.GET("/tasks/:id", h.HandleGetTaskDetails)
	r.PUT("/tasks/:id", h.HandleUpdateTask)
	r.DELETE("/tasks/:id", h.HandleDeleteTask)

}

// TODO: make it work
func initiateOAuthEndpoints(r *gin.Engine, h handler.Handler) {
	r.LoadHTMLGlob("forms/*")
	r.GET("/login", h.HandleLoginGet)
	r.POST("/login", h.HandleLoginPost)
	r.GET("/authenticate", h.HandleAuthenticateGet)
	r.POST("/authenticate", h.HandleAuthenticatePost)
	r.GET("/authorize", h.HandleAuthorize)
	r.GET("/token", h.HandleGetToken)
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("UserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("UserID")
	store.Save()
	return
}
