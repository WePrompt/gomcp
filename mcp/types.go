package mcp

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	JSONRPCVersion = "2.0"
)

// Standard JSON-RPC error codes
const (
	ErrorCodeParseError     = -32700
	ErrorCodeInvalidRequest = -32600
	ErrorCodeMethodNotFound = -32601
	ErrorCodeInvalidParams  = -32602
	ErrorCodeInternalError  = -32603
)

// Base for objects that include optional annotations for the client. The client
// can use annotations to inform how objects are used or displayed
type Annotated struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *AnnotatedAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`
}

type AnnotatedAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AnnotatedAnnotations) UnmarshalJSON(b []byte) error {
	type Plain AnnotatedAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = AnnotatedAnnotations(plain)
	return nil
}

type BlobResourceContents struct {
	// A base64-encoded string representing the binary data of the item.
	Blob string `json:"blob" yaml:"blob" mapstructure:"blob"`

	// The MIME type of this resource, if known.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// The URI of this resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BlobResourceContents) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["blob"]; raw != nil && !ok {
		return fmt.Errorf("field blob in BlobResourceContents: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in BlobResourceContents: required")
	}
	type Plain BlobResourceContents
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = BlobResourceContents(plain)
	return nil
}

// Used by the client to invoke a tool provided by the server.
type CallToolRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CallToolRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CallToolRequestParams struct {
	// Arguments corresponds to the JSON schema field "arguments".
	Arguments CallToolRequestParamsArguments `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

type CallToolRequestParamsArguments map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CallToolRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in CallToolRequestParams: required")
	}
	type Plain CallToolRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CallToolRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CallToolRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CallToolRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CallToolRequest: required")
	}
	type Plain CallToolRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CallToolRequest(plain)
	return nil
}

// The server's response to a tool call.
//
// Any errors that originate from the tool SHOULD be reported inside the result
// object, with `isError` set to true, _not_ as an MCP protocol-level error
// response. Otherwise, the LLM would not be able to see that an error occurred
// and self-correct.
//
// However, any errors in _finding_ the tool, an error indicating that the
// server does not support tool calls, or any other exceptional conditions,
// should be reported as an MCP error response.
type CallToolResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta CallToolResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Content corresponds to the JSON schema field "content".
	Content []interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// Whether the tool call ended in an error.
	//
	// If not set, this is assumed to be false (the call was successful).
	IsError bool `json:"isError,omitempty" yaml:"isError,omitempty" mapstructure:"isError,omitempty"`
}

func (j *CallToolResult) AddTextContent(content TextContent) {
	j.Content = append(j.Content, content)
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type CallToolResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CallToolResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in CallToolResult: required")
	}
	type Plain CallToolResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CallToolResult(plain)
	return nil
}

// This notification can be sent by either side to indicate that it is cancelling a
// previously-issued request.
//
// The request SHOULD still be in-flight, but due to communication latency, it is
// always possible that this notification MAY arrive after the request has already
// finished.
//
// This notification indicates that the result will be unused, so any associated
// processing SHOULD cease.
//
// A client MUST NOT attempt to cancel its `initialize` request.
type CancelledNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CancelledNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CancelledNotificationParams struct {
	// An optional string describing the reason for the cancellation. This MAY be
	// logged or presented to the user.
	Reason *string `json:"reason,omitempty" yaml:"reason,omitempty" mapstructure:"reason,omitempty"`

	// The ID of the request to cancel.
	//
	// This MUST correspond to the ID of a request previously issued in the same
	// direction.
	RequestId interface{} `json:"requestId" yaml:"requestId" mapstructure:"requestId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CancelledNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["requestId"]; raw != nil && !ok {
		return fmt.Errorf("field requestId in CancelledNotificationParams: required")
	}
	type Plain CancelledNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CancelledNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CancelledNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CancelledNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CancelledNotification: required")
	}
	type Plain CancelledNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CancelledNotification(plain)
	return nil
}

// Capabilities a client may support. Known capabilities are defined here, in this
// schema, but this is not a closed set: any client can define its own, additional
// capabilities.
type ClientCapabilities struct {
	// Experimental, non-standard capabilities that the client supports.
	Experimental ClientCapabilitiesExperimental `json:"experimental,omitempty" yaml:"experimental,omitempty" mapstructure:"experimental,omitempty"`

	// Present if the client supports listing roots.
	Roots *ClientCapabilitiesRoots `json:"roots,omitempty" yaml:"roots,omitempty" mapstructure:"roots,omitempty"`

	// Present if the client supports sampling from an LLM.
	Sampling ClientCapabilitiesSampling `json:"sampling,omitempty" yaml:"sampling,omitempty" mapstructure:"sampling,omitempty"`
}

// Experimental, non-standard capabilities that the client supports.
type ClientCapabilitiesExperimental map[string]map[string]interface{}

// Present if the client supports listing roots.
type ClientCapabilitiesRoots struct {
	// Whether the client supports notifications for changes to the roots list.
	ListChanged bool `json:"listChanged,omitempty" yaml:"listChanged,omitempty" mapstructure:"listChanged,omitempty"`
}

// Present if the client supports sampling from an LLM.
type ClientCapabilitiesSampling map[string]interface{}

type ClientNotification interface{}

type ClientRequest interface{}

type ClientResult interface{}

// A request from the client to the server, to ask for completion options.
type CompleteRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CompleteRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CompleteRequestParams struct {
	// The argument's information
	Argument CompleteRequestParamsArgument `json:"argument" yaml:"argument" mapstructure:"argument"`

	// Ref corresponds to the JSON schema field "ref".
	Ref interface{} `json:"ref" yaml:"ref" mapstructure:"ref"`
}

// The argument's information
type CompleteRequestParamsArgument struct {
	// The name of the argument
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The value of the argument to use for completion matching.
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteRequestParamsArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in CompleteRequestParamsArgument: required")
	}
	if _, ok := raw["value"]; raw != nil && !ok {
		return fmt.Errorf("field value in CompleteRequestParamsArgument: required")
	}
	type Plain CompleteRequestParamsArgument
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteRequestParamsArgument(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["argument"]; raw != nil && !ok {
		return fmt.Errorf("field argument in CompleteRequestParams: required")
	}
	if _, ok := raw["ref"]; raw != nil && !ok {
		return fmt.Errorf("field ref in CompleteRequestParams: required")
	}
	type Plain CompleteRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CompleteRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CompleteRequest: required")
	}
	type Plain CompleteRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteRequest(plain)
	return nil
}

// The server's response to a completion/complete request
type CompleteResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta CompleteResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Completion corresponds to the JSON schema field "completion".
	Completion CompleteResultCompletion `json:"completion" yaml:"completion" mapstructure:"completion"`
}

type CompleteResultCompletion struct {
	// Indicates whether there are additional completion options beyond those provided
	// in the current response, even if the exact total is unknown.
	HasMore bool `json:"hasMore,omitempty" yaml:"hasMore,omitempty" mapstructure:"hasMore,omitempty"`

	// The total number of completion options available. This can exceed the number of
	// values actually sent in the response.
	Total *int `json:"total,omitempty" yaml:"total,omitempty" mapstructure:"total,omitempty"`

	// An array of completion values. Must not exceed 100 items.
	Values []string `json:"values" yaml:"values" mapstructure:"values"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteResultCompletion) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["values"]; raw != nil && !ok {
		return fmt.Errorf("field values in CompleteResultCompletion: required")
	}
	type Plain CompleteResultCompletion
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteResultCompletion(plain)
	return nil
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type CompleteResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["completion"]; raw != nil && !ok {
		return fmt.Errorf("field completion in CompleteResult: required")
	}
	type Plain CompleteResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteResult(plain)
	return nil
}

