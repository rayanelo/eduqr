package controllers

import (
	"net/http"
	"strconv"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type LoginResponse struct {
	Token string               `json:"token"`
	User  *models.UserResponse `json:"user"`
}

// @Summary Register a new user
// @Description Register a new user with email, password, first name, and last name
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.Register(&req)
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User login data"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := c.userService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get user profile
// @Description Get current user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := c.userService.GetUserByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Update user profile
// @Description Update current user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(userID.(uint), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Get user by ID
// @Description Get user by ID based on role permissions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user can view this user
	if !models.CanViewRole(userRole.(string), user.Role) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	// Get viewable fields for this user
	viewableFields := models.GetViewableFields(userRole.(string), user.Role)
	filteredUser := filterUserFields(*user, viewableFields)

	ctx.JSON(http.StatusOK, filteredUser)
}

// @Summary Get all users
// @Description Get all users based on user role permissions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter users based on role permissions
	filteredUsers := []models.UserResponse{}
	for _, user := range users {
		if models.CanViewRole(userRole.(string), user.Role) {
			// Get viewable fields for this user
			viewableFields := models.GetViewableFields(userRole.(string), user.Role)
			filteredUser := filterUserFields(*user, viewableFields)
			filteredUsers = append(filteredUsers, filteredUser)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"users": filteredUsers})
}

// filterUserFields returns a UserResponse with only the specified fields
func filterUserFields(user models.UserResponse, viewableFields []string) models.UserResponse {
	filteredUser := models.UserResponse{}

	for _, field := range viewableFields {
		switch field {
		case "id":
			filteredUser.ID = user.ID
		case "email":
			filteredUser.Email = user.Email
		case "first_name":
			filteredUser.FirstName = user.FirstName
		case "last_name":
			filteredUser.LastName = user.LastName
		case "phone":
			filteredUser.Phone = user.Phone
		case "address":
			filteredUser.Address = user.Address
		case "avatar":
			filteredUser.Avatar = user.Avatar
		case "role":
			filteredUser.Role = user.Role
		case "created_at":
			filteredUser.CreatedAt = user.CreatedAt
		case "updated_at":
			filteredUser.UpdatedAt = user.UpdatedAt
		}
	}

	return filteredUser
}

// @Summary Create user
// @Description Create a new user based on role permissions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.CreateUserRequest true "User creation data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user can manage the target role
	if !models.CanManageRole(userRole.(string), req.Role) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to create user with this role"})
		return
	}

	user, err := c.userService.CreateUser(&req)
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// @Summary Update user
// @Description Update user by ID based on role permissions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get target user to check permissions
	targetUser, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user can manage the target user
	if !models.CanManageRole(userRole.(string), targetUser.Role) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to update this user"})
		return
	}

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Delete user by ID based on role permissions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get target user to check permissions
	targetUser, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user can manage the target user
	if !models.CanManageRole(userRole.(string), targetUser.Role) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to delete this user"})
		return
	}

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// @Summary Update user role
// @Description Update user role by ID based on role permissions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param role body map[string]string true "Role update data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/role [patch]
func (c *UserController) UpdateUserRole(ctx *gin.Context) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the new role
	if !models.IsValidRole(req.Role) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	// Get target user to check permissions
	targetUser, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user can manage the target user
	if !models.CanManageRole(userRole.(string), targetUser.Role) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to update this user's role"})
		return
	}

	// Check if user can manage the new role
	if !models.CanManageRole(userRole.(string), req.Role) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to assign this role"})
		return
	}

	user, err := c.userService.UpdateUserRole(uint(id), req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
