package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codepnw/go-auth-cookies/internal/models"
	"github.com/codepnw/go-auth-cookies/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JWTToken struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type SessionData struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"userId"`
}

func (h *handlerConfig) SignInHandler(c *gin.Context) {
	var authReq models.UserAuthenReq

	if err := c.ShouldBindJSON(&authReq); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate
	validateErrors := utils.ValidateUserReq(authReq)
	if len(validateErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validateErrors,
		})
		return
	}

	// Find user by email
	foundUser, err := h.db.FindUserByEmail(c, authReq.Email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "email or password is invalid",
		})
		return
	}

	if foundUser.Password != authReq.Password {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email or password is invalid",
		})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Email: authReq.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	sessionID := uuid.New().String()

	sessionData := map[string]any{
		"userId": foundUser.ID,
		"token":  tokenStr,
	}

	sessionJson, err := json.Marshal(sessionData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "fail to encode the session data",
		})
		return
	}

	// Redis set
	err = h.redisClient.Set(c, sessionID, sessionJson, time.Until(expirationTime)).Err()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "fail to save session data to the redis",
		})
		return
	}

	c.SetCookie("session_id", sessionID, int(time.Until(expirationTime).Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"expires": expirationTime,
	})
}

func (h *handlerConfig) LogoutHandler(c *gin.Context) {
	// Retrieve the session from cookies
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized request",
		})
		return
	}

	err = h.redisClient.Del(c, sessionID).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "failed to end session",
		})
		return
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successful",
	})
}

func (h *handlerConfig) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - no session",
			})
			return
		}

		sessionData, err := h.redisClient.Get(c, sessionID).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired session",
			})
			return
		}

		var data SessionData
		err = json.Unmarshal([]byte(sessionData), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to decode the session data",
			})
			return
		}

		token, err := jwt.ParseWithClaims(data.Token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("userId", data.UserID)
		c.Next()
	}
}
