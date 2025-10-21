package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupProjectRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	projects := api.Group(deps.Routes.ProjectBase)
	projects.Use(deps.AuthMiddleware)
	{
		projects.POST("", deps.ProjectHandler.CreateProject)
		projects.GET("", deps.ProjectHandler.GetAllProjects)
		projects.GET(deps.Routes.ParamID, deps.ProjectHandler.GetProjectByID)
		projects.PUT(deps.Routes.ParamID, deps.ProjectHandler.UpdateProject)
		projects.DELETE(deps.Routes.ParamID, deps.ProjectHandler.DeleteProject)

		// --- Project Applicant Routes ---
		projects.POST(deps.Routes.ProjectApply, deps.ProjectApplicantHandler.ApplyToProject)
		projects.DELETE(deps.Routes.ProjectApply, deps.ProjectApplicantHandler.WithdrawApplication)
		projects.GET(deps.Routes.ProjectApplicants, deps.ProjectApplicantHandler.GetApplicantsForProject)

		// --- Project Region Routes --- // <--- ADD THIS BLOCK
		regions := projects.Group(deps.Routes.ProjectRegions)
		{
			regions.POST("", deps.ProjectRegionHandler.AddRegionToProject)
			regions.GET("", deps.ProjectRegionHandler.GetRegionsForProject)
			regions.DELETE(deps.Routes.ParamRegionID, deps.ProjectRegionHandler.RemoveRegionFromProject)
		}

		// --- Project Skill Routes --- // <--- ADD THIS BLOCK
		skills := projects.Group(deps.Routes.ProjectSkills)
		{
			skills.POST("", deps.ProjectSkillHandler.AddProjectSkill)
			skills.GET("", deps.ProjectSkillHandler.GetProjectSkills)
			skills.PUT(deps.Routes.ParamSkillID, deps.ProjectSkillHandler.UpdateProjectSkill)
			skills.DELETE(deps.Routes.ParamSkillID, deps.ProjectSkillHandler.RemoveProjectSkill)
		}

		// Nested routes for project members
		members := projects.Group(deps.Routes.ProjectMembers)
		{
			members.POST("", deps.ProjectMemberHandler.AddProjectMember)
			members.GET("", deps.ProjectMemberHandler.GetProjectMembers)
			members.GET(deps.Routes.ParamUserID, deps.ProjectMemberHandler.GetProjectMember)
			members.PUT(deps.Routes.ParamUserID, deps.ProjectMemberHandler.UpdateProjectMemberRole)
			members.DELETE(deps.Routes.ParamUserID, deps.ProjectMemberHandler.RemoveProjectMember)
		}
	}
}
