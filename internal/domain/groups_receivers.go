package domain

type GroupsReceivers struct {
	GroupID    string `json:"group_id" db:"group_id"`
	ReceiverID string `json:"receiver_id" db:"receiver_id"`
}
