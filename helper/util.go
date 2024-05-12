package helper

import (
	"be-assignment/model"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/JPratama7/util/convert"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
	"log"
	"strings"
	"time"
)

const PASSKEYLEN = 32
const SALTKEYLEN = 16

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Tangkap error yang muncul di middleware atau handler
		err := c.Errors.Last()
		if err != nil {
			// Tangani error di sini
			var e *model.Error
			switch {
			case errors.As(err.Err, &e):
				// Tangani custom error
				c.JSON(e.Code, gin.H{
					"error": e.Message,
				})
			default:
				// Tangani error lainnya
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
			}

			// Hentikan eksekusi konteks
			c.Abort()
		}
	}
}

func SetContextData[T any](key string, data *T) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, data)
		c.Next()
	}
}

func GetContextData[T any](c *gin.Context, key string) (*T, error) {
	data, ok := c.Get(key)
	if !ok {
		return nil, errors.New("data not found")
	}

	return data.(*T), nil
}

func HashPassword(password string) (hashed, salted string, err error) {
	salt, err := generateRandomBytes(SALTKEYLEN)
	if err != nil {
		return
	}
	hash := generateFromPassword(convert.UnsafeBytes(password), salt, PASSKEYLEN)

	hashed, salted = base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt)
	return
}

func VerifyPassword(password, hash, salt string) (res bool) {

	salted, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return
	}

	pass, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return
	}

	newHash := generateFromPassword(convert.UnsafeBytes(password), salted, PASSKEYLEN)

	if !bytes.Equal(newHash, pass) {
		return false
	}

	return true
}

func generateFromPassword(password, salt []byte, keylen uint32) []byte {
	return argon2.IDKey(password, salt, 3, 64*1024, 2, keylen)
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func ConvertToHandler(fn func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			log.Printf("Err: %+v\n", c.Error(err))
		}
	}
}

func GetAuthHeader(c *gin.Context) (token string, err error) {
	token = c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("token not found")
	}
	split := strings.Split(token, " ")
	if len(split) < 1 {
		return "", errors.New("invalid token")
	}
	token = split[1]
	return
}

func NewTimestamp() primitive.Timestamp {
	return primitive.Timestamp{T: uint32(time.Now().Unix())}
}
