package main

import (
	"net/http"
	"strconv"

	"github.com/Splinter0/identity/bankid"
	"github.com/gin-gonic/gin"
)

const CSRF_HEADER = "X-BankID-CSRF"

// Uses custom header csrf protection
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#employing-custom-request-headers-for-ajaxapi
func CsrfProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only using post to change data
		if c.Request.Method == "POST" && c.Request.Header.Get(CSRF_HEADER) == "" {
			c.JSON(401, gin.H{"message": "CSRF check failed"})
			c.Abort()
		}
	}
}

func main() {
	p := bankid.NewBankIDProvider()

	r := gin.Default()
	r.Use(CsrfProtection())
	r.LoadHTMLGlob("templates/bankid/*")
	r.Static("/js/", "static/js/")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.POST("/start", func(c *gin.Context) {
		// If had a transaction ongoing, cancel it
		transactionKey, err := c.Cookie("bankidTransaction")
		if err == nil {
			p.Cancel(transactionKey)
		}
		sameValue, ok := c.GetQuery("same")
		var sameDevice bool
		if !ok {
			sameDevice = true
		} else {
			var err error
			sameDevice, err = strconv.ParseBool(sameValue)
			if err != nil {
				sameDevice = true
			}
		}
		authResponse := p.Authenticate(bankid.BankIDAuthenticationRequest{
			SameDevice:     sameDevice,
			UserIp:         c.ClientIP(),
			MessageForUser: "Log into blah",
			RedirectURL:    "null",
			Mobile:         bankid.IsMobileUserAgent(c.Request.UserAgent()),
		})
		var code int
		if !authResponse.Success {
			code = 401
		} else {
			c.SetSameSite(http.SameSiteStrictMode)
			c.SetCookie(
				"bankidTransaction",
				authResponse.TransactionKey,
				30,
				"/",
				"dev.mastersplinter.work",
				false,
				true,
			)
			code = 200
		}
		c.JSON(code, authResponse)
	})
	r.GET("/status", func(c *gin.Context) {
		transactionKey, err := c.Cookie("bankidTransaction")
		if err != nil {
			c.JSON(400, bankid.BankIDStatusResponse{
				Message: "Transaction expired",
				Status:  bankid.FAILED,
			})
			return
		}
		statusResponse := p.Status(transactionKey)
		var code int
		if statusResponse.Status == bankid.FAILED {
			code = 401
		} else {
			code = 200
		}
		c.JSON(code, statusResponse)
	})
	r.POST("/cancel", func(c *gin.Context) {
		transactionKey, err := c.Cookie("bankidTransaction")
		if err != nil {
			c.JSON(400, bankid.BankIDStatusResponse{
				Message: "Transaction not started",
				Status:  bankid.FAILED,
			})
			return
		}
		p.Cancel(transactionKey)
		c.JSON(204, nil)
	})
	r.Run("0.0.0.0:8080")
}
