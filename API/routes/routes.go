package routes

import (
	"database/sql"
	"net/http"

	"forum/handlers"
	"forum/middleware"
	"forum/repository"
	"forum/repository/notification"
	"forum/repository/session"
	"forum/repository/user"
)

func SetupRoutes(db *sql.DB) http.Handler {
	// Create repositories
	userRepo := user.NewUserRepository(db)
	sessionRepo := session.NewSessionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	reactionRepo := repository.NewReactionRepository(db)
	imageRepo := repository.NewImageRepository(db)
	notificationRepo := notification.NewRepository(db)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userRepo, sessionRepo)
	oauthHandler := handlers.NewOAuthHandler(userRepo, sessionRepo, authHandler)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo, postRepo, imageRepo)
	postHandler := handlers.NewPostHandler(postRepo)
	myPostsHandler := handlers.NewMyPostsHandler(postRepo, commentRepo, reactionRepo, imageRepo)
	likedPostsHandler := handlers.NewLikedPostsHandler(postRepo, commentRepo, reactionRepo, imageRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo, postRepo, notificationRepo, userRepo)
	reactionHandler := handlers.NewReactionHandler(reactionRepo, postRepo, notificationRepo, userRepo)
	imageHandler := handlers.NewImageHandler(imageRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationRepo)
	guestHandler := handlers.NewGuestHandler(categoryRepo, postRepo, commentRepo, reactionRepo, imageRepo)

	// Create middleware
	registerLimiter := middleware.NewRateLimiter()
	authMiddleware := middleware.NewAuthMiddleware(sessionRepo, userRepo)
	// Corrected: CSRF is a method on AuthMiddleware, not a standalone function
	// csrfMiddleware is now directly authMiddleware.CSRF
	corsMiddleware := middleware.NewCORSMiddleware("http://localhost:8081")

	// Create router
	mux := http.NewServeMux()

	// Serve uploaded images from the API container
	fs := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes
	mux.Handle("/forum/api/categories", corsMiddleware.Handler(http.HandlerFunc(categoryHandler.GetCategories)))
	mux.Handle("/forum/api/category", corsMiddleware.Handler(http.HandlerFunc(categoryHandler.GetCategoryByID)))
	mux.Handle("/forum/api/feed", corsMiddleware.Handler(http.HandlerFunc(guestHandler.GetGuestData)))

	// Authentication routes (guest only)
	guestOnly := func(h http.Handler) http.Handler {
		return corsMiddleware.Handler(authMiddleware.RequireGuest(h))
	}

	mux.Handle("/forum/api/register", guestOnly(http.HandlerFunc(registerLimiter.Limit(authHandler.Register))))
	mux.Handle("/forum/api/session/login", guestOnly(http.HandlerFunc(authHandler.Login)))

	// OAuth routes (guest only)
	mux.Handle("/auth/google/login", guestOnly(http.HandlerFunc(oauthHandler.GoogleLogin)))
	// OAuth callbacks do not need RequireGuest or CSRF, as the 'state' parameter handles CSRF
	mux.Handle("/auth/google/callback", corsMiddleware.Handler(http.HandlerFunc(oauthHandler.GoogleCallback)))
	mux.Handle("/auth/github/login", guestOnly(http.HandlerFunc(oauthHandler.GitHubLogin)))
	mux.Handle("/oauth/github/callback", corsMiddleware.Handler(http.HandlerFunc(oauthHandler.GitHubCallback)))

	// Session management routes
	mux.Handle("/forum/api/session/logout", corsMiddleware.Handler(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("/forum/api/session/verify", corsMiddleware.Handler(http.HandlerFunc(authHandler.VerifySession)))

	// Protected routes with CSRF
	protected := func(h http.Handler) http.Handler {
		// Ensure CSRF middleware is active for protected routes
		return corsMiddleware.Handler(authMiddleware.RequireAuth(authMiddleware.CSRF(h)))
		// Temporarily for testing — remove authMiddleware.CSRF(h) if you want to bypass CSRF,
		// but remember to re-enable it for security.
		// return corsMiddleware.Handler(authMiddleware.RequireAuth(h))
	}

	// Protected user routes
	mux.Handle("/forum/api/posts/create", protected(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("/forum/api/posts/delete/", protected(http.HandlerFunc(postHandler.DeletePost)))            // DELETE /forum/api/posts/delete/{id}
	mux.Handle("/forum/api/posts/edit-title/", protected(http.HandlerFunc(postHandler.EditPostTitle)))     // PUT /forum/api/posts/edit-title/{id}
	mux.Handle("/forum/api/posts/edit-content/", protected(http.HandlerFunc(postHandler.EditPostContent))) // PUT /forum/api/posts/edit-content/{id}
	mux.Handle("/forum/api/user/posts", protected(http.HandlerFunc(myPostsHandler.GetMyPosts)))
	mux.Handle("/forum/api/user/liked", protected(http.HandlerFunc(likedPostsHandler.GetLikedPosts)))
	mux.Handle("/forum/api/user/disliked", protected(http.HandlerFunc(likedPostsHandler.GetDislikedPosts)))
	mux.Handle("/forum/api/comments/create", protected(http.HandlerFunc(commentHandler.CreateComment)))
	mux.Handle("/forum/api/comments/edit/", protected(http.HandlerFunc(commentHandler.EditComment)))     // PUT /forum/api/comments/edit/{id}
	mux.Handle("/forum/api/comments/delete/", protected(http.HandlerFunc(commentHandler.DeleteComment))) // DELETE /forum/api/comments/delete/{id}
	mux.Handle("/forum/api/react", protected(http.HandlerFunc(reactionHandler.CreateReact)))
	mux.Handle("/forum/api/images/upload", protected(http.HandlerFunc(imageHandler.Upload)))
	mux.Handle("/forum/api/user/commented", protected(http.HandlerFunc(myPostsHandler.GetCommentedPosts)))
	mux.Handle("/forum/api/images/delete/", protected(http.HandlerFunc(imageHandler.DeleteImagesByPost))) // DELETE /forum/api/images/delete/{post_id}

	// Notification routes
	mux.Handle("/forum/api/user/notifications", protected(http.HandlerFunc(notificationHandler.GetUserNotifications)))
	mux.Handle("/forum/api/notifications/read/", protected(http.HandlerFunc(notificationHandler.MarkRead))) // POST /forum/api/notifications/read/{id}
	mux.Handle("/forum/api/notifications/read-all", protected(http.HandlerFunc(notificationHandler.MarkAllRead)))
	mux.Handle("/forum/api/notifications/delete/", protected(http.HandlerFunc(notificationHandler.Delete))) // DELETE /forum/api/notifications/delete/{id}
	mux.Handle("/forum/api/notifications/delete-all", protected(http.HandlerFunc(notificationHandler.DeleteAll)))

	// Additional protected routes for user management
	mux.Handle("/forum/api/user/profile", protected(http.HandlerFunc(authHandler.GetProfile)))
	mux.Handle("/forum/api/session/logout-all", protected(http.HandlerFunc(authHandler.LogoutAll)))

	return authMiddleware.Authenticate(mux)

}
