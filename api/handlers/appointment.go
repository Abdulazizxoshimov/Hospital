package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// @Summary Create an appointment
// @Description API for creating a new appointment
// @Tags Appointment
// @Accept json
// @Produce json
// @Param appointment body entity.Appointment true "Appointment details"
// @Success 201 {object} entity.Appointment
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /appointment [post]
func (h *HandlerV1) CreateAppointment(c *gin.Context) {
	var appointment entity.Appointment

	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Invalid request data",
		})
		h.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	createdAppointment, err := h.Service.Appointment().CreateAppointment(ctx, &appointment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: "Failed to create appointment",
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdAppointment)
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			Appointment
// @Accept 			json
// @Produce 		json
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Success 200 {array} entity.ListAppointments
// @Failure 400 {object} entity.Error
// @Router /appointments [get]
func (h *HandlerV1) GetAppointments(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid page number"})
		h.Logger.Error(err.Error())
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid limit value"})
		h.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	listApp, totalCount, err := h.Service.Appointment().ListAppointments(ctx, page, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{Message: "Appointments not found"})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.ListAppointments{
		Appointments: listApp,
		TotalCount:   int64(totalCount),
	})
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			Appointment
// @Accept 			json
// @Produce 		json
// @Param id path int true "Appointment ID"
// @Success 200 {object} entity.Appointment
// @Failure 404 {object} entity.Error
// @Router /appointment/{id} [get]
func (h *HandlerV1) GetAppointmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid ID format"})
		h.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	appointment, err := h.Service.Appointment().GetAppointment(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{Message: "Appointment not found"})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			Appointment
// @Accept 			json
// @Produce 		json
// @Param id path int true "Appointment ID"
// @Param date query string true "New appointment date"
// @Success 200 {object} entity.UserCreateResponse
// @Failure 400 {object} entity.Error
// @Router /appointment [put]
func (h *HandlerV1) UpdateAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid ID format"})
		h.Logger.Error(err.Error())
		return
	}

	newDate := c.Query("date")
	parsedDate, err := time.Parse(time.RFC3339, newDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid date format"})
		h.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	err = h.Service.Appointment().UpdateAppointment(ctx, id, parsedDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to update appointment"})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.UserCreateResponse{ID: strconv.Itoa(id)})
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			Appointment
// @Accept 			json
// @Produce 		json
// @Param id path int true "Appointment ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} entity.Error
// @Router /appointment/{id} [delete]
func (h *HandlerV1) DeleteAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid ID format"})
		h.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	err = h.Service.Appointment().DeleteAppointment(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{Message: "Appointment not found"})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			Appointment
// @Accept 			json
// @Produce 		json
// @Param  page query string true "Page"
// @Param  limit query string true "Limit"
// @Success 200 {array} entity.ListAvailabilities
// @Failure 400 {object} entity.Error
// @Router /availabilities [get]
func (h *HandlerV1) GetDoctorAvailabilities(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		h.Logger.Error(err.Error())
		return
	}
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		h.Logger.Error(err.Error())
		return
	}
	ListAvailabilities, totalCount, err := h.Service.Appointment().ListAvailabilities(ctx, pageInt, intLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.ObjectNotFount,
		})
		h.Logger.Error(err.Error())
		return
	}
	c.JSON(http.StatusAccepted, entity.ListAvailabilities{
		Availabilities: ListAvailabilities,
		TotalCount:   int64(totalCount),
	})
}

// @Security  		BearerAuth
// @Summary   		List User
// @Description 	Api for getting list user
// @Tags 			Appointment
// @Accept 			json
// @Produce 		json
// @Param id path int true "doctor_availability ID"
// @Success 200 {object} entity.Availability
// @Failure 404 {object} map[string]string
// @Router /availability/{id} [get]
func (h *HandlerV1) GetAvailabilityByID(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), h.Config.Context.Timeout)
	defer cancel()

	id := c.Param("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: err.Error(),
		})
		h.Logger.Error(err.Error())
		return
	}
	availabilitie, err := h.Service.Appointment().GetAvailability(ctx, ID)

	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: entity.ObjectNotFount,
		})
		return
	}
	c.JSON(http.StatusOK, availabilitie)
}
