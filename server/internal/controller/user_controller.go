package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/auth"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/dto"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/email"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	baseController BaseController[schema.User, dto.User]
}

func NewUserController() UserController {
	return UserController{
		baseController: BaseController[schema.User, dto.User]{
			collectionName: "users",
			displayName:    "User",
		},
	}
}

// Query godoc
// @Summary      Query users
// @Description  Retrieve users that match the specified query parameters. Supported parameters: `id`, `role`.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    query     string  false  "User ID (ObjectID)"
// @Param        role  query     string  false  "User role (e.g. company, jobSeeker)"
// @Success      200   {array}   schema.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/query [get]
func (jc UserController) Query(c *gin.Context) {
	userFilter, shouldReturn := userFilter(c)
	if shouldReturn {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := repository.FindAll[schema.User](ctx, userFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func userFilter(c *gin.Context) (bson.M, bool) {
	allowedParams := map[string]func(string) (any, error){
		"id": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("id parameter is empty")
			}
			return primitive.ObjectIDFromHex(v)
		},
		"role": func(v string) (any, error) {
			if v == "" {
				return nil, nil
			}
			return bson.M{"$eq": v}, nil
		},
	}

	filter := bson.M{}

	// Loop through query params
	for key, value := range c.Request.URL.Query() {
		if fn, ok := allowedParams[key]; ok {
			val, err := fn(value[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return nil, true
			}
			if val != nil {
				if key == "id" {
					filter["_id"] = val
				} else {
					filter[key] = val
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported query parameter: " + key})
			return nil, true
		}
	}
	return filter, false
}

// Update godoc
// @Summary      Update an existing user
// @Description  Modify user data by providing the user ID in the path and updated JSON body.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "User ID"
// @Param        user  body      schema.User  true  "Updated user data"
// @Success      200   {object}  schema.User
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{id} [put]
func (jc UserController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Create godoc
// @Summary      Create a new user
// @Description  Add a new user document to the database.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      schema.User  true  "User object"
// @Success      201   {object}  schema.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/ [post]
func (jc UserController) Create(c *gin.Context) {
	jc.baseController.Create(c)
}

// Delete godoc
// @Summary      Delete a user
// @Description  Remove a user document from the database by ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (jc UserController) Delete(c *gin.Context) {
	uid := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, _ := repository.FindOne[schema.User](ctx, uid)

	emailBody := fmt.Sprintf(
		"Dear %s,\nYour account has been deleted by the administrator. If you beleive this is a mistake, please reply to this email immediately.\nRegards,\nJob Applier 3000", user.Name,
	)

	if err := email.Send(user.Email, "User Deletion Notice", emailBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email"})
		return
	}
	jc.baseController.Delete(c)
}

// RetrieveAll godoc
// @Summary      Retrieve all users
// @Description  Fetch all user documents from the database.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200   {array}   schema.User
// @Failure      500   {object}  map[string]string
// @Router       /users/ [get]
func (jc UserController) RetrieveAll(c *gin.Context) {
	jc.baseController.RetrieveAll(c)
}

// RetrieveOne godoc
// @Summary      Get a user by ID
// @Description  Retrieve a single user document by its ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  schema.User
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func (jc UserController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}

// VerifyUser godoc
// @Summary      Verify or unverify a user
// @Description  Change only the 'verified' status of a user (admin only)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id        path      string          true  "User ID"
// @Param        verified  body      map[string]bool true  "Verified status"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/verify [patch]
func (jc UserController) VerifyUser(c *gin.Context) {
	tokenStr, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No auth token"})
		return
	}

	claims, err := auth.ParseJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	jc.baseController.Update(c)

	uid := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, _ := repository.FindOne[schema.User](ctx, uid)
	var verificationStatus string
	if user.Verified {
		verificationStatus = "verified"
	} else {
		verificationStatus = "unverified"
	}
	emailBody := fmt.Sprintf(
		"Dear %s,\nYour account has been %s.", user.Name, verificationStatus,
	)

	if err := email.Send(user.Email, "User Account Verification Notice", emailBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email"})
		return
	}

}

// EditPermission godoc
// @Summary      Change a user's role (admin only)
// @Description  Modify only the 'role' of a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string            true  "User ID"
// @Param        role  body      map[string]string true  "New role"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/role [patch]
func (jc UserController) EditPermission(c *gin.Context) {
	tokenStr, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No auth token"})
		return
	}

	claims, err := auth.ParseJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	jc.baseController.Update(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	uid := c.Param("id")
	user, _ := repository.FindOne[schema.User](ctx, uid)

	emailBody := fmt.Sprintf(
		"Dear %s, \n Your account permission has been changed to %s", user.Name, user.Role,
	)

	if err := email.Send(user.Email, "User Permission Change Notice", emailBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email"})
		return
	}
}

// GetPublicInfo godoc
// @Summary      Get minimal public info of a user/company
// @Description  Returns name, role, profile image, and userInfo (custom name/logo if company)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       //{id} [get]
func (jc UserController) GetPublicInfo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := repository.FindOne[schema.User](ctx, objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Extract custom name and logo from userInfo if available
	customName, customLogo := extractUserInfo(user.UserInfo)

	resp := gin.H{
		"name":         user.Name, // default top-level name
		"role":         user.Role,
		"profileImage": user.AvatarURL, // default avatar
	}

	// Override with custom info if available
	userInfo := gin.H{}
	if customName != "" {
		userInfo["name"] = customName
		resp["name"] = customName // optionally, override top-level name
	}
	if customLogo != "" {
		userInfo["logo"] = customLogo
		resp["profileImage"] = customLogo
	}

	if len(userInfo) > 0 {
		resp["userInfo"] = userInfo
	}

	c.JSON(http.StatusOK, resp)
}

func extractUserInfo(info bson.M) (name string, logo string) {
	name, _ = info["name"].(string)
	logo, _ = info["logo"].(string)
	return
}
