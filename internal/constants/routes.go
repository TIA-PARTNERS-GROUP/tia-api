package constants

type Routes struct {
	APIPrefix string

	AuthBase            string
	UsersBase           string
	BusinessBase        string
	ProjectBase         string
	PublicationBase     string 
	SkillToggleStatus   string 
	SkillsBase          string 
	SubscriptionBase    string 
	UserConfigBase      string 
	UserSkillsBase      string 
	TagsBase            string
	ConnectBase         string
	EventBase           string
	FeedbackBase        string
	L2EBase             string
	NotifyBase          string
	DailyActBase        string
	InferredBase        string
	ContextKeyUser      string
	ContextKeyUserID    string
	ContextKeySessionID string

	Login  string
	Logout string
	Me     string

	PublicationByID       string 
	PublicationBySlug     string 
	SubscriptionSubscribe string 

	BusinessTags       string
	BusinessConnects   string
	ProjectMembers     string
	ProjectApplicants  string 
	ProjectApply       string 
	ProjectRegions     string 
	ProjectSkills      string 
	ProjectMemberships string 

	DailyActEnrol string

	ConnectAccept string
	ConnectReject string

	UserEnrolments    string
	UserL2EResponses  string
	UserNotifications string
	UserApplications  string 
	UserSubscriptions string 

	UserNotifyReadAll      string
	ParamKeyID             string
	ParamKeyNotificationID string
	ParamKeyEntityType     string
	ParamKeyEntityID       string
	ParamKeyUserID         string
	ParamKeyRegionID       string 
	ParamKeySkillID        string 
	ParamKeySlug           string 
	ParamKeySubscriptionID string 
	ParamKeyConfigType     string 

	ParamID                string
	ParamNotificationID    string
	ParamUserID            string
	ParamRegionID          string 
	ParamSkillID           string 
	ParamSubscriptionID    string 
	ParamConfigType        string 
	ParamSlug              string
	UserNotifyReadOne      string
	InferredBySource       string
	SkillToggleStatusRoute string 

	UserSubscriptionCancel string 
}

var AppRoutes = Routes{
	APIPrefix: "/api/v1",

	AuthBase:         "/auth",
	UsersBase:        "/users",
	BusinessBase:     "/businesses",
	ProjectBase:      "/projects",
	PublicationBase:  "/publications", 
	TagsBase:         "/tags",
	SkillsBase:       "/skills",        
	SubscriptionBase: "/subscriptions", 
	UserConfigBase:   "/:id/config",    
	UserSkillsBase:   "/:id/skills",    
	ConnectBase:      "/connections",
	EventBase:        "/events",
	FeedbackBase:     "/feedback",
	L2EBase:          "/l2e-responses",
	NotifyBase:       "/notifications",
	DailyActBase:     "/daily-activities",
	InferredBase:     "/inferred-connections",

	SkillToggleStatus: "/toggle-status", 

	Login:                  "/login",
	SubscriptionSubscribe:  "/subscribe", 
	Logout:                 "/logout",
	Me:                     "/me",
	PublicationByID:        "/id/:id",     
	PublicationBySlug:      "/slug/:slug", 
	BusinessTags:           "/:id/tags",
	BusinessConnects:       "/:id/connections",
	ProjectMembers:         "/:id/members",
	ProjectApplicants:      "/:id/applicants", 
	ProjectApply:           "/:id/apply",      
	ProjectRegions:         "/:id/regions",    
	ProjectSkills:          "/:id/skills",     
	DailyActEnrol:          "/:id/enrolments",
	ConnectAccept:          "/:id/accept",
	ConnectReject:          "/:id/reject",
	UserEnrolments:         "/:id/enrolments",
	UserL2EResponses:       "/:id/l2e-responses",
	UserNotifications:      "/:id/notifications",
	UserApplications:       "/:id/applications", 
	UserSubscriptions:      "/:id/subscriptions",
	UserSubscriptionCancel: "/:id/subscriptions/:userSubscriptionID",
	ProjectMemberships:     "/:id/project-memberships", 
	SkillToggleStatusRoute: "/:id/toggle-status",       
	UserNotifyReadAll:      "/read-all",
	ParamKeyID:             "id",
	ParamKeyNotificationID: "notificationID",
	ParamKeyEntityType:     "entityType",
	ParamKeyEntityID:       "entityID",
	ParamKeyUserID:         "userID",
	ParamKeyRegionID:       "regionID",           
	ParamKeySkillID:        "skillID",            
	ParamKeySlug:           "slug",               
	ParamKeySubscriptionID: "userSubscriptionID", 
	ParamKeyConfigType:     "configType",         
	ParamID:                "/:id",
	ParamNotificationID:    "/:notificationID",
	ParamUserID:            "/:userID",
	ParamRegionID:          "/:regionID",           
	ParamSkillID:           "/:skillID",            
	ParamSlug:              "/:slug",               
	ParamSubscriptionID:    "/:userSubscriptionID", 
	ParamConfigType:        "/:configType",         
	UserNotifyReadOne:      "/:notificationID/read",
	InferredBySource:       "/source/:entityType/:entityID",
	ContextKeyUser:         "user",
	ContextKeyUserID:       "userID",
	ContextKeySessionID:    "sessionID",
}