// A request from the server to sample an LLM via the client. The client has full
// discretion over which model to select. The client should also inform the user
// before beginning sampling, to allow them to inspect the request (human in the
// loop) and decide whether to approve it.
type CreateMessageRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CreateMessageRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CreateMessageRequestParams struct {
	// A request to include context from one or more MCP servers (including the
	// caller), to be attached to the prompt. The client MAY ignore this request.
	IncludeContext *CreateMessageRequestParamsIncludeContext `json:"includeContext,omitempty" yaml:"includeContext,omitempty" mapstructure:"includeContext,omitempty"`

	// The maximum number of tokens to sample, as requested by the server. The client
	// MAY choose to sample fewer tokens than requested.
	MaxTokens int `json:"maxTokens" yaml:"maxTokens" mapstructure:"maxTokens"`

	// Messages corresponds to the JSON schema field "messages".
	Messages []SamplingMessage `json:"messages" yaml:"messages" mapstructure:"messages"`

	// Optional metadata to pass through to the LLM provider. The format of this
	// metadata is provider-specific.
	Metadata CreateMessageRequestParamsMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	// The server's preferences for which model to select. The client MAY ignore these
	// preferences.
	ModelPreferences *ModelPreferences `json:"modelPreferences,omitempty" yaml:"modelPreferences,omitempty" mapstructure:"modelPreferences,omitempty"`

	// StopSequences corresponds to the JSON schema field "stopSequences".
	StopSequences []string `json:"stopSequences,omitempty" yaml:"stopSequences,omitempty" mapstructure:"stopSequences,omitempty"`

	// An optional system prompt the server wants to use for sampling. The client MAY
	// modify or omit this prompt.
	SystemPrompt *string `json:"systemPrompt,omitempty" yaml:"systemPrompt,omitempty" mapstructure:"systemPrompt,omitempty"`

	// Temperature corresponds to the JSON schema field "temperature".
	Temperature *float64 `json:"temperature,omitempty" yaml:"temperature,omitempty" mapstructure:"temperature,omitempty"`
}

type CreateMessageRequestParamsIncludeContext string

const CreateMessageRequestParamsIncludeContextAllServers CreateMessageRequestParamsIncludeContext = "allServers"
const CreateMessageRequestParamsIncludeContextNone CreateMessageRequestParamsIncludeContext = "none"
const CreateMessageRequestParamsIncludeContextThisServer CreateMessageRequestParamsIncludeContext = "thisServer"

var enumValues_CreateMessageRequestParamsIncludeContext = []interface{}{
	"allServers",
	"none",
	"thisServer",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageRequestParamsIncludeContext) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_CreateMessageRequestParamsIncludeContext {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_CreateMessageRequestParamsIncludeContext, v)
	}
	*j = CreateMessageRequestParamsIncludeContext(v)
	return nil
}

// Optional metadata to pass through to the LLM provider. The format of this
// metadata is provider-specific.
type CreateMessageRequestParamsMetadata map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["maxTokens"]; raw != nil && !ok {
		return fmt.Errorf("field maxTokens in CreateMessageRequestParams: required")
	}
	if _, ok := raw["messages"]; raw != nil && !ok {
		return fmt.Errorf("field messages in CreateMessageRequestParams: required")
	}
	type Plain CreateMessageRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateMessageRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CreateMessageRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CreateMessageRequest: required")
	}
	type Plain CreateMessageRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateMessageRequest(plain)
	return nil
}

// The client's response to a sampling/create_message request from the server. The
// client should inform the user before returning the sampled message, to allow
// them to inspect the response (human in the loop) and decide whether to allow the
// server to see it.
type CreateMessageResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta CreateMessageResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Content corresponds to the JSON schema field "content".
	Content interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// The name of the model that generated the message.
	Model string `json:"model" yaml:"model" mapstructure:"model"`

	// Role corresponds to the JSON schema field "role".
	Role Role `json:"role" yaml:"role" mapstructure:"role"`

	// The reason why sampling stopped, if known.
	StopReason *string `json:"stopReason,omitempty" yaml:"stopReason,omitempty" mapstructure:"stopReason,omitempty"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type CreateMessageResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in CreateMessageResult: required")
	}
	if _, ok := raw["model"]; raw != nil && !ok {
		return fmt.Errorf("field model in CreateMessageResult: required")
	}
	if _, ok := raw["role"]; raw != nil && !ok {
		return fmt.Errorf("field role in CreateMessageResult: required")
	}
	type Plain CreateMessageResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateMessageResult(plain)
	return nil
}

// An opaque token used to represent a cursor for pagination.
type Cursor string

// The contents of a resource, embedded into a prompt or tool call result.
//
// It is up to the client how best to render embedded resources for the benefit
// of the LLM and/or the user.
type EmbeddedResource struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *EmbeddedResourceAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// Resource corresponds to the JSON schema field "resource".
	Resource interface{} `json:"resource" yaml:"resource" mapstructure:"resource"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

type EmbeddedResourceAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EmbeddedResourceAnnotations) UnmarshalJSON(b []byte) error {
	type Plain EmbeddedResourceAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = EmbeddedResourceAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EmbeddedResource) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["resource"]; raw != nil && !ok {
		return fmt.Errorf("field resource in EmbeddedResource: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in EmbeddedResource: required")
	}
	type Plain EmbeddedResource
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = EmbeddedResource(plain)
	return nil
}

// Used by the client to get a prompt provided by the server.
type GetPromptRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params GetPromptRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type GetPromptRequestParams struct {
	// Arguments to use for templating the prompt.
	Arguments GetPromptRequestParamsArguments `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// The name of the prompt or prompt template.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// Arguments to use for templating the prompt.
type GetPromptRequestParamsArguments map[string]string

// UnmarshalJSON implements json.Unmarshaler.
func (j *GetPromptRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in GetPromptRequestParams: required")
	}
	type Plain GetPromptRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = GetPromptRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GetPromptRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in GetPromptRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in GetPromptRequest: required")
	}
	type Plain GetPromptRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = GetPromptRequest(plain)
	return nil
}

// The server's response to a prompts/get request from the client.
type GetPromptResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta GetPromptResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An optional description for the prompt.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Messages corresponds to the JSON schema field "messages".
	Messages []PromptMessage `json:"messages" yaml:"messages" mapstructure:"messages"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type GetPromptResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GetPromptResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["messages"]; raw != nil && !ok {
		return fmt.Errorf("field messages in GetPromptResult: required")
	}
	type Plain GetPromptResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = GetPromptResult(plain)
	return nil
}

// An image provided to or from an LLM.
type ImageContent struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *ImageContentAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// The base64-encoded image data.
	Data string `json:"data" yaml:"data" mapstructure:"data"`

	// The MIME type of the image. Different providers may support different image
	// types.
	MimeType string `json:"mimeType" yaml:"mimeType" mapstructure:"mimeType"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

type ImageContentAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ImageContentAnnotations) UnmarshalJSON(b []byte) error {
	type Plain ImageContentAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = ImageContentAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ImageContent) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["data"]; raw != nil && !ok {
		return fmt.Errorf("field data in ImageContent: required")
	}
	if _, ok := raw["mimeType"]; raw != nil && !ok {
		return fmt.Errorf("field mimeType in ImageContent: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ImageContent: required")
	}
	type Plain ImageContent
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ImageContent(plain)
	return nil
}

// Describes the name and version of an MCP implementation.
type Implementation struct {
	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Version corresponds to the JSON schema field "version".
	Version string `json:"version" yaml:"version" mapstructure:"version"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Implementation) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Implementation: required")
	}
	if _, ok := raw["version"]; raw != nil && !ok {
		return fmt.Errorf("field version in Implementation: required")
	}
	type Plain Implementation
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Implementation(plain)
	return nil
}

// This request is sent from the client to the server when it first connects,
// asking it to begin initialization.
type InitializeRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params InitializeRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type InitializeRequestParams struct {
	// Capabilities corresponds to the JSON schema field "capabilities".
	Capabilities ClientCapabilities `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`

	// ClientInfo corresponds to the JSON schema field "clientInfo".
	ClientInfo Implementation `json:"clientInfo" yaml:"clientInfo" mapstructure:"clientInfo"`

	// The latest version of the Model Context Protocol that the client supports. The
	// client MAY decide to support older versions as well.
	ProtocolVersion string `json:"protocolVersion" yaml:"protocolVersion" mapstructure:"protocolVersion"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializeRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities in InitializeRequestParams: required")
	}
	if _, ok := raw["clientInfo"]; raw != nil && !ok {
		return fmt.Errorf("field clientInfo in InitializeRequestParams: required")
	}
	if _, ok := raw["protocolVersion"]; raw != nil && !ok {
		return fmt.Errorf("field protocolVersion in InitializeRequestParams: required")
	}
	type Plain InitializeRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializeRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializeRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in InitializeRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in InitializeRequest: required")
	}
	type Plain InitializeRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializeRequest(plain)
	return nil
}

// After receiving an initialize request from the client, the server sends this
// response.
type InitializeResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta InitializeResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Capabilities corresponds to the JSON schema field "capabilities".
	Capabilities ServerCapabilities `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`

	// Instructions describing how to use the server and its features.
	//
	// This can be used by clients to improve the LLM's understanding of available
	// tools, resources, etc. It can be thought of like a "hint" to the model. For
	// example, this information MAY be added to the system prompt.
	Instructions *string `json:"instructions,omitempty" yaml:"instructions,omitempty" mapstructure:"instructions,omitempty"`

	// The version of the Model Context Protocol that the server wants to use. This
	// may not match the version that the client requested. If the client cannot
	// support this version, it MUST disconnect.
	ProtocolVersion string `json:"protocolVersion" yaml:"protocolVersion" mapstructure:"protocolVersion"`

	// ServerInfo corresponds to the JSON schema field "serverInfo".
	ServerInfo Implementation `json:"serverInfo" yaml:"serverInfo" mapstructure:"serverInfo"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type InitializeResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializeResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities in InitializeResult: required")
	}
	if _, ok := raw["protocolVersion"]; raw != nil && !ok {
		return fmt.Errorf("field protocolVersion in InitializeResult: required")
	}
	if _, ok := raw["serverInfo"]; raw != nil && !ok {
		return fmt.Errorf("field serverInfo in InitializeResult: required")
	}
	type Plain InitializeResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializeResult(plain)
	return nil
}

// This notification is sent from the client to the server after initialization has
// finished.
type InitializedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *InitializedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type InitializedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta InitializedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type InitializedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in InitializedNotification: required")
	}
	type Plain InitializedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializedNotification(plain)
	return nil
}

// A response to a request that indicates an error occurred.
type JSONRPCError struct {
	// Error corresponds to the JSON schema field "error".
	Error JSONRPCErrorData `json:"error" yaml:"error" mapstructure:"error"`

	// Id corresponds to the JSON schema field "id".
	Id interface{} `json:"id" yaml:"id" mapstructure:"id"`

	// Jsonrpc corresponds to the JSON schema field "jsonrpc".
	Jsonrpc string `json:"jsonrpc" yaml:"jsonrpc" mapstructure:"jsonrpc"`
}

