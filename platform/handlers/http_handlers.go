package handlers

import (
	_ "PragatiIot/platform/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"

	"PragatiIot/platform/middleware"
	"PragatiIot/platform/models"
	"PragatiIot/platform/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser @Summary Register new user
// @Description Register a new user with username, password, and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User required "User Info"
// @Success 201 {object} models.ApiResponse "User registered successfully with token"
// @Failure 400 {object} models.ApiResponse "Invalid request payload"
// @Failure 500 {object} models.ApiResponse "Failed to register user or hash password"
// @Router /register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, models.ApiResponse{Error: "Invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	if err := h.userService.AddUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{Message: "User registered successfully"})
}

// LoginUser logs in a user with username and password to receive a token
// @Summary User login
// @Description Login with username and password to receive a token
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User required "User Credentials"
// @Success 200 {object} models.ApiResponse "message": "Login successful", "token": "JWT_TOKEN"
// @Failure 400 {object} models.ApiResponse "Invalid request payload"
// @Failure 401 {object} models.ApiResponse "Invalid username or password"
// @Router /login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{Error: "Invalid request payload"})
		return
	}

	dbUser, err := h.userService.GetUserByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{Error: "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.PasswordHash)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{Error: "Invalid username or password"})
		return
	}

	token, err := middleware.CreateToken(dbUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{Message: "Login successful", Token: token})
}

type HomeHandler struct {
	homeService *services.HomeService
}

func NewHomeHandler(homeService *services.HomeService) *HomeHandler {
	return &HomeHandler{homeService: homeService}
}

// AddHome adds a new home to the system
// @Summary Add a home
// @Description Adds a new home to the system
// @Tags homes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param home body models.Home required "Home Info"
// @Success 201 {object} map[string]string "Home added successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Failed to add home"
// @Router /auth/home [post]
func (h *HomeHandler) AddHome(c *gin.Context) {
	var home models.Home
	if err := c.ShouldBindJSON(&home); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.homeService.AddHome(home); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add home"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Home added successfully"})
}

// AddUserToHome adds a user to a home with a specific role
// @Summary Add user to home
// @Description Adds a user to a home with a specific role
// @Tags homes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body models.AddUserToHomeRequest true "Home and User Info"
// @Success 200 {object} models.ApiResponse "User added to home successfully"
// @Failure 400 {object} models.ApiResponse "Invalid request payload"
// @Failure 500 {object} models.ApiResponse "Failed to add user to home"
// @Router /auth/home/add-user [post]
func (h *HomeHandler) AddUserToHome(c *gin.Context) {
	var req models.AddUserToHomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{Error: "Invalid request payload"})
		return
	}

	if err := h.homeService.AddUserToHome(req.HomeID, req.UserID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to add user to home"})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{Message: "User added to home successfully"})
}

// GetHomesByUserID retrieves homes associated with a specific user ID
// @Summary Get homes by user ID
// @Description Retrieves homes associated with a specific user ID
// @Tags homes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id query int true "User ID"
// @Success 200 {array} models.Home "List of homes associated with the user ID"
// @Failure 400 {object} models.ApiResponse "Invalid user ID provided"
// @Failure 500 {object} models.ApiResponse "Failed to retrieve homes due to a server error"
// @Router /auth/home/list [get]
func (h *HomeHandler) GetHomesByUserID(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{Error: "Invalid user ID"})
		return
	}

	homes, err := h.homeService.GetHomesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to get homes"})
		return
	}

	c.JSON(http.StatusOK, homes)
}

type DeviceHandler struct {
	deviceService *services.DeviceService
	homeService   *services.HomeService
}

func NewDeviceHandler(deviceService *services.DeviceService, homeService *services.HomeService) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService, homeService: homeService}
}

