package models

import (
	"time"

	"gorm.io/datatypes"
)

type BusinessType string
type BusinessCategory string
type BusinessPhase string
type BusinessConnectionType string
type BusinessConnectionStatus string
type PublicationType string
type NotificationType string
type RelatedEntityType string
type ProjectStatus string
type ProjectMemberRole string
type ProjectSkillImportance string
type UserSkillProficiency string
type IdeaStatus string
type IdeaVoteType string
type BusinessTagType string
type DailyActivityProgressStatus string

const (
	BusinessTypeConsulting       BusinessType                = "Consulting"
	BusinessTypeRetail           BusinessType                = "Retail"
	BusinessTypeTechnology       BusinessType                = "Technology"
	BusinessTypeManufacturing    BusinessType                = "Manufacturing"
	BusinessTypeServices         BusinessType                = "Services"
	BusinessTypeOther            BusinessType                = "Other"
	BusinessCategoryB2B          BusinessCategory            = "B2B"
	BusinessCategoryB2C          BusinessCategory            = "B2C"
	BusinessCategoryNonProfit    BusinessCategory            = "Non_Profit"
	BusinessCategoryGovernment   BusinessCategory            = "Government"
	BusinessCategoryMixed        BusinessCategory            = "Mixed"
	BusinessPhaseStartup         BusinessPhase               = "Startup"
	BusinessPhaseGrowth          BusinessPhase               = "Growth"
	BusinessPhaseMature          BusinessPhase               = "Mature"
	BusinessPhaseExit            BusinessPhase               = "Exit"
	ConnectionTypePartnership    BusinessConnectionType      = "Partnership"
	ConnectionTypeSupplier       BusinessConnectionType      = "Supplier"
	ConnectionTypeClient         BusinessConnectionType      = "Client"
	ConnectionTypeReferral       BusinessConnectionType      = "Referral"
	ConnectionTypeCollaboration  BusinessConnectionType      = "Collaboration"
	ConnectionStatusPending      BusinessConnectionStatus    = "pending"
	ConnectionStatusActive       BusinessConnectionStatus    = "active"
	ConnectionStatusRejected     BusinessConnectionStatus    = "rejected"
	ConnectionStatusInactive     BusinessConnectionStatus    = "inactive"
	ProjectStatusPlanning        ProjectStatus               = "planning"
	ProjectStatusActive          ProjectStatus               = "active"
	ProjectStatusOnHold          ProjectStatus               = "on_hold"
	ProjectStatusCompleted       ProjectStatus               = "completed"
	ProjectStatusCancelled       ProjectStatus               = "cancelled"
	ProjectMemberRoleManager     ProjectMemberRole           = "manager"
	ProjectMemberRoleContributor ProjectMemberRole           = "contributor"
	ProjectMemberRoleReviewer    ProjectMemberRole           = "reviewer"
	SkillImportanceRequired      ProjectSkillImportance      = "required"
	SkillImportancePreferred     ProjectSkillImportance      = "preferred"
	SkillImportanceOptional      ProjectSkillImportance      = "optional"
	ProficiencyBeginner          UserSkillProficiency        = "beginner"
	ProficiencyIntermediate      UserSkillProficiency        = "intermediate"
	ProficiencyAdvanced          UserSkillProficiency        = "advanced"
	ProficiencyExpert            UserSkillProficiency        = "expert"
	PublicationPost              PublicationType             = "post"
	PublicationCaseStudy         PublicationType             = "case_study"
	PublicationTestimonial       PublicationType             = "testimonial"
	PublicationArticle           PublicationType             = "article"
	IdeaStatusOpen               IdeaStatus                  = "open"
	IdeaStatusUnderReview        IdeaStatus                  = "under_review"
	IdeaStatusPlanned            IdeaStatus                  = "planned"
	IdeaStatusInProgress         IdeaStatus                  = "in_progress"
	IdeaStatusCompleted          IdeaStatus                  = "completed"
	IdeaStatusRejected           IdeaStatus                  = "rejected"
	IdeaVoteUp                   IdeaVoteType                = "up"
	IdeaVoteDown                 IdeaVoteType                = "down"
	BusinessTagClient            BusinessTagType             = "client"
	BusinessTagService           BusinessTagType             = "service"
	BusinessTagSpecialty         BusinessTagType             = "specialty"
	ProgressStatusNotStarted     DailyActivityProgressStatus = "not_started"
	ProgressStatusInProgress     DailyActivityProgressStatus = "in_progress"
	ProgressStatusCompleted      DailyActivityProgressStatus = "completed"
)

