package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type MpesaService struct {
	ShortCode   string
	Passkey     string
	CallbackURL string
	Sandbox     bool

	// token caching
	accessToken    string
	accessTokenExp time.Time
	mu             sync.Mutex
}

var consumerKey = os.Getenv("CONSUMER_KEY")
var consumerSecret = os.Getenv("CONSUMER_SECRET")

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

type STKPushRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            string `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

type STKPushResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}
	if strings.HasPrefix(phone, "07") {
		phone = "254" + phone[1:]
	}
	return phone
}

func (m *MpesaService) getAccessToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.accessToken != "" && time.Now().Before(m.accessTokenExp) {
		return m.accessToken, nil
	}

	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	if !m.Sandbox {
		url = "https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(consumerKey, consumerSecret)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Println("==> Access Token Response:", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mpesa auth failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var tokenRes AccessTokenResponse
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return "", fmt.Errorf("invalid token response: %s", string(body))
	}

	if tokenRes.AccessToken == "" {
		return "", fmt.Errorf("empty access token received")
	}

	expSec, _ := time.ParseDuration(tokenRes.ExpiresIn + "s")
	m.accessToken = tokenRes.AccessToken
	m.accessTokenExp = time.Now().Add(expSec - 10*time.Second)

	return tokenRes.AccessToken, nil
}

func (m *MpesaService) STKPush(phone string, amount string) (*STKPushResponse, error) {
	phone = normalizePhone(phone)

	token, err := m.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	passkey := os.Getenv("PASSKEY")

	m.ShortCode = "174379"
	m.Passkey = passkey
	timestamp := time.Now().Format("20060102150405")
	password := base64.StdEncoding.EncodeToString([]byte(m.ShortCode + m.Passkey + timestamp))

	reqBody := STKPushRequest{
		BusinessShortCode: m.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            amount,
		PartyA:            phone,
		PartyB:            "254708374149",
		PhoneNumber:       phone,
		CallBackURL:       "https://elritch-xerically-wilfredo.ngrok-free.dev/payments/mpesa-callback",
		AccountReference:  "FASTABIZ SERVICES",
		TransactionDesc:   "Payment from app",
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	log.Println("==> STK Push Request Body:", string(data))

	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	if !m.Sandbox {
		url = "https://api.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// LOG full response body
	log.Println("==> STK Push Response Body:", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("mpesa stk push failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return &STKPushResponse{
		ResponseCode:        "0",
		ResponseDescription: "STK Push request sent",
		CustomerMessage:     "Check your phone to complete payment",
	}, nil
}
