package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/Abdulazizxoshimov/Hospital/pkg/gmail"
	"github.com/Abdulazizxoshimov/Hospital/pkg/token"
	"github.com/Abdulazizxoshimov/Hospital/pkg/validation"
	govalidator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/spf13/cast"
)

// @Summary 		Register
// @Description 	Api for register user
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			User body entity.UserRegister true "Register User"
// @Success 		200 {object} string
// @Failure 		400 {object} entity.Error
// @Failure         409 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/register [POST]
func (h *HandlerV1) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var (
		body entity.UserRegister
	)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	pp.Println(body.Email)
	valid := govalidator.IsEmail(body.Email)
	if !valid {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Bad email",
		})
		log.Println(err)
		return
	}

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	status := validation.PasswordValidation(body.Password)
	if !status {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Password should be 8-20 characters long and contain at least one lowercase letter, one uppercase letter, and one digit",
		})
		log.Println(err)
		return
	}

	usernameStatus := validation.ValidateUsername(body.Username)
	if !usernameStatus {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Username is invalid!!!",
		})
		log.Println("Username is invalid!!!")
		return
	}

	filter := map[string]string{
		"email": body.Email,
	}
	exists, err := h.Service.User().CheckUnique(ctx, &entity.GetRequest{
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: "Oops something went wrong!!!",
		})
		log.Println(err)
		return
	}

	if exists {
		c.JSON(http.StatusConflict, entity.Error{
			Message: "This email already in use:",
		})
		return
	}
	radomNumber, err := gmail.SendCodeGmail(body.Email, "Hospital\n", "./pkg/gmail/emailotp.html", h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	err = h.redisStorage.Set(ctx, radomNumber, body, time.Second*300)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, entity.VerificationCodeSentYourEmail)
}

