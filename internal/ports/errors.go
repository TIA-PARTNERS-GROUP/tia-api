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
	
	ErrUserNotFound      = &ApiError{StatusCode: 404, Message: "User not found"}
	ErrUserAlreadyExists = &ApiError{StatusCode: 409, Message: "User with this email already exists"}
	
	ErrInvalidCredentials = &ApiError{StatusCode: 401, Message: "Invalid email or password"}
	ErrIncorrectPassword  = &ApiError{StatusCode: 401, Message: "Current password is incorrect"}
	ErrPasswordComplexity = &ApiError{StatusCode: 422, Message: "Password does not meet complexity requirements"}
	ErrAccountDeactivated = &ApiError{StatusCode: 401, Message: "Account is deactivated"}
	ErrInvalidToken       = &ApiError{StatusCode: 401, Message: "Invalid or expired authentication token"}
	ErrInvalidSession     = &ApiError{StatusCode: 401, Message: "Invalid or expired session"}
	ErrTokenGeneration    = &ApiError{StatusCode: 500, Message: "Failed to generate authentication token"}
	
	ErrBusinessNotFound = &ApiError{StatusCode: 404, Message: "Business not found"}
	ErrBusinessInUse    = &ApiError{StatusCode: 409, Message: "Cannot delete business, it is currently in use"}
	ErrOperatorNotFound = &ApiError{StatusCode: 404, Message: "Operator user not found"}
	
	ErrProjectNotFound     = &ApiError{StatusCode: 404, Message: "Project not found"}
	ErrProjectNameExists   = &ApiError{StatusCode: 409, Message: "A project with this name already exists"}
	ErrMemberAlreadyExists = &ApiError{StatusCode: 409, Message: "User is already a member of this project"}
	ErrMemberNotFound      = &ApiError{StatusCode: 404, Message: "Project member not found"}
	ErrManagerNotFound     = &ApiError{StatusCode: 400, Message: "Manager user not found"}
	
	ErrSkillNotFound   = &ApiError{StatusCode: 404, Message: "Skill not found"}
	ErrSkillNameExists = &ApiError{StatusCode: 409, Message: "A skill with this name already exists"}
	ErrSkillInUse      = &ApiError{StatusCode: 409, Message: "Cannot delete skill, it is currently in use"}
	
	ErrPublicationNotFound       = &ApiError{StatusCode: 404, Message: "Publication not found"}
	ErrPublicationSlugExists     = &ApiError{StatusCode: 409, Message: "A publication with this title/slug already exists"}
	ErrPublicationAuthorNotFound = &ApiError{StatusCode: 400, Message: "Author user not found"}
	
	ErrIdeaNotFound          = &ApiError{StatusCode: 404, Message: "Idea not found"}
	ErrIdeaSubmitterNotFound = &ApiError{StatusCode: 400, Message: "Submitter user not found"}
	
	ErrNotificationNotFound = &ApiError{StatusCode: 404, Message: "Notification not found"}
	ErrReceiverNotFound     = &ApiError{StatusCode: 400, Message: "Notification receiver not found"}
	
	ErrUserSkillNotFound      = &ApiError{StatusCode: 404, Message: "User skill not found"}
	ErrUserSkillAlreadyExists = &ApiError{StatusCode: 409, Message: "User already has this skill"}
	ErrInvalidProficiency     = &ApiError{StatusCode: 400, Message: "Invalid proficiency level"}
	
	ErrProjectSkillNotFound      = &ApiError{StatusCode: 404, Message: "Project skill not found"}
	ErrProjectSkillAlreadyExists = &ApiError{StatusCode: 409, Message: "Project already has this skill"}
	ErrInvalidImportance         = &ApiError{StatusCode: 400, Message: "Invalid importance level"}
	
	ErrProjectMemberNotFound      = &ApiError{StatusCode: 404, Message: "Project member not found"}
	ErrProjectMemberAlreadyExists = &ApiError{StatusCode: 409, Message: "User is already a member of this project"}
	ErrCannotRemoveManager        = &ApiError{StatusCode: 400, Message: "Cannot remove project manager"}
	ErrInvalidRole                = &ApiError{StatusCode: 400, Message: "Invalid role"}
	
	ErrBusinessConnectionNotFound      = &ApiError{StatusCode: 404, Message: "Business connection not found"}
	ErrBusinessConnectionAlreadyExists = &ApiError{StatusCode: 409, Message: "Business connection already exists"}
	ErrInvalidConnectionType           = &ApiError{StatusCode: 400, Message: "Invalid connection type"}
	ErrInvalidConnectionStatus         = &ApiError{StatusCode: 400, Message: "Invalid connection status"}
	ErrCannotConnectToSelf             = &ApiError{StatusCode: 400, Message: "Cannot create connection to same business"}
	ErrConnectionNotPending            = &ApiError{StatusCode: 400, Message: "Connection is not in pending status"}
	
	ErrBusinessTagNotFound      = &ApiError{StatusCode: 404, Message: "Business tag not found"}
	ErrBusinessTagAlreadyExists = &ApiError{StatusCode: 409, Message: "Business tag already exists"}
	ErrInvalidTagType           = &ApiError{StatusCode: 400, Message: "Invalid tag type"}
	
	ErrIdeaVoteNotFound      = &ApiError{StatusCode: 404, Message: "Idea vote not found"}
	ErrIdeaVoteAlreadyExists = &ApiError{StatusCode: 409, Message: "User has already voted on this idea"}
	ErrInvalidVoteType       = &ApiError{StatusCode: 400, Message: "Invalid vote type"}
	ErrCannotVoteOwnIdea     = &ApiError{StatusCode: 400, Message: "Cannot vote on your own idea"}
	
	ErrSessionNotFound = &ApiError{StatusCode: 404, Message: "Session not found"}
	
	ErrL2EResponseNotFound = &ApiError{StatusCode: 404, Message: "L2E response not found"}
	
	ErrFeedbackNotFound = &ApiError{StatusCode: 404, Message: "Feedback not found"}
	
	ErrAlreadyApplied        = &ApiError{StatusCode: 409, Message: "User has already applied to this project"}
	ErrApplicationNotFound   = &ApiError{StatusCode: 404, Message: "Project application not found"}
	ErrProjectOrUserNotFound = &ApiError{StatusCode: 400, Message: "Project or user not found"}
	
	ErrDailyActivityNotFound = &ApiError{StatusCode: 404, Message: "Daily activity not found"}
	ErrActivityNameExists    = &ApiError{StatusCode: 409, Message: "An activity with this name already exists"}
	ErrAlreadyEnrolled       = &ApiError{StatusCode: 409, Message: "User is already enrolled in this activity"}
	ErrEnrolmentNotFound     = &ApiError{StatusCode: 404, Message: "User is not enrolled in this activity"}
	
	ErrEventNotFound = &ApiError{StatusCode: 404, Message: "Event not found"}
	
	ErrSubscriptionNotFound   = &ApiError{StatusCode: 404, Message: "Subscription plan not found"}
	ErrSubscriptionNameExists = &ApiError{StatusCode: 409, Message: "A subscription plan with this name already exists"}
	
	ErrUserSubscriptionNotFound = &ApiError{StatusCode: 404, Message: "User subscription not found"}
	
	ErrUserConfigNotFound = &ApiError{StatusCode: 404, Message: "User configuration not found"}
	
	ErrRegionAlreadyAdded      = &ApiError{StatusCode: 409, Message: "Region is already associated with this project"}
	ErrProjectRegionNotFound   = &ApiError{StatusCode: 404, Message: "This project is not associated with the specified region"}
	ErrProjectOrRegionNotFound = &ApiError{StatusCode: 400, Message: "Project or Region not found"}
	ErrForbidden = &ApiError{StatusCode: 403, Message: "Forbidden"}
	
	ErrDatabase     = &ApiError{StatusCode: 500, Message: "A database error occurred"}
	ErrNoUpdateData = &ApiError{StatusCode: 400, Message: "No valid fields provided for update"}
)
