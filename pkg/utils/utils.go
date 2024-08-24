package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/Ayobami6/pickitup_v3/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
)

var rdb *redis.Client
var ctx = context.Background()

var Validate = validator.New()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_URL", "localhost:6379"),
        DB:       0,
	})
}


func GenerateAndCacheVerificationCode(email string) (int, error) {
	rand.NewSource(time.Now().UnixNano())

	randomNumber := rand.Intn(9000) + 1000

	numberStr := fmt.Sprintf("%d", randomNumber)

	err := rdb.Set(ctx, email, numberStr, 15*time.Minute).Err()
	if err != nil {
		return 0, err
	}


	return randomNumber, nil
}

func GetCachedVerificationCode(email string) (int, error) {
	val, err := rdb.Get(ctx, email).Result()
	if err != nil {
		return 0, err
	}

	var randomNumber int
	_, err = fmt.Sscanf(string(val), "%d", &randomNumber)
	if err != nil {
		return 0, err
	}

	return randomNumber, nil
}



func WriteJSON(w http.ResponseWriter, status int, status_msg, data any, others ...string) error {
	message := ""
	if len(others) > 0 {
        message = others[0]
    }
	res := map[string]interface{}{
		"status": status_msg,
        "data":  data,
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)

}


func Response(statusCode int, data map[string]any, message any) map[string]any {
	var status string
	switch {
	case statusCode >= 200 && statusCode <= 299:
		status = "success"
	case statusCode >= 300 && statusCode <= 399:
		status = "redirect"
	case statusCode == 400:
		status = "error"
	case statusCode == 404:
		status = "not found"
	case statusCode >= 405 && statusCode <= 499:
		status = "error"
	case statusCode == 401 || statusCode == 403:
		status = "unauthorized"
	case statusCode >= 500:
		status = "error"
		message = "This is from us!, please contact admin"
	default:
		status = "error"
		message = "This is from us!, please contact admin"
	}
	res := map[string]any{
        "status": status,
        "data":  data,
        "message": message,
		"status_code": statusCode,
    }
	return res

}


func WriteError(w http.ResponseWriter, status int, err... string) {
	var errMessage string
	if len(err) > 0 {
        errMessage = err[0]
    } else {
		errMessage = "Don't Panic This is From Us!"
	}
	log.Println(err)
    WriteJSON(w, status, "error", nil, errMessage)
}

func ThrowError(err error) error{
	log.Println(err)
	return err
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("request body is missing")
	}
	return json.NewDecoder(r.Body).Decode(payload)

}

func SendMail(recipient string, subject string, username string, message string) error {
	tmpl, err := os.ReadFile("templates/verification_template.html")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return err
	}
	
	t, err := template.New("email").Parse(string(tmpl))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}
	
    data := map[string]interface{}{
        "UserName": username,
        "Message":  message,
    }
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}
	
	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", "sainthaywon80@gmail.com")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)

	// Set the HTML body
	m.SetBody("text/html", body.String())
	smtpHost := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := 465
	smtpUser := config.GetEnv("SMTP_USER", "protected@gmail.com")
	smtpPass := config.GetEnv("SMTP_PWD", "protected")

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email and handle errors
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	// Success message
	fmt.Println("Email sent successfully!")

	return nil

}


func GetTokenFromRequest(c *gin.Context) (string, error) {
	tokenAuth := c.GetHeader("Authorization")
    tokenQuery := c.Query("token")

    if tokenAuth!= "" {
        return tokenAuth, nil
    }

    if tokenQuery!= "" {
        return tokenQuery, nil
    }

    return "", errors.New("token not found in request")
}