type JSONRPCErrorData struct {
	// The error type that occurred.
	Code int `json:"code" yaml:"code" mapstructure:"code"`

	// Additional information about the error. The value of this member is defined by
	// the sender (e.g. detailed error information, nested errors etc.).
	Data interface{} `json:"data,omitempty" yaml:"data,omitempty" mapstructure:"data,omitempty"`

	// A short description of the error. The message SHOULD be limited to a concise
	// single sentence.
	Message string `json:"message" yaml:"message" mapstructure:"message"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JSONRPCErrorData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["code"]; raw != nil && !ok {
		return fmt.Errorf("field code in JSONRPCErrorData: required")
	}
	if _, ok := raw["message"]; raw != nil && !ok {
		return fmt.Errorf("field message in JSONRPCErrorData: required")
	}
	type Plain JSONRPCErrorData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = JSONRPCErrorData(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JSONRPCError) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["error"]; raw != nil && !ok {
		return fmt.Errorf("field error in JSONRPCError: required")
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in JSONRPCError: required")
	}
	if _, ok := raw["jsonrpc"]; raw != nil && !ok {
		return fmt.Errorf("field jsonrpc in JSONRPCError: required")
	}
	type Plain JSONRPCError
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = JSONRPCError(plain)
	return nil
}

type JSONRPCMessage interface{}

// A notification which does not expect a response.
type JSONRPCNotification struct {
	// Jsonrpc corresponds to the JSON schema field "jsonrpc".
	Jsonrpc string `json:"jsonrpc" yaml:"jsonrpc" mapstructure:"jsonrpc"`

	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *JSONRPCNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type JSONRPCNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta JSONRPCNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type JSONRPCNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JSONRPCNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["jsonrpc"]; raw != nil && !ok {
		return fmt.Errorf("field jsonrpc in JSONRPCNotification: required")
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in JSONRPCNotification: required")
	}
	type Plain JSONRPCNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = JSONRPCNotification(plain)
	return nil
}

// A request that expects a response.
type JSONRPCRequest struct {
	// Id corresponds to the JSON schema field "id".
	Id interface{} `json:"id" yaml:"id" mapstructure:"id"`

	// Jsonrpc corresponds to the JSON schema field "jsonrpc".
	Jsonrpc string `json:"jsonrpc" yaml:"jsonrpc" mapstructure:"jsonrpc"`

	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params json.RawMessage `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JSONRPCRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in JSONRPCRequest: required")
	}
	if _, ok := raw["jsonrpc"]; raw != nil && !ok {
		return fmt.Errorf("field jsonrpc in JSONRPCRequest: required")
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in JSONRPCRequest: required")
	}
	type Plain JSONRPCRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = JSONRPCRequest(plain)
	return nil
}

// A successful (non-error) response to a request.
type JSONRPCResponse struct {
	// Id corresponds to the JSON schema field "id".
	Id interface{} `json:"id" yaml:"id" mapstructure:"id"`

	// Jsonrpc corresponds to the JSON schema field "jsonrpc".
	Jsonrpc string `json:"jsonrpc" yaml:"jsonrpc" mapstructure:"jsonrpc"`

	// Result corresponds to the JSON schema field "result".
	Result interface{} `json:"result" yaml:"result" mapstructure:"result"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JSONRPCResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in JSONRPCResponse: required")
	}
	if _, ok := raw["jsonrpc"]; raw != nil && !ok {
		return fmt.Errorf("field jsonrpc in JSONRPCResponse: required")
	}
	if _, ok := raw["result"]; raw != nil && !ok {
		return fmt.Errorf("field result in JSONRPCResponse: required")
	}
	type Plain JSONRPCResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = JSONRPCResponse(plain)
	return nil
}

// Sent from the client to request a list of prompts and prompt templates the
// server has.
type ListPromptsRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListPromptsRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListPromptsRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListPromptsRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListPromptsRequest: required")
	}
	type Plain ListPromptsRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListPromptsRequest(plain)
	return nil
}

// The server's response to a prompts/list request from the client.
type ListPromptsResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListPromptsResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// Prompts corresponds to the JSON schema field "prompts".
	Prompts []Prompt `json:"prompts" yaml:"prompts" mapstructure:"prompts"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListPromptsResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListPromptsResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["prompts"]; raw != nil && !ok {
		return fmt.Errorf("field prompts in ListPromptsResult: required")
	}
	type Plain ListPromptsResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListPromptsResult(plain)
	return nil
}

// Sent from the client to request a list of resource templates the server has.
type ListResourceTemplatesRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListResourceTemplatesRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListResourceTemplatesRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourceTemplatesRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListResourceTemplatesRequest: required")
	}
	type Plain ListResourceTemplatesRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourceTemplatesRequest(plain)
	return nil
}

// The server's response to a resources/templates/list request from the client.
type ListResourceTemplatesResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListResourceTemplatesResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// ResourceTemplates corresponds to the JSON schema field "resourceTemplates".
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates" yaml:"resourceTemplates" mapstructure:"resourceTemplates"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListResourceTemplatesResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourceTemplatesResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["resourceTemplates"]; raw != nil && !ok {
		return fmt.Errorf("field resourceTemplates in ListResourceTemplatesResult: required")
	}
	type Plain ListResourceTemplatesResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourceTemplatesResult(plain)
	return nil
}

// Sent from the client to request a list of resources the server has.
type ListResourcesRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListResourcesRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListResourcesRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourcesRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListResourcesRequest: required")
	}
	type Plain ListResourcesRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourcesRequest(plain)
	return nil
}