// @Summary            Verify
// @Description        Api for verify register
// @Tags               registration
// @Accept             json
// @Produce            json
// @Param              email query string true "email"
// @Param              code query string true "code"
// @Success            201 {object} entity.UserResponse
// @Failure            400 {object} entity.Error
// @Failure            500 {object} entity.Error
// @Router             /users/verify [post]
func (h *HandlerV1) Verify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	code := c.Query("code")
	email := c.Query("email")

	userData, err := h.redisStorage.Get(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	var user entity.User

	err = json.Unmarshal(userData, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	if user.Email != email {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "The email did not match ",
		})
		log.Println(err)
		return
	}

	id := uuid.NewString()

	h.RefreshToken = token.JWTHandler{
		Sub:        id,
		Role:       "user",
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      user.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusConflict, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	hashPassword, err := validation.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	claims, err := token.ExtractClaim(access, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadGateway, entity.Error{
			Message: err.Error(),
		})
	}

	_, err = h.Service.User().Create(ctx, &entity.User{
		ID:           id,
		FullName:     user.FullName,
		UserName:     user.UserName,
		Email:        user.Email,
		Password:     hashPassword,
		RefreshToken: refresh,
		Role:         cast.ToString(claims["role"]),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	respUser := &entity.UserResponse{
		ID:           id,
		UserName:     user.UserName,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Role:         cast.ToString(claims["role"]),
		AccesToken:   access,
		RefreshToken: refresh,
	}

	c.JSON(http.StatusCreated, respUser)
}

// @Summary 		Login
// @Description 	Api for login user
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			login body entity.Login true "Login Model"
// @Success 		200 {object} entity.UserResponse
// @Failure 		400 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/login [POST]
func (h *HandlerV1) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	var body entity.Login

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	var filter map[string]string
	if govalidator.IsEmail(body.UserNameOrEmail) {
		filter = map[string]string{
			"email": body.UserNameOrEmail,
		}
	} else {
		filter = map[string]string{
			"username": body.UserNameOrEmail,
		}
	}

	response, err := h.Service.User().Get(ctx, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if !(validation.CheckHashPassword(body.Password, response.Password)) {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Incorrect Password",
		})
		return
	}

	h.RefreshToken = token.JWTHandler{
		Sub:        response.ID,
		Role:       response.Role,
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      response.Email,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	respUser := &entity.UserResponse{
		ID:           response.ID,
		FullName:     response.FullName,
		UserName:     response.UserName,
		Email:        response.Email,
		PhoneNumber:  response.PhoneNumber,
		Role:         response.Role,
		RefreshToken: refresh,
		AccesToken:   access,
	}
	_, err = h.Service.User().UpdateRefresh(ctx, &entity.UpdateRefresh{
		UserID:       response.ID,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, respUser)
}

// @Summary 		Forget Password
// @Description 	Api for sending otp
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			email path string true "Email"
// @Success 		200 {object} string
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/forgot/{email} [POST]
func (h *HandlerV1) Forgot(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	email := c.Param("email")

	email, err := validation.EmailValidation(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	filter := map[string]string{
		"email": email,
	}

	status, err := h.Service.User().CheckUnique(ctx, &entity.GetRequest{
		Filter: filter,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if !status {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "This user is not registered",
		})
		return
	}

	radomNumber, err := gmail.SendCodeGmail(email, "Univer\n", "./pkg/gmail/forgotpassword.html", h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if err := h.redisStorage.Set(ctx, radomNumber, cast.ToString(email), time.Second*300); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, "We have sent otp your email")
}

// @Summary 		Verify OTP
// @Description 	Api for verify user
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			email query string true "Email"
// @Param 			otp query string true "OTP"
// @Success 		200 {object} bool
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/verify [POST]
func (h *HandlerV1) VerifyOTP(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	otp := c.Query("otp")
	email := c.Query("email")

	userData, err := h.redisStorage.Get(ctx, otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	var redisEmail string

	err = json.Unmarshal(userData, &redisEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if redisEmail != email {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "The email did not match",
		})
		log.Println("The email did not match")
		return
	}

	c.JSON(http.StatusCreated, true)
}

// @Summary 		Reset Password
// @Description 	Api for reset password
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			User body entity.ResetPassword true "Reset Password"
// @Success 		200 {object} bool
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/reset-password [PUT]
func (h *HandlerV1) ResetPassword(c *gin.Context) {
	var (
		body entity.ResetPassword
	)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(7))
	defer cancel()

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	status := validation.PasswordValidation(body.NewPassword)
	if !status {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Password should be 8-20 characters long and contain at least one lowercase letter, one uppercase letter, one symbol, and one digit",
		})
		log.Println(err)
		return
	}

	user, err := h.Service.User().Get(ctx, map[string]string{
			"email": body.Email,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	hashPassword, err := validation.HashPassword(body.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadGateway, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	responseStatus, err := h.Service.User().UpdatePassword(ctx, &entity.UpdatePassword{
		UserID:      user.ID,
		NewPassword: hashPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	if !responseStatus.Status {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: "Password doesn't updated",
		})
		log.Println("Password doesn't updated")
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Summary 		New Token
// @Description 	Api for updated acces token
// @Tags 			registration
// @Accept 			json
// @Produce 		json
// @Param 			refresh path string true "Refresh Token"
// @Success 		200 {object} entity.TokenResp
// @Failure 		400 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		409 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/token/{refresh} [GET]
func (h *HandlerV1) Token(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	RToken := c.Param("refresh")

	user, err := h.Service.User().Get(ctx, map[string]string{
			"refresh_token": RToken,
		},
	)

	if err != nil {
		c.JSON(500, entity.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	resclaim, err := token.ExtractClaim(RToken, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(500, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	Now_time := time.Now().Unix()
	exp := (resclaim["exp"])
	if exp.(float64)-float64(Now_time) > 0 {
		h.RefreshToken = token.JWTHandler{
			Sub:        user.ID,
			Role:       user.Role,
			SigningKey: h.Config.Token.SignInKey,
			Log:        h.Logger,
			Email:      user.Email,
		}

		access, refresh, err := h.RefreshToken.GenerateAuthJWT()
		if err != nil {
			c.JSON(http.StatusConflict, entity.Error{
				Message: err.Error(),
			})
			log.Println(err)
			return
		}

		_, err = h.Service.User().UpdateRefresh(ctx, &entity.UpdateRefresh{
			UserID:       user.ID,
			RefreshToken: refresh,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: err.Error(),
			})
			log.Println(err)
			return
		}

		respUser := &entity.TokenResp{
			ID:      user.ID,
			Role:    user.Role,
			Refresh: refresh,
			Access:  access,
		}

		c.JSON(http.StatusCreated, respUser)
	} else {
		c.JSON(http.StatusUnauthorized, entity.Error{
			Message: "refresh token expired",
		})
		log.Println("refresh token expired")
		return
	}
}
