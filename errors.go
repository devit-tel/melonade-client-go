package melonade_client_go

import (
	"github.com/devit-tel/goerror"
)

var (
	ErrUnableParseUrlRequest            = goerror.DefineBadRequest("UnableParseUrlRequest", "unable parse url request")
	ErrUnableParseInputPayload          = goerror.DefineBadRequest("UnableParseInputPayload", "unable parse input payload")
	ErrUnableParseOutputPayload         = goerror.DefineInternalServerError("UnableParseOutputPayload", "unable parse output payload")
	ErrUnableCreateRequestStartWorkflow = goerror.DefineInternalServerError("UnableCreateRequestStartWorkflow", "unable create request start workflow")
	ErrUnableStartWorkflow              = goerror.DefineInternalServerError("UnableStartWorkflow", "unable start workflow")
	ErrUnableParseKafkaVersion          = goerror.DefineInternalServerError("UnableParseKafkaVersion", "unable parse kafka version")
	ErrUnableToCreateProducer           = goerror.DefineInternalServerError("UnableToCreateProducer", "unable to create producer")
	ErrUnableToCreateConsumer           = goerror.DefineInternalServerError("UnableToCreateConsumer", "unable to create consumer")
	ErrRequestFailed                    = goerror.DefineInternalServerError("ErrRequestFailed", "request failed")
	ErrReadResponseBodyFailed           = goerror.DefineInternalServerError("ErrReadResponseBodyFailed", "read response body failed")
	ErrParseResponseBodyFailed          = goerror.DefineInternalServerError("ErrParseResponseBodyFailed", "parse response body failed")
	ErrParseRequestBodyFailed           = goerror.DefineInternalServerError("ErrParseRequestBodyFailed", "parse request body failed")
	ErrRequestNotSuccess                = goerror.DefineInternalServerError("ErrRequestNotSuccess", "request not success")
	ErrTransactionNotFound              = goerror.DefineInternalServerError("ErrTransactionNotFound", "transaction not found")
)