// The server's response to a resources/list request from the client.
type ListResourcesResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListResourcesResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// Resources corresponds to the JSON schema field "resources".
	Resources []Resource `json:"resources" yaml:"resources" mapstructure:"resources"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListResourcesResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourcesResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["resources"]; raw != nil && !ok {
		return fmt.Errorf("field resources in ListResourcesResult: required")
	}
	type Plain ListResourcesResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourcesResult(plain)
	return nil
}

// Sent from the server to request a list of root URIs from the client. Roots allow
// servers to ask for specific directories or files to operate on. A common example
// for roots is providing a set of repositories or directories a server should
// operate
// on.
//
// This request is typically used when the server needs to understand the file
// system
// structure or access specific locations that the client has permission to read
// from.
type ListRootsRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListRootsRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListRootsRequestParams struct {
	// Meta corresponds to the JSON schema field "_meta".
	Meta *ListRootsRequestParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

type ListRootsRequestParamsMeta struct {
	// If specified, the caller is requesting out-of-band progress notifications for
	// this request (as represented by notifications/progress). The value of this
	// parameter is an opaque token that will be attached to any subsequent
	// notifications. The receiver is not obligated to provide these notifications.
	ProgressToken *ProgressToken `json:"progressToken,omitempty" yaml:"progressToken,omitempty" mapstructure:"progressToken,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListRootsRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListRootsRequest: required")
	}
	type Plain ListRootsRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListRootsRequest(plain)
	return nil
}

// The client's response to a roots/list request from the server.
// This result contains an array of Root objects, each representing a root
// directory
// or file that the server can operate on.
type ListRootsResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListRootsResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Roots corresponds to the JSON schema field "roots".
	Roots []Root `json:"roots" yaml:"roots" mapstructure:"roots"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListRootsResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListRootsResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["roots"]; raw != nil && !ok {
		return fmt.Errorf("field roots in ListRootsResult: required")
	}
	type Plain ListRootsResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListRootsResult(plain)
	return nil
}

// Sent from the client to request a list of tools the server has.
type ListToolsRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListToolsRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListToolsRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListToolsRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListToolsRequest: required")
	}
	type Plain ListToolsRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListToolsRequest(plain)
	return nil
}

// The server's response to a tools/list request from the client.
type ListToolsResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListToolsResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// Tools corresponds to the JSON schema field "tools".
	Tools []Tool `json:"tools" yaml:"tools" mapstructure:"tools"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListToolsResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListToolsResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["tools"]; raw != nil && !ok {
		return fmt.Errorf("field tools in ListToolsResult: required")
	}
	type Plain ListToolsResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListToolsResult(plain)
	return nil
}

type LoggingLevel string

const LoggingLevelAlert LoggingLevel = "alert"
const LoggingLevelCritical LoggingLevel = "critical"
const LoggingLevelDebug LoggingLevel = "debug"
const LoggingLevelEmergency LoggingLevel = "emergency"
const LoggingLevelError LoggingLevel = "error"
const LoggingLevelInfo LoggingLevel = "info"
const LoggingLevelNotice LoggingLevel = "notice"
const LoggingLevelWarning LoggingLevel = "warning"

var enumValues_LoggingLevel = []interface{}{
	"alert",
	"critical",
	"debug",
	"emergency",
	"error",
	"info",
	"notice",
	"warning",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingLevel) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_LoggingLevel {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_LoggingLevel, v)
	}
	*j = LoggingLevel(v)
	return nil
}

// Notification of a log message passed from server to client. If no
// logging/setLevel request has been sent from the client, the server MAY decide
// which messages to send automatically.
type LoggingMessageNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params LoggingMessageNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type LoggingMessageNotificationParams struct {
	// The data to be logged, such as a string message or an object. Any JSON
	// serializable type is allowed here.
	Data interface{} `json:"data" yaml:"data" mapstructure:"data"`

	// The severity of this log message.
	Level LoggingLevel `json:"level" yaml:"level" mapstructure:"level"`

	// An optional name of the logger issuing this message.
	Logger *string `json:"logger,omitempty" yaml:"logger,omitempty" mapstructure:"logger,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingMessageNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["data"]; raw != nil && !ok {
		return fmt.Errorf("field data in LoggingMessageNotificationParams: required")
	}
	if _, ok := raw["level"]; raw != nil && !ok {
		return fmt.Errorf("field level in LoggingMessageNotificationParams: required")
	}
	type Plain LoggingMessageNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = LoggingMessageNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingMessageNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in LoggingMessageNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in LoggingMessageNotification: required")
	}
	type Plain LoggingMessageNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = LoggingMessageNotification(plain)
	return nil
}

// Hints to use for model selection.
//
// Keys not declared here are currently left unspecified by the spec and are up
// to the client to interpret.
type ModelHint struct {
	// A hint for a model name.
	//
	// The client SHOULD treat this as a substring of a model name; for example:
	//  - `claude-3-5-sonnet` should match `claude-3-5-sonnet-20241022`
	//  - `sonnet` should match `claude-3-5-sonnet-20241022`,
	// `claude-3-sonnet-20240229`, etc.
	//  - `claude` should match any Claude model
	//
	// The client MAY also map the string to a different provider's model name or a
	// different model family, as long as it fills a similar niche; for example:
	//  - `gemini-1.5-flash` could match `claude-3-haiku-20240307`
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
}

// The server's preferences for model selection, requested of the client during
// sampling.
//
// Because LLMs can vary along multiple dimensions, choosing the "best" model is
// rarely straightforward.  Different models excel in different areas—some are
// faster but less capable, others are more capable but more expensive, and so
// on. This interface allows servers to express their priorities across multiple
// dimensions to help clients make an appropriate selection for their use case.
//
// These preferences are always advisory. The client MAY ignore them. It is also
// up to the client to decide how to interpret these preferences and how to
// balance them against other considerations.
type ModelPreferences struct {
	// How much to prioritize cost when selecting a model. A value of 0 means cost
	// is not important, while a value of 1 means cost is the most important
	// factor.
	CostPriority *float64 `json:"costPriority,omitempty" yaml:"costPriority,omitempty" mapstructure:"costPriority,omitempty"`

	// Optional hints to use for model selection.
	//
	// If multiple hints are specified, the client MUST evaluate them in order
	// (such that the first match is taken).
	//
	// The client SHOULD prioritize these hints over the numeric priorities, but
	// MAY still use the priorities to select from ambiguous matches.
	Hints []ModelHint `json:"hints,omitempty" yaml:"hints,omitempty" mapstructure:"hints,omitempty"`

	// How much to prioritize intelligence and capabilities when selecting a
	// model. A value of 0 means intelligence is not important, while a value of 1
	// means intelligence is the most important factor.
	IntelligencePriority *float64 `json:"intelligencePriority,omitempty" yaml:"intelligencePriority,omitempty" mapstructure:"intelligencePriority,omitempty"`

	// How much to prioritize sampling speed (latency) when selecting a model. A
	// value of 0 means speed is not important, while a value of 1 means speed is
	// the most important factor.
	SpeedPriority *float64 `json:"speedPriority,omitempty" yaml:"speedPriority,omitempty" mapstructure:"speedPriority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ModelPreferences) UnmarshalJSON(b []byte) error {
	type Plain ModelPreferences
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.CostPriority != nil && 1 < *plain.CostPriority {
		return fmt.Errorf("field %s: must be <= %v", "costPriority", 1)
	}
	if plain.CostPriority != nil && 0 > *plain.CostPriority {
		return fmt.Errorf("field %s: must be >= %v", "costPriority", 0)
	}
	if plain.IntelligencePriority != nil && 1 < *plain.IntelligencePriority {
		return fmt.Errorf("field %s: must be <= %v", "intelligencePriority", 1)
	}
	if plain.IntelligencePriority != nil && 0 > *plain.IntelligencePriority {
		return fmt.Errorf("field %s: must be >= %v", "intelligencePriority", 0)
	}
	if plain.SpeedPriority != nil && 1 < *plain.SpeedPriority {
		return fmt.Errorf("field %s: must be <= %v", "speedPriority", 1)
	}
	if plain.SpeedPriority != nil && 0 > *plain.SpeedPriority {
		return fmt.Errorf("field %s: must be >= %v", "speedPriority", 0)
	}
	*j = ModelPreferences(plain)
	return nil
}

