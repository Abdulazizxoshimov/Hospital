package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/Abdulazizxoshimov/Hospital/pkg/validation"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Security  		BearerAuth
// @Summary   		Create User
// @Description 	Api for create a new user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			user body entity.UserRegister true "Create User Model"
// @Success 		201 {object} entity.UserCreateResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	var (
		body entity.UserRegister
	)

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	filter := map[string]string{
		"email": body.Email,
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
	if status {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Email already used",
		})
		log.Println(err)
		return
	}

	statusPassword := validation.PasswordValidation(body.Password)
	if !statusPassword {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.WeakPasswordMessage,
		})
		log.Println(entity.WeakPasswordMessage)
		return
	}

	hashpassword, err := validation.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	if !validation.ValidateUsername(body.Username) {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Invalid Username",
		})
		log.Println(err)
		return
	}

	userServiceCreateResponse, err := h.Service.User().Create(ctx, &entity.User{
		ID:       uuid.New().String(),
		UserName: body.Username,
		Email:    body.Email,
		Password: hashpassword,
		Role:     "user",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, entity.UserCreateResponse{
		ID: userServiceCreateResponse.ID,
	})
}

// @Security  		BearerAuth
// @Summary   		Update User
// @Description 	Api for update a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			user body entity.UserUpdate true "Update User Model"
// @Success 		200 {object} entity.UserUpdate
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user [PUT]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	var (
		body entity.UserUpdate
	)

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	body.Email, err = validation.EmailValidation(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	filter := map[string]string{
		"id": body.ID,
	}
	user, err := h.Service.User().Get(ctx, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
	}
	mp := map[string]string{
		"email": body.Email,
	}
	if user.Email != body.Email {
		status, err := h.Service.User().CheckUnique(ctx, &entity.GetRequest{
			Filter: mp,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
		if status {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: "email already used",
			})
			return
		}
	}

	if body.PhoneNumber != "" {
		status := validation.PhoneUz(body.PhoneNumber)
		if !status {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: "phone number is invalid",
			})
			log.Println("phone number is invalid")
			return
		}
	}

	updatedUser, err := h.Service.User().Update(ctx, &entity.User{
		ID:          body.ID,
		UserName:    body.UserName,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		Role:        "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.UserUpdate{
		ID:          updatedUser.ID,
		UserName:    updatedUser.UserName,
		Email:       updatedUser.Email,
		PhoneNumber: updatedUser.PhoneNumber,
	})
}

// @Security  		BearerAuth
// @Summary   		Delete User
// @Description 	Api for delete a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} bool
// @Failure 		404 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user/{id} [DELETE]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	userID := c.Param("id")

	user, err := h.Service.User().Get(ctx, map[string]string{
		"id": userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	if user.Role == "admin" {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Wrong request",
		})
		return
	}

	err = h.Service.User().Delete(ctx, &entity.DeleteRequest{
		Id: userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Security  		BearerAuth
// @Summary   		Get User
// @Description 	Api for getting a user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} entity.UserResponse
// @Failure 		404 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user/{id} [GET]
func (h *HandlerV1) GetUser(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	userID := c.Param("id")

	filter := make(map[string]string)

	if govalidator.IsEmail(userID) {
		filter["email"] = userID
	} else if validation.ValidateUUID(userID) {
		filter["id"] = userID
	} else {
		filter["username"] = userID
	}

	response, err := h.Service.User().Get(ctx, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.UserResponse{
		ID:           userID,
		UserName:     response.UserName,
		Email:        response.Email,
		PhoneNumber:  response.PhoneNumber,
		RefreshToken: response.RefreshToken,
		Role:         response.Role,
	})
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			page query string true "Page"
// @Param 			limit query string true "Limit"
// @Success 		200 {object} entity.ListUserRes
// @Failure 		404 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/users [GET]
func (h *HandlerV1) ListUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()
	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	offset := (pageInt - 1) * limitInt
	filter := map[string]string{
		"role": "user",
	}
	listUsers, err := h.Service.User().List(ctx, &entity.ListRequest{
		Offset: offset,
		Limit:  limitInt,
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var users []*entity.User
	for _, user := range listUsers.User {
		users = append(users, &entity.User{
			ID:           user.ID,
			UserName:     user.UserName,
			Email:        user.Email,
			PhoneNumber:  user.PhoneNumber,
			RefreshToken: user.RefreshToken,
			Role:         user.Role,
		})
	}
	

	c.JSON(http.StatusOK, entity.ListUserRes{
		User: users,
		TotalCount: listUsers.TotalCount,
	})
}

// @Security        BearerAuth
// @Summary         Update Password
// @Description     Api for updating user's password
// @Tags            users
// @Accept          json
// @Produce         json
// @Param 			user body entity.UpdatePassword true "Update User Password"
// @Success 		200 {object} string
// @Failure 		404 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user/password [PUT]
func (h *HandlerV1) UpdatePassword(c *gin.Context) {
	var (
		body entity.UpdatePassword
	)

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	_, err = h.Service.User().UpdatePassword(ctx, &entity.UpdatePassword{
		UserID:      body.UserID,
		NewPassword: body.NewPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.YourProfileHasChangedSuccusfully)
}
