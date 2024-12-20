package domain

const (
	RoleCreator     = "creator"     //	the creator of the event
	RoleAdmin       = "admin"       //	a participant with admin permissions
	RoleParticipant = "participant" //  a normal participant
	RoleAssistant   = "Assistant"   //	only an assistant
)

const (
	PermissionAssignTask     = "permission to assign task"
	PermissionAddParticipant = "permission to add a participant"
	PermissionInvitePeople   = "permission to invite people"
	PermissionSendMessages   = "permission to use the chat"
	PermissionEditEvent      = "permission to edit the event"
)

var RolesPermissions = map[string][]string{
	RoleCreator: {
		PermissionAssignTask,
		PermissionAddParticipant,
		PermissionInvitePeople,
		PermissionSendMessages,
		PermissionEditEvent,
	},
	RoleAdmin: {
		PermissionAssignTask,
		PermissionAddParticipant,
		PermissionInvitePeople,
		PermissionSendMessages,
	},
	RoleParticipant: {
		PermissionInvitePeople,
		PermissionSendMessages,
	},
	RoleAssistant: {},
}