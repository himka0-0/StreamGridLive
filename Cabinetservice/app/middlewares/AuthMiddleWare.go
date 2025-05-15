package middlewares

import (
	"Cabinetservice/app/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func validateToken(token string) (*models.Respauth, error) {
	req, _ := http.NewRequest("POST", os.Getenv("AUTH_URL")+"/validation", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("validation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validation returned status %d", resp.StatusCode)
	}

	var result models.Respauth
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func refreshToken(c *gin.Context, refreshToken string) (string, error) {
	body, _ := json.Marshal(gin.H{})
	req, _ := http.NewRequest("POST", os.Getenv("AUTH_URL")+"/refresh", bytes.NewReader(body))
	req.AddCookie(&http.Cookie{
		Name:   "refresh_token",
		Value:  refreshToken,
		Path:   "/",
		Domain: os.Getenv("DOMAIN"),
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("refresh returned status %d", resp.StatusCode)
	}

	// Ищем новый access_token в Set-Cookie
	for _, ck := range resp.Cookies() {
		if ck.Name == "access_token" {
			http.SetCookie(c.Writer, ck)
			return ck.Value, nil
		}
	}
	return "", fmt.Errorf("no access_token in response cookies")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if accessToken, err := c.Cookie("access_token"); err == nil {
			if authRes, err := validateToken(accessToken); err == nil {
				c.Set("email", authRes.Email)
				c.Set("username", authRes.Name)
				c.Next()
				return
			}
		}

		if refreshTokenValue, err := c.Cookie("refresh_token"); err == nil && refreshTokenValue != "" {
			if newAccess, err := refreshToken(c, refreshTokenValue); err == nil {
				if authRes, err := validateToken(newAccess); err == nil {
					c.Set("email", authRes.Email)
					c.Set("username", authRes.Name)
					c.Next()
					return
				}
			}
		}

		c.Redirect(http.StatusFound, os.Getenv("AUTH_URL")+"/login")
		c.Abort()
	}
}
