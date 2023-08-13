package sid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	Client = &http.Client{}
)

func CheckSidStatus(Ngrok string, Sid string) (string, error) {
	var Response CheckSidResponse

	req, err := http.NewRequest("POST", Ngrok+"/sid-status", bytes.NewBuffer([]byte(fmt.Sprintf(`{"Sid": "%v"}`, Sid))))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "", fmt.Errorf("CheckSidStatus: %v", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("CheckSidStatus: %v", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&Response); err != nil {
		return "", fmt.Errorf("CheckSidStatus: %v", err)
	}
	return Response.Status, nil
}

func GetSid(TargetNumber string, TargetName string, Service string, Ngrok string, LenDigits int, Callback bool) (string, error) {
	var Response GetSidresponse
	var Url string

	if Callback {
		Url = Ngrok + "/call-back"
	} else {
		Url = Ngrok + "/initiate-call"
	}

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer([]byte(fmt.Sprintf(`{"RecipientNumber": "%v", "TargetName": "%v", "Service": "%v", "LenDigits": %v}`, TargetNumber, TargetName, Service, LenDigits))))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "", fmt.Errorf("GetSid: %v", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GetSid: %v", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&Response); err != nil {
		return "", fmt.Errorf("GetSid: %v", err)
	}

	return Response.Sid, nil
}
