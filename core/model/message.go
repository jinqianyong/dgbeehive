/**
  @author:jinqianyong
  @date:2/14/23
*/
package model

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	InsertOperation        = "insert"
	DeleteOperation        = "delete"
	QueryOperation         = "query"
	UpdateOperation        = "update"
	ResponseOperation      = "response"
	ResponseErrorOperation = "error"

	ResourceTypePod          = "pod"
	ResourceTypeConfigmap    = "configmap"
	ResourceTypeNode         = "node"
	ResourceTypePodlist      = "podlist"
	ResourceTypePodStatus    = "podstatus"
	ResourceTypeRule         = "rule"
	ResourceTypeRuleEndpoint = "ruleendpoint"
	ResourceTypeRuleStatus   = "rulestatus"
)

// Message struct
type Message struct {
	Header  MessageHeader `json:"header"`
	Router  MessageRoute  `json:"router"`
	Content interface{}   `json:"content"`
}
type MessageRoute struct {
	// where the message come from
	Source string `json:"source"`
	// where the message will broadcast to
	Group string `json:"group,omitempty"`
	//what's the operation on resource
	Operation string `json:"operation,omitempty"`
	//what's the resource want to operate
	Resource string `json:"resource,omitempty"`
}

type MessageHeader struct {
	//the message uuid
	ID string `json:"msg_id"`
	// the response message parentid must be same with message received
	// please use NewRespByMessage to new response message
	ParentID string `json:"parent_msg_id,omitempty"`
	// the time of creating
	Timestamp int64 `json:"timestamp"`
	// specific resource version for the message, if any
	// it's currently backed by resource version of the K8S object saved in the content field
	// kubeedge leverages the concept of message resource version to achieve reliable transmission
	ResourceVersion string `json:"resourceversion,omitempty"`
	// the flag will be set in sendsync
	Sync bool `json:"sync,omitempty"`
}

// BuildRouter sets route and resource operation in message
func (msg *Message) BuildRouter(source, group, res, opr string) *Message {
	msg.SetRoute(source, group)
	msg.SetResourceOperation(res, opr)
	return msg
}

// SetResourceOperation sets resource version in message header
func (msg *Message) SetResourceOperation(res, opr string) *Message {
	msg.Router.Resource = res
	msg.Router.Operation = opr
	return msg
}

// SetRoute sets source and group in message
func (msg *Message) SetRoute(source, group string) *Message {
	msg.Router.Source = source
	msg.Router.Group = group
	return msg
}

// SetResourceVersion sets resource in message header
func (msg *Message) SetResourceVersion(resourceVersion string) *Message {
	msg.Header.ResourceVersion = resourceVersion
	return msg
}

// IsSync: msg.Header.Sync will be set in sendsync
func (msg *Message) IsSync() bool {
	return msg.Header.Sync
}

// GetResource returns message route resource
func (msg *Message) GetResource() string {
	return msg.Router.Resource
}
func (msg *Message) GetOperation() string {
	return msg.Router.Operation
}
func (msg *Message) GetSource() string {
	return msg.Router.Source
}
func (msg *Message) GetGroup() string {
	return msg.Router.Group
}
func (msg *Message) GetID() string {
	return msg.Header.ID
}
func (msg *Message) GetParentID() string {
	return msg.Header.ParentID
}
func (msg *Message) GetTimestamp() int64 {
	return msg.Header.Timestamp
}
func (msg *Message) GetContent() interface{} {
	return msg.Content
}
func (msg *Message) GetContentData() ([]byte, error) {
	if data, ok := msg.Content.([]byte); ok {
		return data, nil
	}
	data, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, fmt.Errorf("marshal message content failed:%s", err)
	}
	return data, nil
}
func (msg *Message) GetResourceVersion() string {
	return msg.Header.ResourceVersion
}
func (msg *Message) UpdateID() *Message {
	msg.Header.ID = uuid.NewV4().String()
	return msg
}
func (msg *Message) BuildHeader(ID, parentID string, timestamp int64) *Message {
	msg.Header.ID = ID
	msg.Header.ParentID = parentID
	msg.Header.Timestamp = timestamp
	return msg
}
func (msg *Message) FillBody(content interface{}) *Message {
	msg.Content = content
	return msg
}
func NewRawMessage() *Message {
	return &Message{}
}

func NewMessage(parentID string) *Message {
	msg := &Message{}
	msg.Header.ID = uuid.NewV4().String()
	msg.Header.ParentID = parentID
	msg.Header.Timestamp = time.Now().UnixNano() / 1e6
	return msg
}
func (msg *Message) Clone(message *Message) *Message {
	msgID := uuid.NewV4().String()
	return NewRawMessage().BuildHeader(msgID, message.GetParentID(), message.GetTimestamp()).FillBody(message.GetContent())
}

func (msg *Message) NewRespByMessage(message *Message, content interface{}) *Message {
	return NewMessage(message.GetID()).SetRoute(message.GetResource(), message.GetGroup()).
		SetResourceOperation(message.GetResource(), ResponseOperation).FillBody(content)
}
func NewErrorMessage(message *Message, errContent string) *Message {
	return NewMessage(message.Header.ParentID).SetResourceOperation(message.Router.Resource, ResponseErrorOperation).FillBody(errContent)
}