// AddDevice adds a new device to the system
// @Summary Add a device
// @Description Adds a new device to the system
// @Tags devices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param device body models.Device required "Device Info"
// @Success 201 {object} map[string]string "Device added successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Failed to add device"
// @Router /auth/device [post]
func (h *DeviceHandler) AddDevice(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.deviceService.AddDevice(device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Device added successfully"})
}

// AssignDeviceToHome assigns a device to a specified home
// @Summary Assign device to home
// @Description Assigns a device to a specified home
// @Tags devices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body models.AssignDeviceRequest true "Device and Home IDs"
// @Success 200 {object} models.ApiResponse "Device assigned to home successfully"
// @Failure 400 {object} models.ApiResponse "Invalid request payload"
// @Failure 500 {object} models.ApiResponse "Failed to assign device to home"
// @Router /auth/device/assign-home [post]
func (h *DeviceHandler) AssignDeviceToHome(c *gin.Context) {
	var req models.AssignDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.deviceService.AssignDeviceToHome(req.DeviceID, req.HomeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign device to home"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device assigned to home successfully"})
}

// GetDevicesByUserID retrieves devices associated with a specific user ID
// @Summary Get devices by user ID
// @Description Retrieves devices associated with a specific user ID, returning a list of devices.
// @Tags devices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id query int true "User ID" // Make sure to clearly state that this parameter is required.
// @Success 200 {array} models.Device "List of devices associated with the user ID"
// @Failure 400 {object} models.ApiResponse "Invalid user ID provided"
// @Failure 500 {object} models.ApiResponse "Failed to retrieve devices due to server error"
// @Router /auth/device/list [get]
func (h *DeviceHandler) GetDevicesByUserID(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{Error: "Invalid user ID"})
		return
	}

	devices, err := h.deviceService.GetDevicesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to get devices"})
		return
	}

	c.JSON(http.StatusOK, devices)
}

type AnalyticsHandler struct {
	deviceService *services.DeviceService
	homeService   *services.HomeService
}

func NewAnalyticsHandler(deviceService *services.DeviceService, homeService *services.HomeService) *AnalyticsHandler {
	return &AnalyticsHandler{deviceService: deviceService, homeService: homeService}
}

// GetDeviceAnalytics retrieves analytics for a specific device within a specified home
// @Summary Get device analytics
// @Description Retrieves analytics for a specific device within a specified home, including usage statistics or other relevant data.
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param device_id query string true "Device ID required for fetching analytics"
// @Param home_id query int true "Home ID required for contextual analytics within a specific home"
// @Success 200 {object} map[string]interface{} "Analytics data for the specified device within the given home."
// @Failure 400 {object} models.ApiResponse "Invalid home or device ID provided."
// @Failure 401 {object} models.ApiResponse "Unauthorized access attempt detected."
// @Failure 500 {object} models.ApiResponse "Internal server error while retrieving device analytics."
// @Router /auth/device-analytics [get]
func (h *AnalyticsHandler) GetDeviceAnalytics(c *gin.Context) {
	deviceID := c.Query("device_id")
	homeIDStr := c.Query("home_id")
	homeID, err := strconv.Atoi(homeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{Error: "Invalid home ID"})
		return
	}

	username, _ := c.Get("username")
	user, err := h.homeService.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{Error: "Invalid username or password"})
		return
	}

	roleID, err := h.homeService.GetHomeUserRole(homeID, user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{Error: "Unauthorized to access analytics"})
		return
	}

	analytics, err := h.deviceService.GetDeviceAnalytics(deviceID, homeID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{Error: "Failed to get device analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func SetupRoutes(router *gin.Engine, userHandler *UserHandler, homeHandler *HomeHandler, deviceHandler *DeviceHandler, analyticsHandler *AnalyticsHandler) {
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth", middleware.JWTAuthMiddleware())
	{
		auth.POST("/home", homeHandler.AddHome)
		auth.POST("/home/add-user", homeHandler.AddUserToHome)
		auth.GET("/home/list", homeHandler.GetHomesByUserID)

		auth.POST("/device", deviceHandler.AddDevice)
		auth.POST("/device/assign-home", deviceHandler.AssignDeviceToHome)
		auth.GET("/device/list", deviceHandler.GetDevicesByUserID)

		auth.GET("/device-analytics", analyticsHandler.GetDeviceAnalytics)

	}
}
