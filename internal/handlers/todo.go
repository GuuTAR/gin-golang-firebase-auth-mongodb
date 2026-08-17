package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/middleware"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/models"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/services"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/util"
)

// TodoHandler serves the /api/v1/todos endpoints.
type TodoHandler struct {
	svc *services.TodoService
}

// NewTodoHandler returns a TodoHandler backed by svc.
func NewTodoHandler(svc *services.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// Create godoc
//
//	@Summary     Create a todo
//	@Tags        todos
//	@Accept      json
//	@Produce     json
//	@Security    BearerAuth
//	@Param       body  body      models.TodoCreateRequest  true  "Todo"
//	@Success     201   {object}  models.Response{body=models.Todo}
//	@Failure     400   {object}  models.ResponseWithoutBody
//	@Failure     401   {object}  models.ResponseWithoutBody
//	@Failure     500   {object}  models.ResponseWithoutBody
//	@Router      /api/v1/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	uid, _ := middleware.UIDFromContext(c)

	var req models.TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.JSONwithoutBody(c, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.svc.Create(c.Request.Context(), uid, req.Title)
	if err != nil {
		util.RespondError(c, err)
		return
	}

	util.JSON(c, http.StatusCreated, "Todo created successfully", todo)
}

// List godoc
//
//	@Summary     List todos
//	@Tags        todos
//	@Produce     json
//	@Security    BearerAuth
//	@Success     200  {object}  models.Response{body=[]models.Todo}
//	@Failure     401  {object}  models.ResponseWithoutBody
//	@Failure     500  {object}  models.ResponseWithoutBody
//	@Router      /api/v1/todos [get]
func (h *TodoHandler) List(c *gin.Context) {
	uid, _ := middleware.UIDFromContext(c)

	todos, err := h.svc.List(c.Request.Context(), uid)
	if err != nil {
		util.RespondError(c, err)
		return
	}

	util.JSON(c, http.StatusOK, "Todos retrieved successfully", todos)
}

// Toggle godoc
//
//	@Summary     Toggle a todo's completed state
//	@Tags        todos
//	@Produce     json
//	@Security    BearerAuth
//	@Param       id  path  string  true  "Todo ID"
//	@Success     200  {object}  models.Response{body=models.Todo}
//	@Failure     400  {object}  models.ResponseWithoutBody
//	@Failure     401  {object}  models.ResponseWithoutBody
//	@Failure     404  {object}  models.ResponseWithoutBody
//	@Router      /api/v1/todos-toggle/{id} [patch]
func (h *TodoHandler) Toggle(c *gin.Context) {
	uid, _ := middleware.UIDFromContext(c)

	todo, err := h.svc.ToggleCompleted(c.Request.Context(), uid, c.Param("id"))
	if err != nil {
		util.RespondError(c, err)
		return
	}

	util.JSON(c, http.StatusOK, "Todo status toggled successfully", todo)
}

// Reset godoc
//
//	@Summary     Reset all todos to not completed
//	@Tags        todos
//	@Produce     json
//	@Security    BearerAuth
//	@Success     200  {object}  models.ResponseWithoutBody
//	@Failure     401  {object}  models.ResponseWithoutBody
//	@Failure     500  {object}  models.ResponseWithoutBody
//	@Router      /api/v1/todos-reset [patch]
func (h *TodoHandler) Reset(c *gin.Context) {
	uid, _ := middleware.UIDFromContext(c)

	if err := h.svc.ResetCompleted(c.Request.Context(), uid); err != nil {
		util.RespondError(c, err)
		return
	}

	util.JSON(c, http.StatusOK, "Todos reset successfully", nil)
}

// Delete godoc
//
//	@Summary     Delete a todo
//	@Tags        todos
//	@Security    BearerAuth
//	@Param       id  path  string  true  "Todo ID"
//	@Success     200  {object}  models.ResponseWithoutBody
//	@Failure     401  {object}  models.ResponseWithoutBody
//	@Failure     404  {object}  models.ResponseWithoutBody
//	@Router      /api/v1/todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	uid, _ := middleware.UIDFromContext(c)

	if err := h.svc.Delete(c.Request.Context(), uid, c.Param("id")); err != nil {
		util.RespondError(c, err)
		return
	}

	util.JSON(c, http.StatusOK, "Todo deleted successfully", nil)
}
