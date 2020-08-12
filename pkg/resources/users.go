package resources

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/fountain"
	"github.com/gin-gonic/gin"
)

// UsersResource contains handlers for the users endpoint
type UsersResource struct {
	userStore data.UserStore
}

// NewUsersResource creates a new users resource instance
func NewUsersResource(userStore data.UserStore) *UsersResource {
	resource := new(UsersResource)
	resource.userStore = userStore
	return resource
}

// GetUser retrieves a user by ID
func (r *UsersResource) GetUser(c *gin.Context) {
	id := c.Param("id")

	if user, err := r.userStore.User(id); err == nil {
		c.JSON(http.StatusOK, user)
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}
}

type newUser struct {
	ID       string `form:"id" binding:"required,alphanum"`
	FullName string `form:"fullName" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
}

// PostUser pushes a new user into the store
func (r *UsersResource) PostUser(c *gin.Context) {
	var nu newUser

	// validate form input
	if err := c.Bind(&nu); err != nil {
		c.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
		return
	}

	// verify a user with ID does not already exist
	// TODO: this should go in middleware
	if existing, _ := r.userStore.User(nu.ID); existing != nil {
		c.Status(http.StatusConflict)
		return
	}

	// create new user
	u := fountain.User{
		ID:       nu.ID,
		FullName: nu.FullName,
		Email:    nu.Email,
	}

	// insert new user into store
	if err := r.userStore.PutUser(&u); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
	} else {
		loc := fmt.Sprintf("%s/%s", c.Request.URL.Path, u.ID)
		c.Redirect(http.StatusSeeOther, loc)
	}
}
