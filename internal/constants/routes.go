package constants

type Routes struct {
	APIPrefix string

	AuthBase            string
	UsersBase           string
	BusinessBase        string
	ProjectBase         string
	PublicationBase     string // <--- ADD THIS
	SkillToggleStatus   string // <--- ADD THIS
	SkillsBase          string // <--- ADD THIS
	SubscriptionBase    string // <--- ADD THIS
	UserConfigBase      string // <--- ADD THIS
	UserSkillsBase      string // <--- ADD THIS
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

	PublicationByID       string // <--- ADD THIS
	PublicationBySlug     string // <--- ADD THIS
	SubscriptionSubscribe string // <--- ADD THIS

	BusinessTags       string
	BusinessConnects   string
	ProjectMembers     string
	ProjectApplicants  string // <--- ADD THIS
	ProjectApply       string // <--- ADD THIS
	ProjectRegions     string // <--- ADD THIS
	ProjectSkills      string // <--- ADD THIS
	ProjectMemberships string // <--- ADD THIS

	DailyActEnrol string

	ConnectAccept string
	ConnectReject string

	UserEnrolments    string
	UserL2EResponses  string
	UserNotifications string
	UserApplications  string // <--- ADD THIS
	UserSubscriptions string // <--- ADD THIS

	UserNotifyReadAll      string
	ParamKeyID             string
	ParamKeyNotificationID string
	ParamKeyEntityType     string
	ParamKeyEntityID       string
	ParamKeyUserID         string
	ParamKeyRegionID       string // <--- ADD THIS
	ParamKeySkillID        string // <--- ADD THIS
	ParamKeySlug           string // <--- ADD THIS
	ParamKeySubscriptionID string // <--- ADD THIS
	ParamKeyConfigType     string // <--- ADD THIS

	ParamID                string
	ParamNotificationID    string
	ParamUserID            string
	ParamRegionID          string // <--- ADD THIS
	ParamSkillID           string // <--- ADD THIS
	ParamSubscriptionID    string // <--- ADD THIS
	ParamConfigType        string // <--- ADD THIS
	ParamSlug              string
	UserNotifyReadOne      string
	InferredBySource       string
	SkillToggleStatusRoute string // <--- ADD THIS (internal use for handler)

	UserSubscriptionCancel string // <--- ADD THIS (for DELETE /users/:id/subscriptions/:subID)
}

var AppRoutes = Routes{
	APIPrefix: "/api/v1",

	AuthBase:         "/auth",
	UsersBase:        "/users",
	BusinessBase:     "/businesses",
	ProjectBase:      "/projects",
	PublicationBase:  "/publications", // <--- ADD THIS
	TagsBase:         "/tags",
	SkillsBase:       "/skills",        // <--- ADD THIS (The standard REST path)
	SubscriptionBase: "/subscriptions", // <--- ADD THIS
	UserConfigBase:   "/:id/config",    // <--- ADD THIS
	UserSkillsBase:   "/:id/skills",    // <--- ADD THIS
	ConnectBase:      "/connections",
	EventBase:        "/events",
	FeedbackBase:     "/feedback",
	L2EBase:          "/l2e-responses",
	NotifyBase:       "/notifications",
	DailyActBase:     "/daily-activities",
	InferredBase:     "/inferred-connections",

	SkillToggleStatus: "/toggle-status", // <--- ADD THIS

	Login:                  "/login",
	SubscriptionSubscribe:  "/subscribe", // <--- ADD THIS
	Logout:                 "/logout",
	Me:                     "/me",
	PublicationByID:        "/id/:id",     // <--- ADD THIS
	PublicationBySlug:      "/slug/:slug", // <--- ADD THIS
	BusinessTags:           "/:id/tags",
	BusinessConnects:       "/:id/connections",
	ProjectMembers:         "/:id/members",
	ProjectApplicants:      "/:id/applicants", // <--- ADD THIS
	ProjectApply:           "/:id/apply",      // <--- ADD THIS
	ProjectRegions:         "/:id/regions",    // <--- ADD THIS
	ProjectSkills:          "/:id/skills",     // <--- ADD THIS
	DailyActEnrol:          "/:id/enrolments",
	ConnectAccept:          "/:id/accept",
	ConnectReject:          "/:id/reject",
	UserEnrolments:         "/:id/enrolments",
	UserL2EResponses:       "/:id/l2e-responses",
	UserNotifications:      "/:id/notifications",
	UserApplications:       "/:id/applications", // <--- ADD THIS
	UserSubscriptions:      "/:id/subscriptions",
	UserSubscriptionCancel: "/:id/subscriptions/:userSubscriptionID",
	ProjectMemberships:     "/:id/project-memberships", // <--- ADD THIS
	SkillToggleStatusRoute: "/:id/toggle-status",       // <--- ADD THIS (to be used in routes file)
	UserNotifyReadAll:      "/read-all",
	ParamKeyID:             "id",
	ParamKeyNotificationID: "notificationID",
	ParamKeyEntityType:     "entityType",
	ParamKeyEntityID:       "entityID",
	ParamKeyUserID:         "userID",
	ParamKeyRegionID:       "regionID",           // <--- ADD THIS
	ParamKeySkillID:        "skillID",            // <--- ADD THIS
	ParamKeySlug:           "slug",               // <--- ADD THIS
	ParamKeySubscriptionID: "userSubscriptionID", // <--- ADD THIS
	ParamKeyConfigType:     "configType",         // <--- ADD THIS
	ParamID:                "/:id",
	ParamNotificationID:    "/:notificationID",
	ParamUserID:            "/:userID",
	ParamRegionID:          "/:regionID",           // <--- ADD THIS
	ParamSkillID:           "/:skillID",            // <--- ADD THIS
	ParamSlug:              "/:slug",               // <--- ADD THIS
	ParamSubscriptionID:    "/:userSubscriptionID", // <-- MUST be the path segment for DELETE
	ParamConfigType:        "/:configType",         // <--- ADD THIS
	UserNotifyReadOne:      "/:notificationID/read",
	InferredBySource:       "/source/:entityType/:entityID",
	ContextKeyUser:         "user",
	ContextKeyUserID:       "userID",
	ContextKeySessionID:    "sessionID",
}
