package role

type Role uint8

const (
	User Role = iota + 1
	Assistant
)

//go:generate go-enum -type=Role
