package constants
type Routes struct {
	APIPrefix string 
	
	AuthBase     string 
	UsersBase    string 
	BusinessBase string 
	TagsBase     string 
	ConnectBase  string 
	EventBase    string 
	FeedbackBase string 
	L2EBase      string 
	NotifyBase   string 
	DailyActBase string 
	InferredBase string 
	ContextKeyUser      string 
	ContextKeyUserID    string 
	ContextKeySessionID string 
	
	Login  string 
	Logout string 
	Me     string 
	
	BusinessTags     string 
	BusinessConnects string 
	
	DailyActEnrol string 
	
	ConnectAccept string 
	ConnectReject string 
	
	UserEnrolments    string 
	UserL2EResponses  string 
	UserNotifications string 
	
	UserNotifyReadAll string 
	ParamKeyID             string 
	ParamKeyNotificationID string 
	ParamKeyEntityType     string 
	ParamKeyEntityID       string 
	
	
	ParamID             string 
	ParamNotificationID string 
	UserNotifyReadOne   string 
	InferredBySource    string 
}
var AppRoutes = Routes{
	APIPrefix: "/api/v1",
	
	AuthBase:     "/auth",
	UsersBase:    "/users",
	BusinessBase: "/businesses",
	TagsBase:     "/tags",
	ConnectBase:  "/connections",
	EventBase:    "/events",
	FeedbackBase: "/feedback",
	L2EBase:      "/l2e-responses",
	NotifyBase:   "/notifications",
	DailyActBase: "/daily-activities",
	InferredBase: "/inferred-connections",
	
	Login:  "/login",
	Logout: "/logout",
	Me:     "/me",
	BusinessTags:     "/:id/tags",
	BusinessConnects: "/:id/connections",
	DailyActEnrol: "/:id/enrolments",
	ConnectAccept: "/:id/accept",
	ConnectReject: "/:id/reject",
	UserEnrolments:    "/:id/enrolments",
	UserL2EResponses:  "/:id/l2e-responses",
	UserNotifications: "/:id/notifications",
	UserNotifyReadAll: "/read-all",
	ParamKeyID:             "id",
	ParamKeyNotificationID: "notificationID",
	ParamKeyEntityType:     "entityType",
	ParamKeyEntityID:       "entityID",
	ParamID:             "/:id",
	ParamNotificationID: "/:notificationID",
	UserNotifyReadOne:   "/:notificationID/read",
	InferredBySource:    "/source/:entityType/:entityID",
	ContextKeyUser:      "user",
	ContextKeyUserID:    "userID",
	ContextKeySessionID: "sessionID",
}
