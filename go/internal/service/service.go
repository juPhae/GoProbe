package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"goprobe/internal/mqtt"
	"goprobe/internal/util"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GinStart() {
	r := gin.Default()
	// 处理OPTIONS请求
	r.Use(corsMiddleware())
	// 登录接口
	r.POST("/login", login)

	// 获取状态接口
	r.GET("/status", authMiddleware(), getStatus)
	r.POST("/device/start", authMiddleware(), startShellHandler)
	r.POST("/device/stop", authMiddleware(), stopShellHandler)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 解密用户名
	username, err := util.Decrypt([]byte(user.Username), []byte("1234567812345678"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt username"})
		return
	}

	// 解密密码
	password, err := util.Decrypt([]byte(user.Password), []byte("1234567812345678"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt password"})
		return
	}
	log.Println("password ", string(password))
	log.Println("username  ", string(username))
	// 检查用户名和密码是否正确
	if string(username) == "phae" && string(password) == "123456" {
		// 创建token
		tokenString, err := generateToken(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

func getStatus(c *gin.Context) {
	// 遍历队列，取出最新的12个消息并解析为结构体对象
	var data []util.SystemStatus
	for e := mqtt.Queue.Back(); e != nil && len(data) < 12; e = e.Prev() {
		var status util.SystemStatus
		log.Println(string(e.Value.([]byte)))

		err := json.Unmarshal(e.Value.([]byte), &status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal message", "data": err})
			return
		}
		data = append(data, status)
	}
	// 将结构体切片返回给前端
	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": data})
}
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username, ok := claims["username"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}

			c.Set("username", username)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func startShellHandler(c *gin.Context) {
	// 从请求中获取需要的数据
	log.Println("startShellHandler")
	var deviceParams = util.DeviceParams{}
	if err := c.BindJSON(&deviceParams); err != nil {
		// 返回错误响应
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 组装MQTT消息
	topic := "device/status/topic/management"
	//message := fmt.Sprintf(`{"cpu_id":"%s","device":"%s",""}`, deviceParams.CpuID, deviceParams.Device)
	message := fmt.Sprintf(`{"cpu_id":"%s","device":"%s","command":"%s"}`, deviceParams.CpuID, deviceParams.Device, "start")
	// 通过mqtt发送消息
	mqtt.Publish(topic, message)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"result": 1})
}

func stopShellHandler(c *gin.Context) {
	// 从请求中获取需要的数据
	log.Println("stopShellHandler")
	var deviceParams = util.DeviceParams{}
	if err := c.BindJSON(&deviceParams); err != nil {
		// 返回错误响应
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 组装MQTT消息
	topic := "device/status/topic/management"
	//message := fmt.Sprintf(`{"cpu_id":"%s","device":"%s",""}`, deviceParams.CpuID, deviceParams.Device)
	message := fmt.Sprintf(`{"cpu_id":"%s","device":"%s","command":"%s"}`, deviceParams.CpuID, deviceParams.Device, "stop")
	// 通过mqtt发送消息
	mqtt.Publish(topic, message)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"result": 1})
}
