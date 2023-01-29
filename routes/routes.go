package routes

import (
	"fmt"
	"gobooklibrary/controllers"
	"gobooklibrary/middlewares"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func SetupAdminAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddlewareAdmin, err := jwt.New(middlewares.GinJwtMiddlewareHandler(middlewares.Admin))

	if err != nil {
		return nil, fmt.Errorf("New authMiddlewareAdmin Error:" + err.Error())
	}

	errInit := authMiddlewareAdmin.MiddlewareInit()

	if errInit != nil {
		return nil, fmt.Errorf("authMiddlewareAdmin.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddlewareAdmin, nil
}

func SetupLibrarianAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddlewareLibrarian, err := jwt.New(middlewares.GinJwtMiddlewareHandler(middlewares.Librarian))

	if err != nil {
		return nil, fmt.Errorf("New authMiddlewareLibrarian Error:" + err.Error())
	}

	errInit := authMiddlewareLibrarian.MiddlewareInit()

	if errInit != nil {
		return nil, fmt.Errorf("authMiddlewareLibrarian.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddlewareLibrarian, nil
}

func SetupUserAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddlewareUser, err := jwt.New(middlewares.GinJwtMiddlewareHandler(middlewares.User))

	if err != nil {
		return nil, fmt.Errorf("New authMiddlewareUser Error:" + err.Error())
	}

	errInit := authMiddlewareUser.MiddlewareInit()

	if errInit != nil {
		return nil, fmt.Errorf("authMiddlewareUser.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddlewareUser, nil
}

func SetupRouter() *gin.Engine {
	var err error
	r := gin.Default()

	// Admin Auth Middleware
	authMiddlewareAdmin, err := SetupAdminAuth()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Librarian Auth Middleware
	authMiddlewareLibrarian, err := SetupLibrarianAuth()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// User Auth Middleware
	authMiddlewareUser, err := SetupUserAuth()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Route groups
	bookAdmin := r.Group("/book-library/admin")
	bookLibrarian := r.Group("/book-library/librarian")
	bookUser := r.Group("/book-library/customer")

	// Admin
	bookAdmin.POST("/login", authMiddlewareAdmin.LoginHandler)
	bookAdmin.GET("/token", authMiddlewareAdmin.RefreshHandler)

	// Librarian
	bookLibrarian.POST("/login", authMiddlewareLibrarian.LoginHandler)
	bookLibrarian.GET("/token", authMiddlewareLibrarian.RefreshHandler)

	// User
	bookUser.POST("/login", authMiddlewareUser.LoginHandler)
	bookUser.POST("/token", authMiddlewareUser.RefreshHandler)

	bookAdmin.Use(authMiddlewareAdmin.MiddlewareFunc())
	{
		bookAdmin.GET("user", controllers.GetUsers)
		bookAdmin.POST("user", controllers.AddUser)

		bookAdmin.GET("employee", controllers.GetEmployees)
		bookAdmin.POST("employee", controllers.AddEmployee)
	}

	bookLibrarian.Use(authMiddlewareLibrarian.MiddlewareFunc())
	{
		bookLibrarian.POST("book", controllers.AddBook)
		bookLibrarian.GET("book", controllers.GetBooks)

		bookLibrarian.GET("category", controllers.GetCategories)
		bookLibrarian.POST("category", controllers.AddCategory)

		bookLibrarian.GET("borrow", controllers.GetBorrows)
		bookLibrarian.POST("borrow", controllers.AddBorrow)

		bookLibrarian.POST("return", controllers.ReturnBook)

		bookLibrarian.GET("copy", controllers.GetCopies)

		bookLibrarian.GET("avaliable", controllers.GetAvaliable)

		bookLibrarian.POST("copy", controllers.AddBookCopy)

		bookLibrarian.POST("renew", controllers.RenewBook)
	}

	bookUser.Use(authMiddlewareUser.MiddlewareFunc())
	{
		bookUser.GET("book", controllers.GetBooks)

		bookUser.GET("category", controllers.GetCategories)

		bookUser.GET("avaliable", controllers.GetAvaliable)

		bookUser.POST("renew", controllers.RenewBook)

		//bookUser.POST("info", controllers.GetUserProfileInfo)
	}

	return r
}
