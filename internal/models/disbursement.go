package models

import "time"

type DisbursementStatus string

const (
	StatusPending  DisbursementStatus = "pending"
	StatusApproved DisbursementStatus = "approved"
	StatusRejected DisbursementStatus = "rejected"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Role         string    `gorm:"size:30;not null" json:"role"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Disbursement struct {
	ID              uint               `gorm:"primaryKey" json:"id"`
	RequesterID     uint               `gorm:"not null;index" json:"requester_id"`
	Requester       *User              `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	BeneficiaryName string             `gorm:"size:100;not null" json:"beneficiary_name"`
	BankName        string             `gorm:"size:100;not null" json:"bank_name"`
	AccountNumber   string             `gorm:"size:30;not null" json:"account_number"`
	Amount          float64            `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description     string             `gorm:"type:text" json:"description"`
	Status          DisbursementStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	ProcessedByID   *uint              `gorm:"index" json:"processed_by_id,omitempty"`
	ProcessedBy     *User              `gorm:"foreignKey:ProcessedByID" json:"processed_by,omitempty"`
	RejectionReason *string            `gorm:"type:text" json:"rejection_reason,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}
