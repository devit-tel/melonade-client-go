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
)
