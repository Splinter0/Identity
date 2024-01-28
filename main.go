package main

import (
	"log"

	"github.com/Splinter0/identity/bankid"
	"github.com/Splinter0/identity/provider"
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
			provider.AGE,
		},
	}
	p.Authenticate(provider.BankIDAuthenticationRequest{
		UserIp:         "213.102.85.9",
		MessageForUser: "Welcome to our service!",
		SameDevice:     true,
	})
}
