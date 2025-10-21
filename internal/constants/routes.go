package constants

// This struct holds all route paths, grouped by resource.
// This becomes your single source of truth.
type Routes struct {
	APIPrefix string // /api/v1

	// Resource Base Paths
	AuthBase     string // /auth
	UsersBase    string // /users
	BusinessBase string // /businesses
	TagsBase     string // /tags
	ConnectBase  string // /connections
	EventBase    string // /events
	FeedbackBase string // /feedback
	L2EBase      string // /l2e-responses
	NotifyBase   string // /notifications
	DailyActBase string // /daily-activities
	InferredBase string // /inferred-connections

	ContextKeyUser      string // "user"
	ContextKeyUserID    string // "userID"
	ContextKeySessionID string // "sessionID"

	// Auth Sub-routes
	Login  string // /login
	Logout string // /logout
	Me     string // /me

	// Business Sub-routes
	BusinessTags     string // /:id/tags
	BusinessConnects string // /:id/connections

	// Daily Activity Sub-routes
	DailyActEnrol string // /:id/enrolments

	// Connection Sub-routes
	ConnectAccept string // /:id/accept
	ConnectReject string // /:id/reject

	// Add to your 'type Routes struct'
	UserEnrolments    string // /:id/enrolments
	UserL2EResponses  string // /:id/l2e-responses
	UserNotifications string // /:id/notifications

	// Notification Sub-sub-routes
	UserNotifyReadAll string // /read-all

	ParamKeyID             string // "id"
	ParamKeyNotificationID string // "notificationID"
	ParamKeyEntityType     string // "entityType"
	ParamKeyEntityID       string // "entityID"

	// --- Route Segments (for routes files) ---
	// These are your existing param fields
	ParamID             string // /:id
	ParamNotificationID string // /:notificationID
	UserNotifyReadOne   string // /:notificationID/read
	InferredBySource    string // /source/:entityType/:entityID
}

// AppRoutes defines all the routes used in the application.
var AppRoutes = Routes{
	APIPrefix: "/api/v1",

	// Bases
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

	// Sub-routes
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
