package types

type Algorithm interface {
	InitializeNode()
	StartElections()
	SendHeartbeat()
	GetServer() any
}
