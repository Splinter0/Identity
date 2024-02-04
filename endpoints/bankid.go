package endpoints

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Splinter0/identity/bankid"
	"github.com/gin-gonic/gin"
)

const CSRF_HEADER = "X-BankID-CSRF"

// Uses custom header csrf protection
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#employing-custom-request-headers-for-ajaxapi
func csrfProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only using post to change data
		if c.Request.Method == "POST" && c.Request.Header.Get(CSRF_HEADER) == "" {
			c.JSON(401, gin.H{"message": "CSRF check failed"})
			c.Abort()
		}
	}
}

func RegisterBankIDEndpoints(r *gin.Engine, config *Config) {
	if config.BankID.Domain == nil {
		log.Fatal("Cannot register BankID provider without a 'domain' in config.yml")
	}
	if config.BankID.Env == bankid.TEST {
		log.Println("BankID is configured for testing, not to use in production")
	}
	p := bankid.NewBankIDProvider(config.BankID)
	r.Use(csrfProtection())
	r.GET("/bankid", func(c *gin.Context) {
		c.HTML(200, "bankid.html", gin.H{
			"Service": *config.Service,
		})
	})
	r.POST("/bankid/start", func(c *gin.Context) {
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
			MessageForUser: config.BankID.VisibleMessage,
			RedirectURL:    "null",
			UserAgent:      c.Request.UserAgent(),
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
				"/bankid",
				*config.BankID.Domain,
				config.BankID.Env == bankid.PRODUCTION,
				true,
			)
			code = 200
		}
		c.JSON(code, authResponse)
	})
	r.GET("/bankid/status", func(c *gin.Context) {
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
	r.POST("/bankid/cancel", func(c *gin.Context) {
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
}
