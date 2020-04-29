package melonade_client_go

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	json "github.com/json-iterator/go"
)

type Client struct {
	processManagerEndpoint string
	httpClient             *http.Client
}

func New(processManagerEndpoint string) Service {
	transport := &http.Transport{}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		processManagerEndpoint: processManagerEndpoint,
		httpClient: &http.Client{
			Transport: transport,
		},
	}
}

func (c *Client) StartWorkflow(workflowName, revision, transactionId string, payload interface{}) (*StartWorkflowResponse, error) {
	var data []byte
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		data = jsonData
	}

	path, err := url.Parse(fmt.Sprintf("/v1/transaction/%s/%s?transactionId=%s", workflowName, revision, transactionId))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprint(c.processManagerEndpoint+path.EscapedPath()), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyBytes))
	workflowResp := &StartWorkflowResponse{}
	if err := json.NewDecoder(resp.Body).Decode(workflowResp); err != nil {
		fmt.Println("FAILED DECODE")
		return nil, err
	}

	return workflowResp, nil
}
