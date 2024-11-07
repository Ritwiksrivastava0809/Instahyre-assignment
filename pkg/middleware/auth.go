package middleware

import (
	"fmt"
	"net/http"
	"spam-search/pkg/constants"
	"spam-search/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthTokenMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		tokenHeader := c.GetHeader(constants.Authorization)
		if len(tokenHeader) == 0 {
			log.Error().Msg("Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			return
		}

		// Split the tokenHeader into fields
		fields := strings.Fields(tokenHeader)
		if len(fields) < 2 {
			log.Error().Msg("Invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Check if the token is prefixed with "Bearer"
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != strings.ToLower(constants.Bearer) {
			log.Error().Msgf("Unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unsupported authorization type: %s", authorizationType)})
			return
		}

		// Extract the actual token from the header
		accessToken := fields[1]

		// Parse and verify the token
		claims, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			fmt.Println("error:: ", err)
			log.Error().Err(err).Msg(constants.InvalidToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set the claims in the context for downstream use
		c.Set(constants.AuthorizationPayloadKey, claims)

		// Continue with the next middleware or handler
		c.Next()
	}
}
