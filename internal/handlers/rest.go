package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"parcel-service/internal/core/domain"
	"parcel-service/internal/core/ports"
	"parcel-service/pkg/authorization"
	"parcel-service/pkg/dto"
)

type HTTPHandler struct {
	parcelService ports.ParcelService
	router        *gin.Engine
}

func NewRest(parcelService ports.ParcelService, router *gin.Engine) *HTTPHandler {
	return &HTTPHandler{
		parcelService: parcelService,
		router:        router,
	}
}

func (handler *HTTPHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/parcels", handler.GetAll)
	api.GET("/parcels/:id", handler.Get)
	api.POST("/parcels", handler.Create)
	api.DELETE("/parcels/:id", handler.CancelParcel)
	api.GET("/customers/:id/parcels", handler.GetFromCustomer)
}

func (handler *HTTPHandler) GetAll(c *gin.Context) {
	if authorization.NewRest(c).AuthorizeAdmin() {

		parcels, err := handler.parcelService.GetAll()

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, parcels)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (handler *HTTPHandler) Get(c *gin.Context) {
	auth := authorization.NewRest(c)

	parcel, err := handler.parcelService.Get(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(parcel.OwnerId) {
		c.JSON(http.StatusOK, parcel)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (handler *HTTPHandler) Create(c *gin.Context) {
	body := dto.BodyCreateParcel{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(body.OwnerId) {

		parcel, err := handler.parcelService.Create(body.OwnerId, body.Name, body.Description, domain.Dimensions(body.Size), body.Weight)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BuildResponseCreateParcel(parcel))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (handler *HTTPHandler) CancelParcel(c *gin.Context) {
	parcel, err := handler.parcelService.Get(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	auth := authorization.NewRest(c)
	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(parcel.OwnerId) {

		err := handler.parcelService.CancelParcel(parcel.ID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (handler *HTTPHandler) GetFromCustomer(c *gin.Context) {
	auth := authorization.NewRest(c)
	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {
		parcels, err := handler.parcelService.GetAllFromCustomer(c.Param("id"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, parcels)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
