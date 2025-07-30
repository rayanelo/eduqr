package services

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/pkg/utils"
	"errors"
	"time"
)

type UserService struct {
	userRepo      *repositories.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewUserService(userRepo *repositories.UserRepository, jwtSecret string, jwtExpiration time.Duration) *UserService {
	return &UserService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

func (s *UserService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Address:   req.Address,
		Avatar:    "/assets/images/avatars/default-avatar.png",
		Role:      "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) Login(req *models.LoginRequest) (string, *models.UserResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return "", nil, err
	}

	return token, s.toUserResponse(user), nil
}

func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

func (s *UserService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.ContactEmail != "" {
		user.ContactEmail = req.ContactEmail
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) GetAllUsers() ([]*models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var userResponses []*models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, s.toUserResponse(&user))
	}

	return userResponses, nil
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Address:   req.Address,
		Avatar:    "/assets/images/avatars/default-avatar.png",
		Role:      req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(user.ID)
}

func (s *UserService) UpdateUserRole(id uint, role string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Role = role

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) UpdateProfile(id uint, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ContactEmail != "" {
		user.ContactEmail = req.ContactEmail
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) ChangePassword(id uint, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify current password
	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		return errors.New("mot de passe actuel incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword

	return s.userRepo.Update(user)
}

func (s *UserService) toUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		ContactEmail: user.ContactEmail,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Address:      user.Address,
		Avatar:       user.Avatar,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
