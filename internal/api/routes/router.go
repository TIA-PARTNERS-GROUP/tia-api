package routes

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	AuthMiddleware                gin.HandlerFunc
	UserHandler                   *handlers.UserHandler
	AuthHandler                   *handlers.AuthHandler
	BusinessHandler               *handlers.BusinessHandler
	ProjectHandler                *handlers.ProjectHandler
	BusinessConnectionHandler     *handlers.BusinessConnectionHandler
	BusinessTagHandler            *handlers.BusinessTagHandler
	DailyActivityHandler          *handlers.DailyActivityHandler
	DailyActivityEnrolmentHandler *handlers.DailyActivityEnrolmentHandler
	EventHandler                  *handlers.EventHandler
	FeedbackHandler               *handlers.FeedbackHandler
	InferredConnectionHandler     *handlers.InferredConnectionHandler
	L2EHandler                    *handlers.L2EHandler
	NotificationHandler           *handlers.NotificationHandler
	ProjectApplicantHandler       *handlers.ProjectApplicantHandler 
	ProjectMemberHandler          *handlers.ProjectMemberHandler    
	ProjectRegionHandler          *handlers.ProjectRegionHandler    
	ProjectSkillHandler           *handlers.ProjectSkillHandler     
	PublicationHandler            *handlers.PublicationHandler      
	SkillHandler                  *handlers.SkillHandler            
	SubscriptionHandler           *handlers.SubscriptionHandler     
	UserSubscriptionHandler       *handlers.UserSubscriptionHandler 
	UserConfigHandler             *handlers.UserConfigHandler       
	UserSkillHandler              *handlers.UserSkillHandler        
	Routes                        constants.Routes
}

func RegisterRoutes(router *gin.Engine, deps *RouterDependencies) {
	router.RedirectTrailingSlash = false
	api := router.Group(deps.Routes.APIPrefix)
	SetupAuthRoutes(api, deps)
	SetupUserRoutes(api, deps)
	SetupBusinessRoutes(api, deps)
	SetupProjectRoutes(api, deps)
	SetupBusinessTagRoutes(api, deps)
	SetupConnectionRoutes(api, deps)
	SetupDailyActivityRoutes(api, deps)
	SetupEventRoutes(api, deps)
	SetupFeedbackRoutes(api, deps)
	SetupInferredConnectionRoutes(api, deps)
	SetupL2ERoutes(api, deps)
	SetupNotificationRoutes(api, deps)
	SetupPublicationRoutes(api, deps)  
	SetupSubscriptionRoutes(api, deps) 
	SetupSkillRoutes(api, deps)        
}
