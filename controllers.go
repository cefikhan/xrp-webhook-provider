package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

//JWT LOGIC
type RegisterUserRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Userpassword string `json:"userpassword" validate:"required"`
}

var jwtKey = []byte("secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"userpassword"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   int
	jwt.StandardClaims
}

// JWT

type InsertURLRequest struct {
	URL string `json:"url" validate:"required"`
}

type InsertAddressRequest struct {
	Address string `json:"address" validate:"required"`
}

type RegisterAddressRequest struct {
	Address  string `json:"address" validate:"required"`
	StreamID int    `json:"streamid" validate:"required"`
}

type CreateStreamRequest struct {
	Webhookurl string `json:"webhookurl" validate:"required"`
}

type UpdateStreamRequest struct {
	Webhookurl string `json:"webhookurl" validate:"required"`
	StreamID   int    `json:"streamid" validate:"required"`
}

type viewAddressesRequest struct {
	StreamID int `json:"streamid" validate:"required"`
}

type deleteAddressesRequest struct {
	StreamID  int `json:"streamid" validate:"required"`
	AddressID int `json:"addressid" validate:"required"`
}

//types
type ResponseData struct {
	Addresses []string `json:"addresses"`
	Url       string   `json:"url"`
}

type UserDetails struct {
	Username string `json:"username"`
}

func RegisterUserHandler(c *gin.Context) {
	var req RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Not registered"})
		return
	}

	res, err := GetUserDetailsQuery(db, req.Username)
	if res.Id != 0 {
		fmt.Println("yser idd", res.Id)
		c.JSON(http.StatusOK, gin.H{"data": "User already Exists"})
		return
	}

	err = RegisterUserQuery(db, req.Username, req.Email, req.Userpassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": "Error in Query"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"data": "registered"})
}

func SignInHandler(c *gin.Context) {

	var credentials Credentials
	err := json.NewDecoder(c.Request.Body).Decode(&credentials)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"data": "Not sign in"})

		return
	}

	res, err := GetUserDetailsQuery(db, credentials.Username)
	if err != nil {
		fmt.Println("Query Error")
	}

	if err != nil || res.Userpassword != credentials.Password {

		c.JSON(http.StatusUnauthorized, gin.H{"data": "Not Authorized"})

		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		UserID:   res.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"data": "Unable to produce jwt"})
		return
	}

	http.SetCookie(c.Writer,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	ud := UserDetails{
		Username: credentials.Username,
	}

	c.JSON(http.StatusOK, gin.H{"data": ud})

}

func CreateStreamHandler(c *gin.Context) {

	var req CreateStreamRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to bind request"})
		return
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {

			c.JSON(http.StatusUnauthorized, gin.H{"data": "No Cookie"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})

		return
	}

	tokenStr := cookie

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})

		return
	}

	if !tkn.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})
		return
	}

	err = CreateStreamQuery(db, claims.UserID, req.Webhookurl)

	if err != nil {
		fmt.Println("Query Error")
	}

	c.JSON(http.StatusOK, gin.H{"data": "Stream created successfully"})

}

func registerAddressToStreamHandler(c *gin.Context) {

	var req RegisterAddressRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to bind request"})
		return
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {

			c.JSON(http.StatusUnauthorized, gin.H{"data": "No Cookie"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})

		return
	}

	tokenStr := cookie

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})

		return
	}

	if !tkn.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})
		return
	}

	res, err := GetUserAllStreamQuery(db, claims.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return
	}

	if res == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return

	}

	streamX := 0

	for i := 0; i < len(res); i++ {
		if res[i].Id == req.StreamID {
			streamX = res[i].Id
		}
	}

	if streamX == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such StreamID"})
		return
	}

	err = registerAddressToStreamQuery(db, streamX, req.Address)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error in DB Query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Address Registered Successfully"})
	return

}

func viewAddressesHandler(c *gin.Context) {

	var req viewAddressesRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to bind request"})
		return
	}

	userID, err := tokenMiddleware(c)
	if err != nil {
		return
	}
	res, err := GetUserAllStreamQuery(db, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return
	}

	if res == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return

	}

	streamX := 0

	for i := 0; i < len(res); i++ {
		if res[i].Id == req.StreamID {
			streamX = res[i].Id
		}
	}

	if streamX == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such StreamID"})
		return
	}

	resx, err := getAddresswithIDsfromStreamQuery(db, streamX)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "No address"})
		return
	}

	if len(resx) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"data": "No address registered"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"data": resx})
}

func viewStreamHandler(c *gin.Context) {

	var req CreateStreamRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to bind request"})
		return
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {

			c.JSON(http.StatusUnauthorized, gin.H{"data": "No Cookie"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})

		return
	}

	tokenStr := cookie

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})

		return
	}

	if !tkn.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})
		return
	}

	//
	res, err := GetUserAllStreamQuery(db, claims.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to get streamID for user"})
		return
	}

	if res == nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "No Stream Created"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"streams": res})

}

func updateStreamURLHandler(c *gin.Context) {

	var req UpdateStreamRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to bind request"})
		return
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {

			c.JSON(http.StatusUnauthorized, gin.H{"data": "No Cookie"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})

		return
	}

	tokenStr := cookie

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})

		return
	}

	if !tkn.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})
		return
	}

	fmt.Println("new webhook url is ", req.Webhookurl)

	res, err := GetUserAllStreamQuery(db, claims.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return
	}

	if res == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return

	}

	streamX := 0

	for i := 0; i < len(res); i++ {
		if res[i].Id == req.StreamID {
			streamX = res[i].Id
		}
	}

	if streamX == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such StreamID"})
		return
	}

	err = UpdateStreamQuery(db, streamX, req.Webhookurl)

	if err != nil {
		fmt.Println("Query Error")
	}

	x := fmt.Sprintf("Webhook url updated successfully for streamID %s", streamX)
	c.JSON(http.StatusOK, gin.H{"data": x})
}

func deleteAddressHandler(c *gin.Context) {

	var req deleteAddressesRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": "Unable to bind request"})
		return
	}

	userID, err := tokenMiddleware(c)

	if err != nil {
		return
	}

	res, err := GetUserAllStreamQuery(db, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return
	}

	if res == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "First create a Stream"})
		return

	}

	streamX := 0

	for i := 0; i < len(res); i++ {
		if res[i].Id == req.StreamID {
			streamX = res[i].Id
		}
	}

	if streamX == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such StreamID"})
		return
	}

	if err != nil {
		return
	}

	addressesswithIDS, err := getAddresswithIDsfromStreamQuery(db, streamX)

	addressID := 0
	for i := 0; i < len(addressesswithIDS); i++ {
		if addressesswithIDS[i].ID == req.AddressID {
			addressID = addressesswithIDS[i].ID
		}
	}

	if addressID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such AddressID exists for you"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "query error"})
	}

	fmt.Println(addressID)

	err = deleteAddressfromStreamQuery(db, streamX, addressID)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": "deleted query error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": addressID})
}

func tokenMiddleware(c *gin.Context) (int, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {

			c.JSON(http.StatusUnauthorized, gin.H{"data": "No Cookie"})
			return 0, err
		}
		c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})

		return 0, err
	}

	tokenStr := cookie

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"data": "StatusUnauthorized"})
			return 0, err
		}
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})

		return 0, err
	}

	if !tkn.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"data": "StatusUnauthorized"})
		return 0, err
	}

	return claims.UserID, nil
}