type Notification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *NotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type NotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta NotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type NotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Notification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in Notification: required")
	}
	type Plain Notification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Notification(plain)
	return nil
}

type PaginatedRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *PaginatedRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type PaginatedRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PaginatedRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in PaginatedRequest: required")
	}
	type Plain PaginatedRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PaginatedRequest(plain)
	return nil
}

type PaginatedResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta PaginatedResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type PaginatedResultMeta map[string]interface{}

// A ping, issued by either the server or the client, to check that the other party
// is still alive. The receiver must promptly respond, or else may be disconnected.
type PingRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *PingRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type PingRequestParams struct {
	// Meta corresponds to the JSON schema field "_meta".
	Meta *PingRequestParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

type PingRequestParamsMeta struct {
	// If specified, the caller is requesting out-of-band progress notifications for
	// this request (as represented by notifications/progress). The value of this
	// parameter is an opaque token that will be attached to any subsequent
	// notifications. The receiver is not obligated to provide these notifications.
	ProgressToken *ProgressToken `json:"progressToken,omitempty" yaml:"progressToken,omitempty" mapstructure:"progressToken,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PingRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in PingRequest: required")
	}
	type Plain PingRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PingRequest(plain)
	return nil
}

// An out-of-band notification used to inform the receiver of a progress update for
// a long-running request.
type ProgressNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params ProgressNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type ProgressNotificationParams struct {
	// The progress thus far. This should increase every time progress is made, even
	// if the total is unknown.
	Progress float64 `json:"progress" yaml:"progress" mapstructure:"progress"`

	// The progress token which was given in the initial request, used to associate
	// this notification with the request that is proceeding.
	ProgressToken ProgressToken `json:"progressToken" yaml:"progressToken" mapstructure:"progressToken"`

	// Total number of items to process (or total progress required), if known.
	Total *float64 `json:"total,omitempty" yaml:"total,omitempty" mapstructure:"total,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProgressNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["progress"]; raw != nil && !ok {
		return fmt.Errorf("field progress in ProgressNotificationParams: required")
	}
	if _, ok := raw["progressToken"]; raw != nil && !ok {
		return fmt.Errorf("field progressToken in ProgressNotificationParams: required")
	}
	type Plain ProgressNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProgressNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProgressNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ProgressNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in ProgressNotification: required")
	}
	type Plain ProgressNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProgressNotification(plain)
	return nil
}

// A progress token, used to associate progress notifications with the original
// request.
type ProgressToken int

// A prompt or prompt template that the server offers.
type Prompt struct {
	// A list of arguments to use for templating the prompt.
	Arguments []PromptArgument `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// An optional description of what this prompt provides
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the prompt or prompt template.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// Describes an argument that a prompt can accept.
type PromptArgument struct {
	// A human-readable description of the argument.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the argument.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Whether this argument must be provided.
	Required bool `json:"required,omitempty" yaml:"required,omitempty" mapstructure:"required,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in PromptArgument: required")
	}
	type Plain PromptArgument
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptArgument(plain)
	return nil
}

// An optional notification from the server to the client, informing it that the
// list of prompts it offers has changed. This may be issued by servers without any
// previous subscription from the client.
type PromptListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *PromptListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type PromptListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta PromptListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type PromptListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in PromptListChangedNotification: required")
	}
	type Plain PromptListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptListChangedNotification(plain)
	return nil
}

// Describes a message returned as part of a prompt.
//
// This is similar to `SamplingMessage`, but also supports the embedding of
// resources from the MCP server.
type PromptMessage struct {
	// Content corresponds to the JSON schema field "content".
	Content interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// Role corresponds to the JSON schema field "role".
	Role Role `json:"role" yaml:"role" mapstructure:"role"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptMessage) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in PromptMessage: required")
	}
	if _, ok := raw["role"]; raw != nil && !ok {
		return fmt.Errorf("field role in PromptMessage: required")
	}
	type Plain PromptMessage
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptMessage(plain)
	return nil
}

// Identifies a prompt.
type PromptReference struct {
	// The name of the prompt or prompt template
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptReference) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in PromptReference: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in PromptReference: required")
	}
	type Plain PromptReference
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptReference(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Prompt) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Prompt: required")
	}
	type Plain Prompt
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Prompt(plain)
	return nil
}

// Sent from the client to the server, to read a specific resource URI.
type ReadResourceRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params ReadResourceRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type ReadResourceRequestParams struct {
	// The URI of the resource to read. The URI can use any protocol; it is up to the
	// server how to interpret it.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ReadResourceRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ReadResourceRequestParams: required")
	}
	type Plain ReadResourceRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ReadResourceRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ReadResourceRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ReadResourceRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in ReadResourceRequest: required")
	}
	type Plain ReadResourceRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ReadResourceRequest(plain)
	return nil
}

// The server's response to a resources/read request from the client.
type ReadResourceResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ReadResourceResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Contents corresponds to the JSON schema field "contents".
	Contents []interface{} `json:"contents" yaml:"contents" mapstructure:"contents"`
}

func (r *ReadResourceResult) AddTextContent(content TextResourceContents) {
	r.Contents = append(r.Contents, content)
}

func (r *ReadResourceResult) AddBlobContent(content BlobResourceContents) {
	r.Contents = append(r.Contents, content)
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ReadResourceResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ReadResourceResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["contents"]; raw != nil && !ok {
		return fmt.Errorf("field contents in ReadResourceResult: required")
	}
	type Plain ReadResourceResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ReadResourceResult(plain)
	return nil
}

type Request struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *RequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type RequestParams struct {
	// Meta corresponds to the JSON schema field "_meta".
	Meta *RequestParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

type RequestParamsMeta struct {
	// If specified, the caller is requesting out-of-band progress notifications for
	// this request (as represented by notifications/progress). The value of this
	// parameter is an opaque token that will be attached to any subsequent
	// notifications. The receiver is not obligated to provide these notifications.
	ProgressToken *ProgressToken `json:"progressToken,omitempty" yaml:"progressToken,omitempty" mapstructure:"progressToken,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Request) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in Request: required")
	}
	type Plain Request
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Request(plain)
	return nil
}

