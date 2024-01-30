package main

import (
	"log"

	"github.com/Splinter0/identity/bankid"
	"github.com/Splinter0/identity/provider"
	"github.com/gin-gonic/gin"
)

func main() {
	/*r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		c.HTML(200, "phish.html", gin.H{})
	})
	r.Run()*/
	rp, err := bankid.NewBankIDRP(bankid.TEST)
	if err != nil {
		log.Fatal(err)
	}

	p := &provider.BankIDProvider{
		Client: rp,
		DefaultUserDetails: []provider.UserDetailType{
			provider.FIRST_NAME,
			provider.LAST_NAME,
		},
		LaunchURLChannel: make(chan string, 1),
		QRCodeChannel:    make(chan []byte, 30),
		MessageChannel:   make(chan string),
		ResponseChannel:  make(chan provider.BankIDAuthenticationResponse, 1),
		Config: provider.BankIDConfig{
			CompanyName:     "Jesus AB",
			RedirectBaseUrl: "http://localhost",
		},
	}

	r := gin.Default()
	r.POST("/same", func(c *gin.Context) {
		go p.Authenticate(provider.BankIDAuthenticationRequest{
			RequestedDetails: []provider.UserDetailType{},
			SameDevice:       true,
			UserIp:           c.ClientIP(),
			MessageForUser:   "Log into blah",
		})

		c.Redirect(302, <-p.LaunchURLChannel)
	})
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, struct {
			Message string `json:"message"`
		}{Message: <-p.MessageChannel})
	})
	r.Run("0.0.0.0:8080")
}
