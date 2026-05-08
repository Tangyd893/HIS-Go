package model

import (
	"time"
)

// ConsultationMessage 问诊消息
type ConsultationMessage struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	ConsultationID string    `gorm:"column:consultation_id;not null;type:varchar(64);index"`
	SenderID       string    `gorm:"column:sender_id;not null;type:varchar(64)"`
	SenderName     string    `gorm:"column:sender_name;size:50"`
	Content        string    `gorm:"column:content;type:text"`
	MsgType        string    `gorm:"column:msg_type;size:20"` // text/image/voice
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ConsultationMessage) TableName() string { return "consultation_messages" }
