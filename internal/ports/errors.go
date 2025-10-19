package ports

import "fmt"

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

var (
	// User errors
	ErrUserNotFound      = &ApiError{StatusCode: 404, Message: "User not found"}
	ErrUserAlreadyExists = &ApiError{StatusCode: 409, Message: "User with this email already exists"}

	// Auth errors
	ErrInvalidCredentials = &ApiError{StatusCode: 401, Message: "Invalid email or password"}
	ErrIncorrectPassword  = &ApiError{StatusCode: 401, Message: "Current password is incorrect"}
	ErrPasswordComplexity = &ApiError{StatusCode: 422, Message: "Password does not meet complexity requirements"}
	ErrAccountDeactivated = &ApiError{StatusCode: 401, Message: "Account is deactivated"}
	ErrInvalidToken       = &ApiError{StatusCode: 401, Message: "Invalid or expired authentication token"}
	ErrInvalidSession     = &ApiError{StatusCode: 401, Message: "Invalid or expired session"}
	ErrTokenGeneration    = &ApiError{StatusCode: 500, Message: "Failed to generate authentication token"}

	// Business errors
	ErrBusinessNotFound = &ApiError{StatusCode: 404, Message: "Business not found"}
	ErrBusinessInUse    = &ApiError{StatusCode: 409, Message: "Cannot delete business, it is currently in use"}
	ErrOperatorNotFound = &ApiError{StatusCode: 404, Message: "Operator user not found"}

	// Project errors
	ErrProjectNotFound     = &ApiError{StatusCode: 404, Message: "Project not found"}
	ErrProjectNameExists   = &ApiError{StatusCode: 409, Message: "A project with this name already exists"}
	ErrMemberAlreadyExists = &ApiError{StatusCode: 409, Message: "User is already a member of this project"}
	ErrMemberNotFound      = &ApiError{StatusCode: 404, Message: "Project member not found"}
	ErrManagerNotFound     = &ApiError{StatusCode: 400, Message: "Manager user not found"}

	// Skill errors
	ErrSkillNotFound   = &ApiError{StatusCode: 404, Message: "Skill not found"}
	ErrSkillNameExists = &ApiError{StatusCode: 409, Message: "A skill with this name already exists"}
	ErrSkillInUse      = &ApiError{StatusCode: 409, Message: "Cannot delete skill, it is currently in use"}

	// Publication errors
	ErrPublicationNotFound       = &ApiError{StatusCode: 404, Message: "Publication not found"}
	ErrPublicationSlugExists     = &ApiError{StatusCode: 409, Message: "A publication with this title/slug already exists"}
	ErrPublicationAuthorNotFound = &ApiError{StatusCode: 400, Message: "Author user not found"}

	// Idea errors
	ErrIdeaNotFound          = &ApiError{StatusCode: 404, Message: "Idea not found"}
	ErrIdeaSubmitterNotFound = &ApiError{StatusCode: 400, Message: "Submitter user not found"}

	// Notifcation errors
	ErrNotificationNotFound = &ApiError{StatusCode: 404, Message: "Notification not found"}
	ErrReceiverNotFound     = &ApiError{StatusCode: 400, Message: "Notification receiver not found"}

	// UserSkill errors
	ErrUserSkillNotFound      = &ApiError{StatusCode: 404, Message: "User skill not found"}
	ErrUserSkillAlreadyExists = &ApiError{StatusCode: 409, Message: "User already has this skill"}
	ErrInvalidProficiency     = &ApiError{StatusCode: 400, Message: "Invalid proficiency level"}

	// ProjectSkill errors
	ErrProjectSkillNotFound      = &ApiError{StatusCode: 404, Message: "Project skill not found"}
	ErrProjectSkillAlreadyExists = &ApiError{StatusCode: 409, Message: "Project already has this skill"}
	ErrInvalidImportance         = &ApiError{StatusCode: 400, Message: "Invalid importance level"}

	// ProjectMember errors
	ErrProjectMemberNotFound      = &ApiError{StatusCode: 404, Message: "Project member not found"}
	ErrProjectMemberAlreadyExists = &ApiError{StatusCode: 409, Message: "User is already a member of this project"}
	ErrCannotRemoveManager        = &ApiError{StatusCode: 400, Message: "Cannot remove project manager"}
	ErrInvalidRole                = &ApiError{StatusCode: 400, Message: "Invalid role"}

	// BusinessConnection errors
	ErrBusinessConnectionNotFound      = &ApiError{StatusCode: 404, Message: "Business connection not found"}
	ErrBusinessConnectionAlreadyExists = &ApiError{StatusCode: 409, Message: "Business connection already exists"}
	ErrInvalidConnectionType           = &ApiError{StatusCode: 400, Message: "Invalid connection type"}
	ErrInvalidConnectionStatus         = &ApiError{StatusCode: 400, Message: "Invalid connection status"}
	ErrCannotConnectToSelf             = &ApiError{StatusCode: 400, Message: "Cannot create connection to same business"}
	ErrConnectionNotPending            = &ApiError{StatusCode: 400, Message: "Connection is not in pending status"}

	// BusinessTag errors
	ErrBusinessTagNotFound      = &ApiError{StatusCode: 404, Message: "Business tag not found"}
	ErrBusinessTagAlreadyExists = &ApiError{StatusCode: 409, Message: "Business tag already exists"}
	ErrInvalidTagType           = &ApiError{StatusCode: 400, Message: "Invalid tag type"}

	// IdeaVote errors
	ErrIdeaVoteNotFound      = &ApiError{StatusCode: 404, Message: "Idea vote not found"}
	ErrIdeaVoteAlreadyExists = &ApiError{StatusCode: 409, Message: "User has already voted on this idea"}
	ErrInvalidVoteType       = &ApiError{StatusCode: 400, Message: "Invalid vote type"}
	ErrCannotVoteOwnIdea     = &ApiError{StatusCode: 400, Message: "Cannot vote on your own idea"}

	// UserSession errors
	ErrSessionNotFound = &ApiError{StatusCode: 404, Message: "Session not found"}

	// Generic errors
	ErrDatabase     = &ApiError{StatusCode: 500, Message: "A database error occurred"}
	ErrNoUpdateData = &ApiError{StatusCode: 400, Message: "No valid fields provided for update"}
)
