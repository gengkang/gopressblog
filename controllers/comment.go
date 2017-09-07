package controllers

import (
	"net/http"

	"blog/models"
	"blog/services"
	"strconv"

	"github.com/fpay/gopress"
	"github.com/labstack/echo"
)

// CommentController comment action
type CommentController struct {
	db     *services.DBService
	user   *models.User
	scRule *services.ScoreRule
	group  *echo.Group
}

// NewCommentController returns comment controller instance.
func NewCommentController(group *echo.Group) *CommentController {
	c := new(CommentController)
	c.group = group

	return c
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *CommentController) RegisterRoutes(app *gopress.App) {
	c.db = app.Services.Get(services.DBServerName).(*services.DBService)
	c.user = app.Services.Get(services.UserServiceName).(*services.UserService).User
	c.scRule = app.Services.Get(services.ScoreServiceName).(*services.ScoreService).Rule
	c.group.POST("/comments/create", c.create)
}

// create Action
func (c *CommentController) create(ctx gopress.Context) error {
	postIDStr := ctx.FormValue("post_id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)

	mentionIDStr := ctx.FormValue("mention_user_id")
	mentionID, _ := strconv.ParseUint(mentionIDStr, 10, 64)

	content := ctx.FormValue("content")
	comment := &models.Comment{}
	err := comment.CommentPost(c.db.ORM, uint(postID), c.user.ID, uint(mentionID), content, c.scRule.Comment)
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/blog/posts/"+postIDStr+"?message=保存失败")
	}

	return ctx.Redirect(http.StatusFound, "/blog/posts/"+postIDStr)
}