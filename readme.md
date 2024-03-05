# GitHub OAuth2 Authentication in Go

This application demonstrates how to implement GitHub OAuth2 authentication in Go using the Gin web framework. It includes route protection and a random state generator to prevent CSRF attacks.

## Overview

The application uses the `gin-gonic/gin` package to create a web server and handle HTTP requests. It uses the `gorilla/sessions` package to manage sessions and the `golang.org/x/oauth2` package to handle the OAuth2 flow.

The application has several routes, some of which are protected by middleware that checks if the user is authenticated. If the user is not authenticated, they are redirected to the GitHub login page.

## Key Features

- **GitHub OAuth2 Authentication**: The application uses the OAuth2 flow to authenticate users with GitHub. This involves redirecting the user to the GitHub login page, where they can authorize the application to access their GitHub account.

- **Route Protection**: The application uses middleware to protect certain routes. If a user tries to access a protected route and they are not authenticated, they are redirected to the GitHub login page.

- **Random State Generator**: The application generates a random state parameter for each authentication request. This state parameter is stored in the user's session and is sent to GitHub as part of the authentication request. When GitHub redirects the user back to the application, it includes the state parameter in the redirect URL. The application then compares the state parameter in the URL with the state parameter in the user's session to ensure they match. This prevents CSRF attacks by ensuring that each authentication response corresponds to a valid authentication request.

## Packages Used

- `gin-gonic/gin`: A HTTP web framework written in Go (Golang).
- `gorilla/sessions`: Provides cookie and filesystem sessions and infrastructure for custom session backends.
- `golang.org/x/oauth2`: The Go language implementation of OAuth2.
- `golang.org/x/oauth2/github`: The GitHub OAuth2 endpoint.

## Running the Application

To run the application, you need to set the `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` environment variables to your GitHub OAuth2 client ID and secret, respectively. Then, you can run the application using the `go run` command:

```bash
go run main.go
```

The application will start a web server on localhost:8080. You can then navigate to http://localhost:8080/ in your web browser to use the application.

## Conclusion
This application demonstrates how to implement GitHub OAuth2 authentication in Go, protect routes with middleware, and prevent CSRF attacks with a random state generator. It provides a solid foundation for any application that needs to authenticate users with GitHub.