type User struct {
	ID                       uint    `gorm:"primaryKey"`
	FirstName                string  `gorm:"size:60;not null"`
	LastName                 *string `gorm:"size:60"`
	LoginEmail               string  `gorm:"size:254;not null;unique"`
	PasswordHash             *string `gorm:"size:255"`
	ContactEmail             *string `gorm:"size:254;index"`
	ContactPhoneNo           *string `gorm:"size:20"`
	AdkSessionID             *string `gorm:"size:128"`
	PasswordResetToken       []byte
	PasswordResetRequestedAt *time.Time
	EmailVerified            bool      `gorm:"default:false;not null"`
	Active                   bool      `gorm:"default:true;not null;index"`
	CreatedAt                time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt                time.Time `gorm:"not null;default:current_timestamp"`

	Businesses              []Business                  `gorm:"foreignKey:OperatorUserID"`
	ManagedProjects         []Project                   `gorm:"foreignKey:ManagedByUserID"`
	InitiatedConnections    []BusinessConnection        `gorm:"foreignKey:InitiatedByUserID"`
	ReceivedNotifications   []Notification              `gorm:"foreignKey:ReceiverUserID"`
	SentNotifications       []Notification              `gorm:"foreignKey:SenderUserID"`
	ProjectMemberships      []ProjectMember             `gorm:"foreignKey:UserID"`
	Publications            []Publication               `gorm:"foreignKey:UserID"`
	UserSessions            []UserSession               `gorm:"foreignKey:UserID"`
	UserSkills              []UserSkill                 `gorm:"foreignKey:UserID"`
	DailyActivityEnrolments []DailyActivityEnrolment    `gorm:"foreignKey:UserID"`
	ProjectApplicants       []ProjectApplicant          `gorm:"foreignKey:UserID"`
	UserSubscriptions       []UserSubscription          `gorm:"foreignKey:UserID"`
	UserConfigs             []UserConfig                `gorm:"foreignKey:UserID"`
	L2EResponses            []L2EResponse               `gorm:"foreignKey:UserID"`
	DailyActivityProgress   []UserDailyActivityProgress `gorm:"foreignKey:UserID"`
}
type Feedback struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `gorm:"size:120;not null"`
	Email         string    `gorm:"size:254;not null"`
	Content       string    `gorm:"type:text;not null"`
	DateSubmitted time.Time `gorm:"not null;default:current_timestamp"`
}
type ProjectApplicant struct {
	ProjectID uint `gorm:"primaryKey"`
	UserID    uint `gorm:"primaryKey"`

	Project Project `gorm:"foreignKey:ProjectID"`
	User    User    `gorm:"foreignKey:UserID"`
}
type DailyActivity struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:60;not null;unique"`
	Description string `gorm:"type:text;not null"`

	Enrolments []DailyActivityEnrolment `gorm:"foreignKey:DailyActivityID"`
}
type DailyActivityEnrolment struct {
	DailyActivityID uint `gorm:"primaryKey"`
	UserID          uint `gorm:"primaryKey"`

	DailyActivity DailyActivity `gorm:"foreignKey:DailyActivityID"`
	User          User          `gorm:"foreignKey:UserID"`
}
type UserDailyActivityProgress struct {
	UserID          uint `gorm:"primaryKey"`
	DailyActivityID uint `gorm:"primaryKey"`

	Date     datatypes.Date              `gorm:"primaryKey;type:date" swaggertype:"string"`
	Status   DailyActivityProgressStatus `gorm:"type:enum('not_started', 'in_progress', 'completed');default:not_started"`
	Progress int                         `gorm:"default:0"`

	User          User          `gorm:"foreignKey:UserID"`
	DailyActivity DailyActivity `gorm:"foreignKey:DailyActivityID"`
}
type Event struct {
	ID        uint   `gorm:"primaryKey"`
	EventType string `gorm:"size:100;not null;index"`

	Payload   datatypes.JSON `gorm:"type:json" swaggertype:"object"`
	Timestamp time.Time      `gorm:"not null;default:current_timestamp;index"`
	UserID    *uint          `gorm:"index"`

	User *User `gorm:"foreignKey:UserID"`
}
type Subscription struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:100;not null;unique"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	ValidDays   *int
	ValidMonths *int
}
type UserSubscription struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null"`
	SubscriptionID uint      `gorm:"not null"`
	DateFrom       time.Time `gorm:"not null"`
	DateTo         time.Time `gorm:"not null"`
	IsTrial        bool      `gorm:"default:false"`

	User         User         `gorm:"foreignKey:UserID"`
	Subscription Subscription `gorm:"foreignKey:SubscriptionID"`
}
type UserConfig struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null;uniqueIndex:uq_user_config_type"`
	ConfigType string `gorm:"size:50;not null;uniqueIndex:uq_user_config_type"`

	Config datatypes.JSON `gorm:"type:json;not null" swaggertype:"object"`

	User User `gorm:"foreignKey:UserID"`
}
type L2EResponse struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`

	Response  datatypes.JSON `gorm:"type:json;not null" swaggertype:"object"`
	DateAdded time.Time      `gorm:"not null;default:current_timestamp"`

	User User `gorm:"foreignKey:UserID"`
}
type Business struct {
	ID               uint             `gorm:"primaryKey"`
	OperatorUserID   uint             `gorm:"not null;index"`
	Name             string           `gorm:"size:100;not null"`
	Tagline          *string          `gorm:"size:100"`
	Website          *string          `gorm:"size:255"`
	ContactName      *string          `gorm:"size:60"`
	ContactPhoneNo   *string          `gorm:"size:20"`
	ContactEmail     *string          `gorm:"size:254"`
	Description      *string          `gorm:"type:text"`
	Address          *string          `gorm:"size:100"`
	City             *string          `gorm:"size:60"`
	State            *string          `gorm:"size:60"`
	Country          *string          `gorm:"size:60"`
	PostalCode       *string          `gorm:"size:20"`
	Value            *float64         `gorm:"type:decimal(15,2)"`
	BusinessType     BusinessType     `gorm:"type:enum('Consulting', 'Retail', 'Technology', 'Manufacturing', 'Services', 'Other');index"`
	BusinessCategory BusinessCategory `gorm:"type:enum('B2B', 'B2C', 'Non_Profit', 'Government', 'Mixed')"`
	BusinessPhase    BusinessPhase    `gorm:"type:enum('Startup', 'Growth', 'Mature', 'Exit')"`
	Active           bool             `gorm:"default:true;not null;index"`
	CreatedAt        time.Time        `gorm:"not null;default:current_timestamp"`
	UpdatedAt        time.Time        `gorm:"not null;default:current_timestamp"`

	OperatorUser          User                 `gorm:"foreignKey:OperatorUserID"`
	BusinessTags          []BusinessTag        `gorm:"foreignKey:BusinessID"`
	Projects              []Project            `gorm:"foreignKey:BusinessID"`
	Publications          []Publication        `gorm:"foreignKey:BusinessID"`
	InitiatingConnections []BusinessConnection `gorm:"foreignKey:InitiatingBusinessID"`
	ReceivingConnections  []BusinessConnection `gorm:"foreignKey:ReceivingBusinessID"`
}
type Project struct {
	ID              uint          `gorm:"primaryKey"`
	ManagedByUserID uint          `gorm:"not null;index"`
	BusinessID      *uint         `gorm:"index"`
	Name            string        `gorm:"size:100;not null"`
	Description     *string       `gorm:"type:text"`
	ProjectStatus   ProjectStatus `gorm:"type:enum('planning', 'active', 'on_hold', 'completed', 'cancelled');default:planning;index"`
	StartDate       *time.Time
	TargetEndDate   *time.Time
	ActualEndDate   *time.Time
	CreatedAt       time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt       time.Time `gorm:"not null;default:current_timestamp"`

	ManagingUser   User            `gorm:"foreignKey:ManagedByUserID"`
	Business       *Business       `gorm:"foreignKey:BusinessID"`
	ProjectMembers []ProjectMember `gorm:"foreignKey:ProjectID"`
	ProjectSkills  []ProjectSkill  `gorm:"foreignKey:ProjectID"`
	ProjectRegions []ProjectRegion `gorm:"foreignKey:ProjectID"`
}
type Region struct {
	ID   string `gorm:"primaryKey;size:3"`
	Name string `gorm:"size:50;not null;unique"`
}
type ProjectRegion struct {
	RegionID  string `gorm:"primaryKey;size:3"`
	ProjectID uint   `gorm:"primaryKey"`

	Region  Region  `gorm:"foreignKey:RegionID"`
	Project Project `gorm:"foreignKey:ProjectID"`
}
type InferredConnection struct {
	ID               uint      `gorm:"primaryKey"`
	SourceEntityType string    `gorm:"size:50;not null;index"`
	SourceEntityID   uint      `gorm:"not null;index"`
	TargetEntityType string    `gorm:"size:50;not null;index"`
	TargetEntityID   uint      `gorm:"not null;index"`
	ConnectionType   string    `gorm:"size:100;not null;index"`
	ConfidenceScore  float64   `gorm:"not null"`
	ModelVersion     string    `gorm:"size:50"`
	CreatedAt        time.Time `gorm:"not null;default:current_timestamp"`
}
type Skill struct {
	ID          uint      `gorm:"primaryKey"`
	Category    string    `gorm:"size:100;not null;index"`
	Name        string    `gorm:"size:100;not null;unique"`
	Description *string   `gorm:"type:text"`
	Active      bool      `gorm:"default:true;not null;index"`
	CreatedAt   time.Time `gorm:"not null;default:current_timestamp"`

	ProjectSkills []ProjectSkill `gorm:"foreignKey:SkillID"`
	UserSkills    []UserSkill    `gorm:"foreignKey:SkillID"`
}
type Publication struct {
	ID              uint            `gorm:"primaryKey"`
	UserID          uint            `gorm:"not null;index"`
	BusinessID      *uint           `gorm:"index"`
	PublicationType PublicationType `gorm:"type:enum('post', 'case_study', 'testimonial', 'article');index"`
	Title           string          `gorm:"size:255;not null"`
	Slug            string          `gorm:"size:300;not null;unique"`
	Excerpt         *string         `gorm:"type:text"`
	Content         string          `gorm:"type:longtext;not null"`
	Thumbnail       *string         `gorm:"size:255"`
	VideoURL        *string         `gorm:"size:255"`
	Published       bool            `gorm:"default:false;not null;index"`
	PublishedAt     *time.Time
	CreatedAt       time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt       time.Time `gorm:"not null;default:current_timestamp"`

	User     User      `gorm:"foreignKey:UserID"`
	Business *Business `gorm:"foreignKey:BusinessID"`
}
type Notification struct {
	ID                uint               `gorm:"primaryKey"`
	SenderUserID      *uint              `gorm:"index"`
	ReceiverUserID    uint               `gorm:"not null;index"`
	NotificationType  NotificationType   `gorm:"type:enum('connection_request', 'project_invite', 'message', 'system')"`
	Title             string             `gorm:"size:255;not null"`
	Message           string             `gorm:"type:text;not null"`
	RelatedEntityType *RelatedEntityType `gorm:"type:enum('business', 'project', 'publication', 'idea')"`
	RelatedEntityID   *uint
	Read              bool      `gorm:"column:read;default:false;not null;index"`
	ActionURL         *string   `gorm:"size:500"`
	CreatedAt         time.Time `gorm:"not null;default:current_timestamp;index"`

	ReceiverUser User  `gorm:"foreignKey:ReceiverUserID"`
	SenderUser   *User `gorm:"foreignKey:SenderUserID"`
}
type UserSkill struct {
	SkillID          uint                 `gorm:"primaryKey"`
	UserID           uint                 `gorm:"primaryKey;index"`
	ProficiencyLevel UserSkillProficiency `gorm:"type:enum('beginner', 'intermediate', 'advanced', 'expert');default:intermediate"`
	CreatedAt        time.Time            `gorm:"not null;default:current_timestamp"`

	Skill Skill `gorm:"foreignKey:SkillID"`
	User  User  `gorm:"foreignKey:UserID"`
}
type ProjectSkill struct {
	ProjectID  uint                   `gorm:"primaryKey"`
	SkillID    uint                   `gorm:"primaryKey;index"`
	Importance ProjectSkillImportance `gorm:"type:enum('required', 'preferred', 'optional');default:preferred"`

	Project Project `gorm:"foreignKey:ProjectID"`
	Skill   Skill   `gorm:"foreignKey:SkillID"`
}
type ProjectMember struct {
	ProjectID uint              `gorm:"primaryKey"`
	UserID    uint              `gorm:"primaryKey;index"`
	Role      ProjectMemberRole `gorm:"type:enum('manager', 'contributor', 'reviewer');default:contributor"`
	JoinedAt  time.Time         `gorm:"not null;default:current_timestamp"`

	Project Project `gorm:"foreignKey:ProjectID"`
	User    User    `gorm:"foreignKey:UserID"`
}
type BusinessConnection struct {
	ID                   uint                     `gorm:"primaryKey"`
	InitiatingBusinessID uint                     `gorm:"not null;uniqueIndex:uq_business_connections_unique;index"`
	ReceivingBusinessID  uint                     `gorm:"not null;uniqueIndex:uq_business_connections_unique;index"`
	ConnectionType       BusinessConnectionType   `gorm:"type:enum('Partnership', 'Supplier', 'Client', 'Referral', 'Collaboration');uniqueIndex:uq_business_connections_unique"`
	Status               BusinessConnectionStatus `gorm:"type:enum('pending', 'active', 'rejected', 'inactive');default:pending;index"`
	InitiatedByUserID    uint                     `gorm:"not null;index"`
	Notes                *string                  `gorm:"type:text"`
	CreatedAt            time.Time                `gorm:"not null;default:current_timestamp"`
	UpdatedAt            time.Time                `gorm:"not null;default:current_timestamp"`

	InitiatingBusiness Business `gorm:"foreignKey:InitiatingBusinessID"`
	ReceivingBusiness  Business `gorm:"foreignKey:ReceivingBusinessID"`
	InitiatedByUser    User     `gorm:"foreignKey:InitiatedByUserID"`
}
type BusinessTag struct {
	ID          uint            `gorm:"primaryKey"`
	BusinessID  uint            `gorm:"not null;index"`
	TagType     BusinessTagType `gorm:"type:enum('client', 'service', 'specialty');index"`
	Description string          `gorm:"size:100;not null"`
	CreatedAt   time.Time       `gorm:"not null;default:current_timestamp"`

	Business Business `gorm:"foreignKey:BusinessID"`
}
type UserSession struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"size:128;not null;unique"`
	IPAddress *string   `gorm:"size:45"`
	UserAgent *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	ExpiresAt time.Time `gorm:"not null;index"`
	RevokedAt *time.Time

	User User `gorm:"foreignKey:UserID"`
}
