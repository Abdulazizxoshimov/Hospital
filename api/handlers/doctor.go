package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/Abdulazizxoshimov/Hospital/pkg/time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

// @Security BearerAuth
// @Summary Create Doctor
// @Description CreateDoctor
// @Tags doctors
// @Accept json
// @Produce json
// @Param doctor body entity.Doctor true "Doctor info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /doctor [post]
func (h *HandlerV1) CreateDoctor(c *gin.Context) {
	var body entity.Doctor
	pp.Println("eeeeeeeeeee")

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid request body"})
		log.Println(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	startTime, EndTime, err := time.ParseWorkTime(body.Working_hour)
	if err != nil {
		println(err.Error())
	}

	pp.Println(startTime, EndTime)
	body.ID = uuid.New().String()
	doctor, err := h.Service.Doctor().Create(ctx, &body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to create doctor"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, entity.UserCreateResponse{
		ID: doctor.ID,
	})
}

// @Security BearerAuth
// @Summary Get Doctor
// @Description Get for doctor
// @Tags doctors
// @Accept json
// @Produce json
// @Param id path string true "Doctor ID"
// @Success 200 {object} entity.Doctor
// @Failure 404 {object} entity.Error
// @Router /doctor/{id} [get]
func (h *HandlerV1) GetDoctor(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	doctor, err := h.Service.Doctor().Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{Message: "Doctor not found"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, doctor)
}

// @Security BearerAuth
// @Summary Update Doctor
// @Description Update for doctor
// @Tags doctors
// @Accept json
// @Produce json
// @Param doctor body entity.Doctor true "Updated Doctor Info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /doctor [put]
func (h *HandlerV1) UpdateDoctor(c *gin.Context) {
	var body entity.Doctor

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid request body"})
		log.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	doctor, err := h.Service.Doctor().Update(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to update doctor"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, entity.UserCreateResponse{
		ID: doctor.ID,
	})
}

// @Summary Delete Doctor
// @Description Delete for doctor
// @Tags doctors
// @Accept json
// @Produce json
// @Param id path string true "Doctor ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /doctor/{id} [delete]
func (h *HandlerV1) DeleteDoctor(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	if err := h.Service.Doctor().Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to delete doctor"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}

// @Security BearerAuth
// @Summary List Doctors
// @Description get Doctors list
// @Tags doctors
// @Accept json
// @Produce json
// @Param  page query string true "Page"
// @Param  limit query string true "Limit"
// @Param   name query string false "Name"
// @Param   specialization query string false "Specialization"
// @Success 200 {object} entity.ListDoctorRes
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /doctors [get]
func (h *HandlerV1) ListDoctors(c *gin.Context) {
	req := &entity.ListRequest{
		Filter: make(map[string]string),
	}
	page := c.Query("page")
	limit := c.Query("limit")
	name := c.Query("name")
	specialization := c.Query("specialization")
	var err error

	req.Offset, err = strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	req.Limit, err = strconv.Atoi(limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	if name != "" {
		req.Filter["name"] = name
	}
	if specialization != "" {
		req.Filter["specialization"] = specialization
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid request parameters"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	doctors, err := h.Service.Doctor().List(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to fetch doctors"})
		return
	}

	c.JSON(http.StatusOK, doctors)
}
