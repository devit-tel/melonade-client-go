package melonade_client_go

import (
	"bytes"
	"context"
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

func (c *Client) StartWorkflow(ctx context.Context, workflowName, revision, transactionId string, payload interface{}) (*StartWorkflowResponse, goerror.Error) {
	var data []byte
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, errWithInput(ErrUnableParseInputPayload, workflowName, revision, transactionId, payload).WithCause(err)
		}

		data = jsonData
	}

	path, err := url.Parse(fmt.Sprintf("%s/v1/transaction/%s/%s", c.processManagerEndpoint, workflowName, revision))
	if err != nil {
		return nil, errWithInput(ErrUnableParseUrlRequest, workflowName, revision, transactionId, payload).WithCause(err)
	}

	params := url.Values{}
	params.Add("transactionId", transactionId)
	path.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path.String(), bytes.NewBuffer(data))
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

func (c *Client) GetWorkflowDefinitions() ([]*WorkflowDefinition, goerror.Error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/v1/definition/workflow", c.processManagerEndpoint))
	if err != nil {
		return nil, ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrReadResponseBodyFailed.WithCause(err)
	}

	var r workflowDefinitionsResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, ErrParseResponseBodyFailed.WithCause(err)
	}

	return r.Data, nil
}

func (c *Client) GetTaskDefinitions() ([]*TaskDefinition, goerror.Error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/v1/definition/task", c.processManagerEndpoint))
	if err != nil {
		return nil, ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrReadResponseBodyFailed.WithCause(err)
	}

	var r taskDefinitionsResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, ErrParseResponseBodyFailed.WithCause(err)
	}

	return r.Data, nil
}

func (c *Client) CreateTaskDefinition(t TaskDefinition) goerror.Error {
	return c.SetTaskDefinition(t)
}

func (c *Client) SetTaskDefinition(t TaskDefinition) goerror.Error {
	b, err := json.Marshal(t)
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/v1/definition/task", c.processManagerEndpoint), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadResponseBodyFailed.WithCause(err)
	}

	var r melonadeResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return ErrParseResponseBodyFailed.WithCause(err)
	}

	if r.Success != true {
		return ErrRequestNotSuccess.WithExtendMsg(fmt.Sprintf("%+v", r.Error))
	}

	return nil
}

func (c *Client) UpdateTaskDefinition(t TaskDefinition) goerror.Error {
	b, err := json.Marshal(t)
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}

	u, err := url.Parse(fmt.Sprintf("%s/v1/definition/task", c.processManagerEndpoint))
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}
	h := http.Header{}
	h.Set("Content-type", "application/json")

	resp, err := c.httpClient.Do(&http.Request{
		Method: http.MethodPut,
		URL:    u,
		Body:   ioutil.NopCloser(bytes.NewReader(b)),
		Header: h,
	})
	if err != nil {
		return ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadResponseBodyFailed.WithCause(err)
	}

	var r melonadeResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return ErrParseResponseBodyFailed.WithCause(err)
	}

	if r.Success != true {
		return ErrRequestNotSuccess.WithExtendMsg(fmt.Sprintf("%+v", r.Error))
	}

	return nil
}

func (c *Client) CreateWorkflowDefinition(w WorkflowDefinition) goerror.Error {
	return c.SetWorkflowDefinition(w)
}

func (c *Client) SetWorkflowDefinition(t WorkflowDefinition) goerror.Error {
	b, err := json.Marshal(t)
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/v1/definition/workflow", c.processManagerEndpoint), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadResponseBodyFailed.WithCause(err)
	}

	var r melonadeResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return ErrParseResponseBodyFailed.WithCause(err)
	}

	if r.Success != true {
		return ErrRequestNotSuccess.WithExtendMsg(fmt.Sprintf("%+v", r.Error))
	}

	return nil
}

func (c *Client) UpdateWorkflowDefinition(w WorkflowDefinition) goerror.Error {
	b, err := json.Marshal(w)
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}

	u, err := url.Parse(fmt.Sprintf("%s/v1/definition/workflow", c.processManagerEndpoint))
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}
	h := http.Header{}
	h.Set("Content-type", "application/json")

	resp, err := c.httpClient.Do(&http.Request{
		Method: http.MethodPut,
		URL:    u,
		Body:   ioutil.NopCloser(bytes.NewReader(b)),
		Header: h,
	})
	if err != nil {
		return ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadResponseBodyFailed.WithCause(err)
	}

	var r melonadeResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return ErrParseResponseBodyFailed.WithCause(err)
	}

	if r.Success != true {
		return ErrRequestNotSuccess.WithExtendMsg(fmt.Sprintf("%+v", r.Error))
	}

	return nil
}

func (c *Client) DeleteWorkflowDefinition(name, rev string) goerror.Error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/definition/workflow/%s/%s", c.processManagerEndpoint, name, rev))
	if err != nil {
		return ErrParseRequestBodyFailed.WithCause(err)
	}
	h := http.Header{}
	h.Set("Content-type", "application/json")

	resp, err := c.httpClient.Do(&http.Request{
		Method: http.MethodDelete,
		URL:    u,
	})
	if err != nil {
		return ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadResponseBodyFailed.WithCause(err)
	}

	var r melonadeResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return ErrParseResponseBodyFailed.WithCause(err)
	}

	if r.Success != true {
		return ErrRequestNotSuccess.WithExtendMsg(fmt.Sprintf("%+v", r.Error))
	}

	return nil
}

func (c *Client) GetTransactionData(tID string) (*Transaction, goerror.Error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/v1/transaction/%s", c.processManagerEndpoint, tID))
	if err != nil {
		return nil, ErrRequestFailed.WithCause(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrReadResponseBodyFailed.WithCause(err)
	}

	var r transactionResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, ErrParseResponseBodyFailed.WithCause(err)
	}

	if r.Success != true {
		return nil, ErrRequestNotSuccess.WithExtendMsg(fmt.Sprintf("%+v", r.Error))
	}

	if r.Data == nil {
		return nil, ErrTransactionNotFound
	}

	return r.Data, nil
}

func errWithInput(err goerror.Error, workflowName, revision, transactionId string, payload interface{}) goerror.Error {
	return err.WithKeyValueInput(
		"workflowName", workflowName,
		"revision", revision,
		"transactionId", transactionId,
		"payload", payload).WithCause(err)
}

type workflowDefinitionsResponse struct {
	Data []*WorkflowDefinition `json:"data"`
	melonadeResponse
}

type taskDefinitionsResponse struct {
	Data []*TaskDefinition `json:"data"`
	melonadeResponse
}

type transactionResponse struct {
	Data *Transaction `json:"data"`
	melonadeResponse
}

type melonadeResponse struct {
	Error   map[string]interface{} `json:"error"`
	Success bool                   `json:"success"`
}
