package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ayobami6/pickitup_v3/config"
	"github.com/Ayobami6/pickitup_v3/pkg/types"
	"github.com/Ayobami6/pickitup_v3/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

var UserKey contextKey= "UserID"

const RiderKey contextKey = "RiderID"


func CreateJWT(secret []byte, userId int) (string, error) {
	exp, err := strconv.ParseInt(config.GetEnv("JWT_EXPIRATION", "25000"), 10, 64)
	if err != nil {
		return "", err
	}
	expiration := time.Second * time.Duration(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}


// func RiderAuth(handlerFunc http.HandlerFunc, riderStore models.RiderRepository) gin.HandlerFunc{
// 	return func(c *gin.Context) {
// 		// get token frome request
// 		tokenString, err := utils.GetTokenFromRequest(r)
// 		if err != nil {
// 			log.Println("TokenString error: ", err)
// 			utils.WriteError(w, http.StatusInternalServerError, "Don't Panic This is From Us!")
// 			return
// 		}
// 		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
//             if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
//                 return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//             }
//             return []byte(config.GetEnv("JWT_SECRET", "")), nil
//         })
// 		if err!= nil ||!token.Valid {
// 			log.Println("TokenValid error: ", err)
//             Forbidden(w)
//             return
//         }
// 		claims, ok := token.Claims.(*jwt.MapClaims)
// 		if !ok {
// 			log.Println("Claims error: ", err)
// 			Forbidden(w)
// 			return
// 		}
// 		userIDStr, ok := (*claims)["UserID"].(string)
// 		if!ok {
// 			log.Println("UserId error: ", err)
//             Forbidden(w)
//             return
//         }
// 		userID, err := strconv.Atoi(userIDStr)
// 		if err!= nil {
// 			log.Println("Atoi Convert error: ", err)
//             utils.WriteError(w, http.StatusInternalServerError, "Don't Panic This is From Us!")
//             return
//         }
// 		var ID uint = uint(userID)
// 		// get rider by the user ID
// 		rider, err := riderStore.GetRiderByUserID(ID)
// 		if err != nil {
// 			Forbidden(w)
// 			return
// 		}
// 		if rider.UserID == 0 {
// 			Forbidden(w)
// 			return
// 		}
// 		// save Rider Id to the request context
// 		ctx := context.WithValue(r.Context(), UserKey, userID)
//         ctx = context.WithValue(ctx, RiderKey, rider.ID)
//         r = r.WithContext(ctx)
//         handlerFunc(w, r)
// 	}
// }


// func UserAuth(handlerFunc http.HandlerFunc, riderStore models.RiderRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString, err := utils.GetTokenFromRequest(r)
// 		if err!= nil {
//             http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
//             return
//         }
// 		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
//             if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
//                 return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//             }
//             return []byte(config.GetEnv("JWT_SECRET", "")), nil
//         })
// 		if err!= nil || !token.Valid {
// 			log.Println("This is sign token error",err)
//             Forbidden(w)
// 			return
//         }
// 		// get claims
// 		claims, ok := token.Claims.(*jwt.MapClaims)
// 		fmt.Println(claims)
// 		if !ok {
// 			log.Println("This token claims error", ok)
// 			Forbidden(w)
//             return
// 		}
// 		userIDStr, ok := (*claims)["UserID"].(string)
// 		if !ok {
// 			log.Println("this userIdstr extract error", ok)
// 			Forbidden(w)
// 			return
// 		}
// 		userID, err := strconv.Atoi(userIDStr)
// 		if err!= nil {
// 			log.Println(err)
//             Forbidden(w)
//             return
//         }
// 		var ID uint = uint(userID)
// 		// get rider by the user ID
// 		_, err = riderStore.GetRiderByUserID(ID)
// 		log.Println("This rider fetch error", err)
// 		if err == nil {
// 			Forbidden(w)
//             return
// 		}
// 		// save User Id to the request context
// 		ctx := context.WithValue(r.Context(), UserKey, userID)
//         handlerFunc(w, r.WithContext(ctx))
// 	}
// }

func Auth(userStore types.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from the request
		tokenString, err := utils.GetTokenFromRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GetEnv("JWT_SECRET", "")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Get the UserID from the claims
		userIDStr, ok := (*claims)["UserID"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Convert the UserID to an integer
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		log.Println(userID)

		// Verify the user exists in the database
		_, err = userStore.GetUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("UserID", userID)

		// Proceed to the next handler
		c.Next()
	}
}

func Forbidden(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, "Unauthorized",)
}


func GetUserIDFromContext(ctx gin.Context) int {
	userID, ok := ctx.Get("UserID")
    if !ok {
        return -1
    }
    return userID.(int)
}

func GetRiderIDFromContext(ctx gin.Context) int {
	riderID, ok := ctx.Get("RiderID")
    if!ok {
        return -1
    }
    return riderID.(int)
}

// func signToken(tokenString string) (*jwt.Token, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
//             return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//         }
//         return []byte(config.GetEnv("JWT_SECRET", "")), nil
//     })
// 	return token, err
// }