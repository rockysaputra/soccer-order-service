package user

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	errWrap "user-service/common/error"
	errConstant "user-service/constants/error"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)

	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()

	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))

	if err != nil {
		return nil, err
	}

	returnData := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return returnData, nil
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// check if username already exist
	checkUser, err := u.repository.GetUser().FindByUsername(ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if checkUser != nil {
		return nil, errWrap.WrapError(errConstant.ErrUsernameExist)
	}

	checkEmail, err := u.repository.GetUser().FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if checkEmail != nil {
		return nil, errWrap.WrapError(errConstant.ErrEmailExist)
	}

	if req.Password != req.ConfirmPassword {
		return nil, errWrap.WrapError(errConstant.ErrPasswordDoesNotMatch)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPasswordStr := string(hashedPassword)

	dataUser := &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPasswordStr,
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.Customer,
	}

	user, err := u.repository.GetUser().Register(ctx, dataUser)

	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        strings.ToLower(user.Role.Code),
		},
	}

	return response, nil
}

func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {

	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPassword            []byte
		user, userResult          *models.User
		err                       error
		data                      dto.UserResponse
	)

	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errWrap.WrapError(errConstant.ErrUserNotFound)
	}
	// validate username

	checkUsername, err = u.repository.GetUser().FindByUsername(ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if checkUsername != nil {
		return nil, errWrap.WrapError(errConstant.ErrUsernameExist)
	}

	checkEmail, err = u.repository.GetUser().FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if checkEmail != nil {
		return nil, errWrap.WrapError(errConstant.ErrEmailExist)
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			password = string(hashedPassword)
		}
	}

	dataUpdate := &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    &password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	userResult, err = u.repository.GetUser().Update(ctx, dataUpdate, uuid)

	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		PhoneNumber: userResult.PhoneNumber,
		Email:       userResult.Email,
	}

	return &data, nil
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		PhoneNumber: userLogin.PhoneNumber,
		Email:       userLogin.Email,
		Role:        userLogin.Role,
	}
	return &data, nil
}
func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errWrap.WrapError(errConstant.ErrUserNotFound)
	}

	responseData := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}

	return responseData, nil
}
