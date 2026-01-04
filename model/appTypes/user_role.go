package appTypes

// RoleID 用户角色
type RoleID int

const (
	Guest RoleID = iota
	User
	Admin
)
