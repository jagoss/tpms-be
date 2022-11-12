package interfaces

type MessageSender interface {
	SendToEnabledUsers(userID string, data map[string]string) error
}
