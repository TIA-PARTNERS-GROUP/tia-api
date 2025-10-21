package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupSkillRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	// --- MODIFIED: Use the correct SkillsBase constant ---
	skills := api.Group(deps.Routes.SkillsBase)
	skills.Use(deps.AuthMiddleware)
	{
		skills.POST("", deps.SkillHandler.CreateSkill)
		skills.GET("", deps.SkillHandler.GetSkills)
		skills.GET(deps.Routes.ParamID, deps.SkillHandler.GetSkillByID)
		skills.PUT(deps.Routes.ParamID, deps.SkillHandler.UpdateSkill)
		skills.DELETE(deps.Routes.ParamID, deps.SkillHandler.DeleteSkill)
		skills.PATCH(deps.Routes.SkillToggleStatusRoute, deps.SkillHandler.ToggleSkillStatus)
	}
}
