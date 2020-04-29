package melonade_client_go

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/devit-tel/goerror"
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
			return nil, errWithInput(ErrUnableParseInputPayload, workflowName, revision, transactionId, payload).WithCause(err)
		}

		data = jsonData
	}

	path, err := url.Parse(fmt.Sprintf("/v1/transaction/%s/%s?transactionId=%s", workflowName, revision, transactionId))
	if err != nil {
		return nil, errWithInput(ErrUnableParseUrlRequest, workflowName, revision, transactionId, payload).WithCause(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprint(c.processManagerEndpoint+path.EscapedPath()), bytes.NewBuffer(data))
	if err != nil {
		return nil, errWithInput(ErrUnableCreateRequestStartWorkflow, workflowName, revision, transactionId, payload).WithCause(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errWithInput(ErrUnableStartWorkflow, workflowName, revision, transactionId, payload).WithCause(err)
	}
	defer resp.Body.Close()

	workflowResp := &StartWorkflowResponse{}
	if err := json.NewDecoder(resp.Body).Decode(workflowResp); err != nil {
		if bodyBytes, errReadIO := ioutil.ReadAll(resp.Body); errReadIO != nil {
			return nil, ErrUnableParseOutputPayload.WithKeyValueInput(
				"responseBody", string(bodyBytes)).WithCause(err)
		}

		return nil, errWithInput(ErrUnableParseOutputPayload, workflowName, revision, transactionId, payload).WithCause(err)
	}

	return workflowResp, nil
}

func errWithInput(err goerror.Error, workflowName, revision, transactionId string, payload interface{}) goerror.Error {
	return err.WithKeyValueInput(
		"workflowName", workflowName,
		"revision", revision,
		"transactionId", transactionId,
		"payload", payload).WithCause(err)
}
