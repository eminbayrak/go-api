package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
	}

	store = sessions.NewCookieStore([]byte("github-auth-session-store"))
)

func generateRandomString(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	randomString := base64.URLEncoding.EncodeToString(buffer)
	return randomString[:length], nil
}

func generateState() string {
	state, err := generateRandomString(16)
	if err != nil {
		log.Println("Error generating state value:", err)
		return "Error generating state value"
	}
	return state
}

func handleGithubLogin(c *gin.Context) {
	state := generateState()
	session, err := store.Get(c.Request, "github-auth-session")
	if err != nil {
		log.Println("Error retrieving session:", err)
		return
	}
	session.Values["state"] = state
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Println("Error saving session:", err)
		return
	}
	log.Printf("Set state: %s", state) // Log the set state
	url := githubOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGithubCallback(c *gin.Context) {
	// Retrieve the state parameters from the request and session
	state := strings.TrimSpace(c.Query("state"))
	sessionState := strings.TrimSpace(getSessionState(c))

	// Log the received state and session state for debugging
	log.Printf("Received state: %s", state)
	log.Printf("Session state: %s", sessionState)

	// Compare the received state with the session state
	if state != sessionState {
		// If the states don't match, abort the request with a 403 Forbidden error
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("invalid state parameter"))
		log.Println("Invalid state parameter")
		return
	}

	// If the states match, proceed with token exchange and session handling
	code := c.Query("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		log.Println("Failed to exchange token:", err)
		return
	}

	// Save the token in the session
	session, err := store.Get(c.Request, "github-auth-session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		log.Println("Failed to get session:", err)
		return
	}
	session.Values["token"] = token.AccessToken
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		log.Println("Failed to save session:", err)
		return
	}

	// Redirect the user to the desired destination
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getSessionState(c *gin.Context) string {
	session, err := store.Get(c.Request, "github-auth-session")
	if err != nil {
		log.Println("Error retrieving session:", err)
		return ""
	}

	state, ok := session.Values["state"].(string)
	if !ok || state == "" {
		state = generateState()
		session.Values["state"] = state
		if err := session.Save(c.Request, c.Writer); err != nil {
			log.Println("Error saving session:", err)
			return ""
		}
	}

	log.Printf("Retrieved state: %s", state)
	return state
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.URL.Path {
		case "/auth/github/login":
			handleGithubLogin(c)
			c.Abort()
		case "/auth/github/callback":
			handleGithubCallback(c)
			c.Abort()
		default:
			session, err := store.Get(c.Request, "github-auth-session")
			if err != nil || session.Values["token"] == nil {
				// No session exists, redirect to login
				c.Redirect(http.StatusTemporaryRedirect, "/auth/github/login")
				c.Abort()
			} else {
				// Session exists, proceed with request
				c.Next()
			}
		}
	}
}
