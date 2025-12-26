package models

type NodeType string

const (
	NodeTypeTrigger   NodeType = "TRIGGER"
	NodeTypeAction    NodeType = "ACTION"
	NodeTypeCondition NodeType = "CONDITION"
)

type ActionType string

const (
	ActionTypeSendEmail ActionType = "Send Email"
	ActionTypeFail      ActionType = "FAIL"
)

type TriggerType string

const (
	TriggerTypeEvent    TriggerType = "EVENT"
	TriggerTypeSchedule TriggerType = "SCHEDULE"
)
