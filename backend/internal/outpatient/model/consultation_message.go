package model

import (
	"time"
)

// ConsultationMessage 问诊消息
type ConsultationMessage struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	ConsultationID string    `gorm:"column:consultation_id;not null;type:varchar(64);index" json:"consultationId"`
	SenderID       string    `gorm:"column:sender_id;not null;type:varchar(64)" json:"senderId"`
	SenderName     string    `gorm:"column:sender_name;size:50" json:"senderName,omitempty"`
	Content        string    `gorm:"column:content;type:text" json:"content"`
	MsgType        string    `gorm:"column:msg_type;size:20" json:"msgType"` // text/image/voice
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (ConsultationMessage) TableName() string { return "consultation_messages" }
