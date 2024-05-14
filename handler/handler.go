package handler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"simple-users-tasks-service/handler/model"
	"simple-users-tasks-service/usecase"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"gopkg.in/oauth2.v3/server"
)

type Handler struct {
	srv     *server.Server
	usecase usecase.Usecase
}

func NewHandler(srv *server.Server, u usecase.Usecase) Handler {

	return Handler{
		srv:     srv,
		usecase: u,
	}
}

// User handlers
func (h *Handler) HandleRegisterUser(c *gin.Context) {
	var reqBody model.RegisterUserRequest
	err := c.Bind(&reqBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to parse request body",
		})
		return
	}

	statusCode, err := h.usecase.ProcessRegisterUser(reqBody)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register successful",
	})
}

func (h *Handler) HandleGetAllUsers(c *gin.Context) {
	data, statusCode, err := h.usecase.ProcessGetAllUsers()
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (h *Handler) HandleGetUserDetails(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ID",
		})
		return
	}
	data, statusCode, err := h.usecase.ProcessGetUserDetails(idInt)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (h *Handler) HandleUpdateUserData(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ID",
		})
		return
	}

	var reqBody model.UpdateUserRequest
	err = c.Bind(&reqBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to parse request body",
		})
		return
	}

	statusCode, err := h.usecase.ProcessUpdateUserData(idInt, reqBody)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) HandleDeleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ID",
		})
		return
	}

	statusCode, err := h.usecase.ProcessDeleteUserData(idInt)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Tash handlers
func (h *Handler) HandleCreateTask(c *gin.Context) {
	var reqBody model.CreateTaskRequest
	err := c.Bind(&reqBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to parse request body",
		})
		return
	}

	statusCode, err := h.usecase.ProcessCreateTask(reqBody)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "new task created",
	})
}

func (h *Handler) HandleGetAllTasks(c *gin.Context) {
	data, statusCode, err := h.usecase.ProcessGetAllTasks()
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

func (h *Handler) HandleGetTaskDetails(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ID",
		})
		return
	}

	data, statusCode, err := h.usecase.ProcessGetTaskDetails(idInt)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

func (h *Handler) HandleUpdateTask(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid task ID",
		})
		return
	}

	var reqBody model.UpdateTaskRequest
	err = c.Bind(&reqBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to parse request body",
		})
		return
	}

	statusCode, err := h.usecase.ProcessUpdateTask(idInt, reqBody)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) HandleDeleteTask(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ID",
		})
		return
	}

	statusCode, err := h.usecase.ProcessDeleteTask(idInt)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusNoContent, nil)
}

// OAuth handlers

func (h *Handler) HandleLoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *Handler) HandleLoginPost(c *gin.Context) {
	store := ginsession.FromContext(c)

	store.Set("LoggedInUserID", "000000")
	err := store.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.Header("Location", "/authenticate")
	c.Writer.WriteHeader(http.StatusFound)
}

func (h *Handler) HandleAuthenticateGet(c *gin.Context) {
	store := ginsession.FromContext(c)

	if _, ok := store.Get("LoggedInUserID"); !ok {
		c.Header("Location", "/login")
		c.Writer.WriteHeader(http.StatusFound)
		return
	}

	c.HTML(http.StatusOK, "authenticate.html", nil)
}

func (h *Handler) HandleAuthenticatePost(c *gin.Context) {
	store := ginsession.FromContext(c)

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	u := new(url.URL)
	u.Path = "/authorize"
	u.RawQuery = form.Encode()
	c.Header("Location", u.String())
	c.Writer.WriteHeader(http.StatusFound)
	store.Delete("Form")

	if v, ok := store.Get("LoggedInUserID"); ok {
		store.Set("UserID", v)
	}
	store.Save()
}

func (h *Handler) HandleAuthorize(c *gin.Context) {
	err := h.srv.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "authorize error",
		})
	}
}

func (h *Handler) HandleGetToken(c *gin.Context) {
	err := h.srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "get token error",
		})
	}
}
