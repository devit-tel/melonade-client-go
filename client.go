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

func (c *Client) StartWorkflowWithTags(ctx context.Context, workflowName, revision, transactionId string, payload interface{}, tags interface{}) (*StartWorkflowResponse, goerror.Error) {
	var data []byte
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, errWithInput(ErrUnableParseInputPayload, workflowName, revision, transactionId, payload).WithCause(err)
		}

		data = jsonData
	}

	tagsJsonString, err := json.MarshalToString(tags)
	if err != nil {
		if err != nil {
			return nil, errWithInputTag(ErrUnableParseInputPayload, workflowName, revision, transactionId, payload, tags).WithCause(err)
		}
	}

	path, err := url.Parse(fmt.Sprintf("%s/v1/transaction/%s/%s", c.processManagerEndpoint, workflowName, revision))
	if err != nil {
		return nil, errWithInput(ErrUnableParseUrlRequest, workflowName, revision, transactionId, payload).WithCause(err)
	}

	params := url.Values{}
	params.Add("transactionId", transactionId)
	params.Add("tags", tagsJsonString)
	path.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodPost, path.String(), bytes.NewBuffer(data))
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

func (c *Client) StartWorkflow(ctx context.Context, workflowName, revision, transactionId string, payload interface{}) (*StartWorkflowResponse, goerror.Error) {
	return c.StartWorkflowWithTags(ctx, workflowName, revision, transactionId, payload, nil)
}

func errWithInput(err goerror.Error, workflowName, revision, transactionId string, payload interface{}) goerror.Error {
	return err.WithKeyValueInput(
		"workflowName", workflowName,
		"revision", revision,
		"transactionId", transactionId,
		"payload", payload).WithCause(err)
}

func errWithInputTag(err goerror.Error, workflowName, revision, transactionId string, payload interface{}, tags interface{}) goerror.Error {
	return err.WithKeyValueInput(
		"workflowName", workflowName,
		"revision", revision,
		"transactionId", transactionId,
		"payload", payload,
		"tags", tags).WithCause(err)
}