// A known resource that the server is capable of reading.
type Resource struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *ResourceAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// A description of what this resource represents.
	//
	// This can be used by clients to improve the LLM's understanding of available
	// resources. It can be thought of like a "hint" to the model.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The MIME type of this resource, if known.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// A human-readable name for this resource.
	//
	// This can be used by clients to populate UI elements.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The URI of this resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

type ResourceAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceAnnotations) UnmarshalJSON(b []byte) error {
	type Plain ResourceAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = ResourceAnnotations(plain)
	return nil
}

// The contents of a specific resource or sub-resource.
type ResourceContents struct {
	// The MIME type of this resource, if known.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// The URI of this resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceContents) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ResourceContents: required")
	}
	type Plain ResourceContents
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceContents(plain)
	return nil
}

// An optional notification from the server to the client, informing it that the
// list of resources it can read from has changed. This may be issued by servers
// without any previous subscription from the client.
type ResourceListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ResourceListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ResourceListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta ResourceListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type ResourceListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ResourceListChangedNotification: required")
	}
	type Plain ResourceListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceListChangedNotification(plain)
	return nil
}

// A reference to a resource or resource template definition.
type ResourceReference struct {
	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`

	// The URI or URI template of the resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceReference) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ResourceReference: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ResourceReference: required")
	}
	type Plain ResourceReference
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceReference(plain)
	return nil
}

// A template description for resources available on the server.
type ResourceTemplate struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *ResourceTemplateAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// A description of what this template is for.
	//
	// This can be used by clients to improve the LLM's understanding of available
	// resources. It can be thought of like a "hint" to the model.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The MIME type for all resources that match this template. This should only be
	// included if all resources matching this template have the same type.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// A human-readable name for the type of resource this template refers to.
	//
	// This can be used by clients to populate UI elements.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// A URI template (according to RFC 6570) that can be used to construct resource
	// URIs.
	UriTemplate string `json:"uriTemplate" yaml:"uriTemplate" mapstructure:"uriTemplate"`
}

type ResourceTemplateAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceTemplateAnnotations) UnmarshalJSON(b []byte) error {
	type Plain ResourceTemplateAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = ResourceTemplateAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceTemplate) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in ResourceTemplate: required")
	}
	if _, ok := raw["uriTemplate"]; raw != nil && !ok {
		return fmt.Errorf("field uriTemplate in ResourceTemplate: required")
	}
	type Plain ResourceTemplate
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceTemplate(plain)
	return nil
}

// A notification from the server to the client, informing it that a resource has
// changed and may need to be read again. This should only be sent if the client
// previously sent a resources/subscribe request.
type ResourceUpdatedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params ResourceUpdatedNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type ResourceUpdatedNotificationParams struct {
	// The URI of the resource that has been updated. This might be a sub-resource of
	// the one that the client actually subscribed to.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceUpdatedNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ResourceUpdatedNotificationParams: required")
	}
	type Plain ResourceUpdatedNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceUpdatedNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceUpdatedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ResourceUpdatedNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in ResourceUpdatedNotification: required")
	}
	type Plain ResourceUpdatedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceUpdatedNotification(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Resource) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Resource: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in Resource: required")
	}
	type Plain Resource
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Resource(plain)
	return nil
}

type Role string

const RoleAssistant Role = "assistant"
const RoleUser Role = "user"

var enumValues_Role = []interface{}{
	"assistant",
	"user",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Role) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_Role {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_Role, v)
	}
	*j = Role(v)
	return nil
}

// Represents a root directory or file that the server can operate on.
type Root struct {
	// An optional name for the root. This can be used to provide a human-readable
	// identifier for the root, which may be useful for display purposes or for
	// referencing the root in other parts of the application.
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`

	// The URI identifying the root. This *must* start with file:// for now.
	// This restriction may be relaxed in future versions of the protocol to allow
	// other URI schemes.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Root) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in Root: required")
	}
	type Plain Root
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Root(plain)
	return nil
}

// A notification from the client to the server, informing it that the list of
// roots has changed.
// This notification should be sent whenever the client adds, removes, or modifies
// any root.
// The server should then request an updated list of roots using the
// ListRootsRequest.
type RootsListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *RootsListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type RootsListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta RootsListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type RootsListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RootsListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in RootsListChangedNotification: required")
	}
	type Plain RootsListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = RootsListChangedNotification(plain)
	return nil
}

