// Some convenient utility for usage of gin framework
//
// JSON service
//
// 	ginConfig := &GinConfig{
// 		Mode: gin.ReleaseMode,
// 		Host: "localhost",
// 		Port: 8080,
// 	}
// 	engine := NewDefaultJsonEngine(ginConfig)
//
// 	// Start service
// 	// StartServiceOrExit(engine, ginConfig)
//
// 	// Binds the engine into existing HTTP service
// 	http.Handle("/root-service", engine)
//
// Panic in Code
//
// By using of "NewDefaultJsonEngine()", any panic code would be output as:
//
// 	{
// 		"http_status": 500,
// 		"error_code": -1,
//		"error_message": fmt.Sprintf("%v", panicObject),
// 	}
//
// And the HTTP engine would keep running.
//
// Special type of panic object
//
// By "DefaultPanicProcessor()", some types of object, which is panic, would be treat specially:
//
// ValidationError - Generated by "ConformAndValidateStruct()", gives "400 Bad Request" and JSON:
// 	{
// 		"http_status": 400,
// 		"error_code": -1,
// 		"error_message": errObject.Error(),
// 	}
//
// BindJsonError - Generated by "BindJson()", gives "400 Bad Request" and JSON:
// 	{
// 		"http_status": 400,
// 		"error_code": -101,
// 		"error_message": errObject.Error(),
// 	}
//
// DataConflictError - You could panic DataConflictError, which would be output as JSON:
//
// 	{
// 		"http_status": 409
// 		"error_code": errObject.ErrorCode,
// 		"error_message": errObject.ErrorMessage,
// 	}
//
// NotFound 404
//
// When not found occurs, output following JSON:
// 	{
// 		"http_status": 404,
// 		"error_code": -1,
// 		"uri": c.Request.RequestURI,
// 	}
package gin

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
)

// Configuration defines the properties on gin framework
type GinConfig struct {
	// The mode of gin framework
	// 	const (
	// 		DebugMode   string = "debug"
	// 		ReleaseMode string = "release"
	// 		TestMode    string = "test"
	// 	)
	Mode string
	// The host could be used to start service(optional)
	Host string
	// The post could be used to start service(optional)
	Port uint16
}

// Gets the address of url
func (config *GinConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// Same as "GetAddress()"
func (config *GinConfig) String() string {
	return config.GetAddress()
}

var corsConfig cors.Config

func init() {
	headers := []string{
		"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Cache-Control", "X-Requested-With",
		"accept", "origin", "Apitoken",
		"page-size", "page-pos", "order-by", "page-ptr", "total-count", "page-more", "previous-page", "next-page",
	}

	corsConfig = cors.Config{
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE", "UPDATE"},
		AllowHeaders:     headers,
		ExposeHeaders:    headers,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	corsConfig.AllowAllOrigins = true
}

// Initialize a router with default JSON response
//
// 	1. The panic code would not cause process to dead
// 	2. Use gin-contrib/cors as middleware for cross-site issue
// 	3. Change (*gin.Engine).NoRoute() with JSON output
// 	4. Change (*gin.Engine).NoMethod() with JSON output
//
// CORS Setting
//
// 	Access-Control-Allow-Origin: *
// 	Access-Control-Allow-Credentials: true
// 	Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With,
// 		accept, origin,
// 		page-size, page-pos, order-by, page-ptr, previous-page, next-page, page-more, total-count
// 	Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT
//  Access-Control-Max-Age": "43200"
func NewDefaultJsonEngine(config *GinConfig) *gin.Engine {
	gin.SetMode(config.Mode)

	router := gin.New()

	router.Use(cors.New(corsConfig))

	router.NoRoute(JsonNoRouteHandler)
	router.NoMethod(JsonNoMethodHandler)

	router.Use(BuildJsonPanicProcessor(DefaultPanicProcessor))

	return router
}

// Try to start the engine with configuration of gin
//
// If some error happened, exit application with "os.Exit(1)"
func StartServiceOrExit(router *gin.Engine, config *GinConfig) {
	if err := router.Run(config.GetAddress()); err != nil {
		logger.Errorf("Cannot start web service: %v", err)
		os.Exit(1)
	}
}
