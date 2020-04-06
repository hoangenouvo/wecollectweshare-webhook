package dialogflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Request is the top-level struct holding all the information
// Basically links a response ID with a query result.
type Request struct {
	Session                     string                      `json:"session,omitempty"`
	ResponseID                  string                      `json:"responseId,omitempty"`
	QueryResult                 QueryResult                 `json:"queryResult,omitempty"`
	OriginalDetectIntentRequest OriginalDetectIntentRequest `json:"originalDetectIntentRequest,omitempty"`
}

// GetParams simply unmarshals the parameters to the given struct and returns
// an error if it's not possible
func (rw *Request) GetParams(i interface{}) interface{} {
	return rw.QueryResult.Parameters
}

// GetContext allows to search in the output contexts of the query
func (rw *Request) GetContext(ctx string, i interface{}) error {
	for _, c := range rw.QueryResult.OutputContexts {
		if strings.HasSuffix(c.Name, ctx) {
			return json.Unmarshal(c.Parameters, &i)
		}
	}
	return errors.New("context not found")
}

// NewContext is a helper function to create a new named context with params
// name and a lifespan
func (rw *Request) NewContext(name string, lifespan int, params interface{}) (*Context, error) {
	var err error
	var b []byte

	if b, err = json.Marshal(params); err != nil {
		return nil, err
	}
	ctx := &Context{
		Name:          fmt.Sprintf("%s/contexts/%s", rw.Session, name),
		LifespanCount: lifespan,
		Parameters:    b,
	}
	return ctx, nil
}

// QueryResult is the dataset sent back by DialogFlow
type QueryResult struct {
	QueryText                 string                 `json:"queryText,omitempty"`
	Action                    string                 `json:"action,omitempty"`
	LanguageCode              string                 `json:"languageCode,omitempty"`
	AllRequiredParamsPresent  bool                   `json:"allRequiredParamsPresent,omitempty"`
	IntentDetectionConfidence float64                `json:"intentDetectionConfidence,omitempty"`
	Parameters                map[string]interface{} `json:"parameters,omitempty"`
	OutputContexts            []*Context             `json:"outputContexts,omitempty"`
	Intent                    Intent                 `json:"intent,omitempty"`
}

// Intent describes the matched intent
type Intent struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// DialogFlowResponseData struct
type DialogFlowResponseData struct {
	Google DialogFlowResponseGoogle `json:"google"`
}

// DialogFlowResponseGoogle struct
type DialogFlowResponseGoogle struct {
	ExpectUserResponse bool                           `json:"expectUserResponse"`
	IsSsml             bool                           `json:"isSsml"`
	SystemIntent       DialogFlowResponseSystemIntent `json:"systemIntent"`
}

// DialogFlowResponseSystemIntent struct
type DialogFlowResponseSystemIntent struct {
	Intent string                             `json:"intent"`
	Data   DialogFlowResponseSystemIntentData `json:"data"`
}

// DialogFlowResponseSystemIntentData struct
type DialogFlowResponseSystemIntentData struct {
	Type        string   `json:"@type"`
	OptContext  string   `json:"optContext"`
	Permissions []string `json:"permissions"`
}
