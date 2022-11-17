package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

func getToken(directory_id string, client_id string, client_secret string) string {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", directory_id)
	method := "GET"
	payloadstr := fmt.Sprintf("client_id=%s&grant_type=client_credentials&client_secret=%s&resource=https://management.azure.com", client_id, client_secret)

	payload := strings.NewReader(payloadstr)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	type tokenResp struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    string `json:"expires_in"`
		ExtExpiresIn string `json:"ext_expires_in"`
		ExpiresOn    string `json:"expires_on"`
		NotBefore    string `json:"not_before"`
		Resource     string `json:"resource"`
		AccessToken  string `json:"access_token"`
	}
	var result tokenResp
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	return string(result.AccessToken)

}

func createStorageAccount(accsToken string, subscription_id string, account_name string) string {
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/cloudformationresources/providers/Microsoft.Storage/storageAccounts/%s?api-version=2018-02-01", subscription_id, account_name)
	method := "PUT"

	payload := strings.NewReader(`{` + "" + `"sku": {` + "" + `"name": "Standard_GRS"` + "" + `},` + "" + `"kind": "StorageV2",` + "" + `"location": "westus"` + "" + `}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	headerToken := fmt.Sprintf("Bearer %s", accsToken)
	req.Header.Add("Authorization", headerToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	fmt.Println(string(body))
	return "storageaccountsuccess"
}

func deleteStorageAccount(accsToken string, subscription_id string, account_name string) string {
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/cloudformationresources/providers/Microsoft.Storage/storageAccounts/%s?api-version=2021-09-01", subscription_id, account_name)
	method := "DELETE"

	payload := strings.NewReader(`{` + "" + `"sku": {` + "" + `"name": "Standard_GRS"` + "" + `},` + "" + `"kind": "StorageV2",` + "" + `"location": "westus"` + "" + `}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	headerToken := fmt.Sprintf("Bearer %s", accsToken)
	req.Header.Add("Authorization", headerToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	fmt.Println(string(body))
	return "storageaccountdeletesuccess"
}

func HandleRequest(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	fmt.Println(event)
	account_name, _ := event.ResourceProperties["StorageAccountName"].(string)
	req_type := event.RequestType
	fmt.Println(req_type)
	directory_id := os.Getenv("directory_id")
	client_id := os.Getenv("client_id")
	client_secret := os.Getenv("client_secret")
	subscription_id := os.Getenv("subscription_id")
	accessToken := getToken(directory_id, client_id, client_secret)
	if req_type == "Create" {
		createStatus := createStorageAccount(accessToken, subscription_id, account_name)
		fmt.Println("Creation done")
		fmt.Println(createStatus)
	} else if req_type == "Delete" {
		delStatus := deleteStorageAccount(accessToken, subscription_id, account_name)
		fmt.Println("deletion done")
		fmt.Println(delStatus)
	}

	return
}

func main() {
	lambda.Start(cfn.LambdaWrap(HandleRequest))
}
