/**
  @author:jinqianyong
  @date:2/14/23
*/
package context

import (
	gocontext "context"
	"dgbeehive/core/model"
	"time"
)

/**
  ModuleContext is interface for context module management
*/
type ModuleContext interface {
	AddModule(module string)
	AddModuleGroup(module, group string)
	Cleanup(module string)
}

type MessageContext interface {
	// async mode
	Send(module string, message model.Message)
	Receive(module string) (model.Message, error)
	// sync mode
	SendSync(module string, message model.Message, timeout time.Duration) (model.Message, error)
	SendResp(message model.Message)

	//group broadcast
	SendToGroup(moduleType string, message model.Message)
	SendToGroupSync(moduleType string, message model.Message, timeout time.Duration) error
}

type beehiveContext struct {
	moduleContext  ModuleContext
	messageContext MessageContext
	ctx            gocontext.Context
	cancel         gocontext.CancelFunc
}
