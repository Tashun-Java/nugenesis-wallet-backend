package auth

import (
	"encoding/json"
)
import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConf = &oauth2.Config{
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
	RedirectURL:  "http://localhost:8080/api/auth/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func Login(c *gin.Context) {
	url := oauthConf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(c *gin.Context) {
	code := c.Query("code")

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Example: Fetch user profile
	client := oauthConf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}