// Describes a message issued to or received from an LLM API.
type SamplingMessage struct {
	// Content corresponds to the JSON schema field "content".
	Content interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// Role corresponds to the JSON schema field "role".
	Role Role `json:"role" yaml:"role" mapstructure:"role"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SamplingMessage) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in SamplingMessage: required")
	}
	if _, ok := raw["role"]; raw != nil && !ok {
		return fmt.Errorf("field role in SamplingMessage: required")
	}
	type Plain SamplingMessage
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SamplingMessage(plain)
	return nil
}

// Capabilities that a server may support. Known capabilities are defined here, in
// this schema, but this is not a closed set: any server can define its own,
// additional capabilities.
type ServerCapabilities struct {
	// Experimental, non-standard capabilities that the server supports.
	Experimental ServerCapabilitiesExperimental `json:"experimental,omitempty" yaml:"experimental,omitempty" mapstructure:"experimental,omitempty"`

	// Present if the server supports sending log messages to the client.
	Logging ServerCapabilitiesLogging `json:"logging,omitempty" yaml:"logging,omitempty" mapstructure:"logging,omitempty"`

	// Present if the server offers any prompt templates.
	Prompts *ServerCapabilitiesPrompts `json:"prompts,omitempty" yaml:"prompts,omitempty" mapstructure:"prompts,omitempty"`

	// Present if the server offers any resources to read.
	Resources *ServerCapabilitiesResources `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// Present if the server offers any tools to call.
	Tools *ServerCapabilitiesTools `json:"tools,omitempty" yaml:"tools,omitempty" mapstructure:"tools,omitempty"`
}

// Experimental, non-standard capabilities that the server supports.
type ServerCapabilitiesExperimental map[string]map[string]interface{}

// Present if the server supports sending log messages to the client.
type ServerCapabilitiesLogging map[string]interface{}

// Present if the server offers any prompt templates.
type ServerCapabilitiesPrompts struct {
	// Whether this server supports notifications for changes to the prompt list.
	ListChanged bool `json:"listChanged,omitempty" yaml:"listChanged,omitempty" mapstructure:"listChanged,omitempty"`
}

// Present if the server offers any resources to read.
type ServerCapabilitiesResources struct {
	// Whether this server supports notifications for changes to the resource list.
	ListChanged bool `json:"listChanged,omitempty" yaml:"listChanged,omitempty" mapstructure:"listChanged,omitempty"`

	// Whether this server supports subscribing to resource updates.
	Subscribe bool `json:"subscribe,omitempty" yaml:"subscribe,omitempty" mapstructure:"subscribe,omitempty"`
}

// Present if the server offers any tools to call.
type ServerCapabilitiesTools struct {
	// Whether this server supports notifications for changes to the tool list.
	ListChanged bool `json:"listChanged,omitempty" yaml:"listChanged,omitempty" mapstructure:"listChanged,omitempty"`
}

type ServerNotification interface{}

type ServerRequest interface{}

type ServerResult interface{}

// A request from the client to the server, to enable or adjust logging.
type SetLevelRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params SetLevelRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type SetLevelRequestParams struct {
	// The level of logging that the client wants to receive from the server. The
	// server should send all logs at this level and higher (i.e., more severe) to the
	// client as notifications/logging/message.
	Level LoggingLevel `json:"level" yaml:"level" mapstructure:"level"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SetLevelRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["level"]; raw != nil && !ok {
		return fmt.Errorf("field level in SetLevelRequestParams: required")
	}
	type Plain SetLevelRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SetLevelRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SetLevelRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in SetLevelRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in SetLevelRequest: required")
	}
	type Plain SetLevelRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SetLevelRequest(plain)
	return nil
}

// Sent from the client to request resources/updated notifications from the server
// whenever a particular resource changes.
type SubscribeRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params SubscribeRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type SubscribeRequestParams struct {
	// The URI of the resource to subscribe to. The URI can use any protocol; it is up
	// to the server how to interpret it.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SubscribeRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in SubscribeRequestParams: required")
	}
	type Plain SubscribeRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SubscribeRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SubscribeRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in SubscribeRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in SubscribeRequest: required")
	}
	type Plain SubscribeRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SubscribeRequest(plain)
	return nil
}

// Text provided to or from an LLM.
type TextContent struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *TextContentAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// The text content of the message.
	Text string `json:"text" yaml:"text" mapstructure:"text"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

type TextContentAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TextContentAnnotations) UnmarshalJSON(b []byte) error {
	type Plain TextContentAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = TextContentAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TextContent) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["text"]; raw != nil && !ok {
		return fmt.Errorf("field text in TextContent: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in TextContent: required")
	}
	type Plain TextContent
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = TextContent(plain)
	return nil
}

type TextResourceContents struct {
	// The MIME type of this resource, if known.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// The text of the item. This must only be set if the item can actually be
	// represented as text (not binary data).
	Text string `json:"text" yaml:"text" mapstructure:"text"`

	// The URI of this resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TextResourceContents) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["text"]; raw != nil && !ok {
		return fmt.Errorf("field text in TextResourceContents: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in TextResourceContents: required")
	}
	type Plain TextResourceContents
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = TextResourceContents(plain)
	return nil
}

// Definition for a tool the client can call.
type Tool struct {
	// A human-readable description of the tool.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// A JSON Schema object defining the expected parameters for the tool.
	InputSchema ToolInputSchema `json:"inputSchema" yaml:"inputSchema" mapstructure:"inputSchema"`

	// The name of the tool.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// A JSON Schema object defining the expected parameters for the tool.
type ToolInputSchema struct {
	// Properties corresponds to the JSON schema field "properties".
	Properties ToolInputSchemaProperties `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

type ToolInputSchemaProperties map[string]map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolInputSchema) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ToolInputSchema: required")
	}
	type Plain ToolInputSchema
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ToolInputSchema(plain)
	return nil
}

// An optional notification from the server to the client, informing it that the
// list of tools it offers has changed. This may be issued by servers without any
// previous subscription from the client.
type ToolListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ToolListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ToolListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta ToolListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type ToolListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ToolListChangedNotification: required")
	}
	type Plain ToolListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ToolListChangedNotification(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Tool) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["inputSchema"]; raw != nil && !ok {
		return fmt.Errorf("field inputSchema in Tool: required")
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Tool: required")
	}
	type Plain Tool
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Tool(plain)
	return nil
}

// Sent from the client to request cancellation of resources/updated notifications
// from the server. This should follow a previous resources/subscribe request.
type UnsubscribeRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params UnsubscribeRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type UnsubscribeRequestParams struct {
	// The URI of the resource to unsubscribe from.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *UnsubscribeRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in UnsubscribeRequestParams: required")
	}
	type Plain UnsubscribeRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = UnsubscribeRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *UnsubscribeRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in UnsubscribeRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in UnsubscribeRequest: required")
	}
	type Plain UnsubscribeRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = UnsubscribeRequest(plain)
	return nil
}