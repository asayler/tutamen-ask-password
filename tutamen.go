package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	COLLECTION        = "fa6f8b4b-1877-4f48-8790-a952d5caf8a2"
	SECRET            = "e30bf9ba-4633-4af4-84a0-5df25dd8c6e2"
	TLS_KEY_PATH      = "/tut.key"
	TLS_CRT_PATH      = "/tut.crt"
	AUTH_URI          = "https://ac.tutamen-test.bdr1.volaticus.net/api/v1/authorizations/"
	SECRET_URI        = "https://ss.tutamen-test.bdr1.volaticus.net/api/v1/collections/" + COLLECTION + "/secrets/" + SECRET + "/versions/latest/"
	TOKEN_APPROVED    = "approved"
	TOKEN_PENDING     = "pending"
	TOKEN_DENIED      = "denied"
)

type AuthorizationRequest struct {
	Objperm string `json:"objperm"`
	Objtype string `json:"objtype"`
	Objuid  string `json:"objuid"`
}

type AuthorizationReply struct {
	Authorizations []string `json:"authorizations"`
}

type TokenReply struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type SecretReply struct {
	Data string `json:"data"`
}

func TutamenGetPassword() (string, error) {

	var auth    string
	var token   string
	var secret  string
	var err     error

	fmt.Println("Using hardcoded SSL certificate:  ", TLS_CRT_PATH)
	fmt.Println("Using hardcoded SSL key:          ", TLS_KEY_PATH)
	fmt.Println("Using hardcoded collection UUID:  ", COLLECTION)
	fmt.Println("Using hardcoded secret UUID:      ", SECRET)
	fmt.Println("Using hardcoded auth URI:         ", AUTH_URI)

	client, err := getClient()
	if err != nil {
		return "", errors.New("Error creating HTTP Client: " + err.Error())
	}

	auth, err = getAuthorization(client)
	if err != nil {
		return "", errors.New("Error getting authorization: " + err.Error())
	}
	fmt.Println("Authorization:", auth)

	for {
		// TODO don't break on pending
		token, err = getToken(client, auth)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Token:", token)

	secret, err = getSecret(client, token)
	if err != nil {
		return "", errors.New("Error getting secret: " + err.Error())
	}

	return secret, nil
}

func getClient() (*http.Client, error) {

	x509, err := tls.LoadX509KeyPair(TLS_CRT_PATH, TLS_KEY_PATH)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport {
		TLSClientConfig: &tls.Config {
			Certificates: []tls.Certificate{ x509 },
		},
	}

	return &http.Client{Transport: tr}, nil
}

func getAuthorization(client *http.Client) (string, error) {

	// Get

	body, err := json.Marshal(AuthorizationRequest{
		Objperm: "col-read",
		Objtype: "collection",
		Objuid:  COLLECTION,
	})

	if err != nil {
		return "", err
	}

	resp, err := client.Post(AUTH_URI, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	rbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	fmt.Println("Response Body:", string(rbody))

	// Parse

	var m AuthorizationReply

	err = json.Unmarshal(rbody, &m)
	if err != nil {
		return "", err
	}

	return m.Authorizations[0], nil
}

func getToken(client *http.Client, authorization string) (string, error) {

	// Get

	url := AUTH_URI + authorization
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	rbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(rbody))

	// Parse

	var m TokenReply

	err = json.Unmarshal(rbody, &m)
	if err != nil {
		return "", err
	}

	switch m.Status {
	case TOKEN_APPROVED:
		return m.Token, nil
	case TOKEN_PENDING:
		return "", errors.New("Token Pending")
	case TOKEN_DENIED:
		return "", errors.New("Token Denied")
	default:
		return "", errors.New("Uknown Token State: " + m.Status)
	}
}

func getSecret(client *http.Client, token string) (string, error) {

	req, err := http.NewRequest("GET", SECRET_URI, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("tutamen-tokens", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	rbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	fmt.Println("Response Body:", string(rbody))

	var m SecretReply

	err = json.Unmarshal(rbody, &m)
	if err != nil {
		return "", err
	}

	return m.Data, nil
}
