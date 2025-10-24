package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/configs"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataGenerator struct {
	db *gorm.DB
}

func main() {
	config := configs.LoadConfig()
	
	db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.Business{},
		&models.Project{},
		&models.Skill{},
		&models.Publication{},
		&models.Notification{},
		&models.UserSkill{},
		&models.ProjectSkill{},
		&models.ProjectMember{},
		&models.BusinessConnection{},
		&models.BusinessTag{},
		&models.Feedback{},
		&models.ProjectApplicant{},
		&models.DailyActivity{},
		&models.DailyActivityEnrolment{},
		&models.UserDailyActivityProgress{},
		&models.Event{},
		&models.Subscription{},
		&models.UserSubscription{},
		&models.UserConfig{},
		&models.L2EResponse{},
		&models.Region{},
		&models.ProjectRegion{},
		&models.InferredConnection{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	generator := &DataGenerator{db: db}
	
	fmt.Println("Starting data generation...")
	
	regions := generator.generateRegions()
	skills := generator.generateSkills()
	users := generator.generateUsers()
	subscriptions := generator.generateSubscriptions()
	generator.generateUserSubscriptions(users, subscriptions)
	businesses := generator.generateBusinesses(users)
	projects := generator.generateProjects(users, businesses)
	generator.generateProjectMembers(users, projects)
	generator.generateProjectSkills(projects, skills)
	generator.generateProjectRegions(projects, regions)
	generator.generateUserSkills(users, skills)
	generator.generateBusinessTags(businesses)
	generator.generateBusinessConnections(businesses, users)
	publications := generator.generatePublications(users, businesses)
	notifications := generator.generateNotifications(users)
	generator.generateFeedback()
	dailyActivities := generator.generateDailyActivities()
	generator.generateDailyActivityEnrolments(users, dailyActivities)
	generator.generateUserDailyActivityProgress(users, dailyActivities)
	generator.generateEvents(users)
	generator.generateL2EResponses(users)
	generator.generateUserConfigs(users)
	generator.generateProjectApplicants(users, projects)
	generator.generateInferredConnections(businesses, projects, users)
	
	fmt.Printf("Data generation completed successfully!\n")
	fmt.Printf("Generated data summary:\n")
	fmt.Printf("   - Users: %d\n", len(users))
	fmt.Printf("   - Businesses: %d\n", len(businesses))
	fmt.Printf("   - Projects: %d\n", len(projects))
	fmt.Printf("   - Skills: %d\n", len(skills))
	fmt.Printf("   - Publications: %d\n", len(publications))
	fmt.Printf("   - Notifications: %d\n", len(notifications))
	fmt.Printf("   - Daily Activities: %d\n", len(dailyActivities))
	fmt.Printf("   - Regions: %d\n", len(regions))
	fmt.Printf("   - Subscriptions: %d\n", len(subscriptions))
}

func (dg *DataGenerator) generateRegions() []models.Region {
	regions := []models.Region{
		{ID: "US", Name: "United States"},
		{ID: "CA", Name: "Canada"},
		{ID: "UK", Name: "United Kingdom"},
		{ID: "DE", Name: "Germany"},
		{ID: "FR", Name: "France"},
		{ID: "AU", Name: "Australia"},
		{ID: "JP", Name: "Japan"},
		{ID: "SG", Name: "Singapore"},
		{ID: "IN", Name: "India"},
		{ID: "BR", Name: "Brazil"},
	}
	
	for i := range regions {
		dg.db.Create(&regions[i])
	}
	
	fmt.Println("Generated regions")
	return regions
}

func (dg *DataGenerator) generateSkills() []models.Skill {
	skills := []models.Skill{
		{Category: "Technology", Name: "Go Programming", Description: stringPtr("Backend development with Go"), Active: true},
		{Category: "Technology", Name: "Python", Description: stringPtr("Python programming language"), Active: true},
		{Category: "Technology", Name: "JavaScript", Description: stringPtr("Frontend and backend JavaScript development"), Active: true},
		{Category: "Technology", Name: "React", Description: stringPtr("React.js frontend framework"), Active: true},
		{Category: "Technology", Name: "Node.js", Description: stringPtr("Server-side JavaScript runtime"), Active: true},
		{Category: "Technology", Name: "Docker", Description: stringPtr("Containerization technology"), Active: true},
		{Category: "Technology", Name: "Kubernetes", Description: stringPtr("Container orchestration"), Active: true},
		{Category: "Technology", Name: "AWS", Description: stringPtr("Amazon Web Services cloud platform"), Active: true},
		{Category: "Technology", Name: "PostgreSQL", Description: stringPtr("Relational database management"), Active: true},
		{Category: "Technology", Name: "Redis", Description: stringPtr("In-memory data structure store"), Active: true},
		
		{Category: "Business", Name: "Project Management", Description: stringPtr("Managing projects and teams"), Active: true},
		{Category: "Business", Name: "Strategic Planning", Description: stringPtr("Long-term business strategy development"), Active: true},
		{Category: "Business", Name: "Financial Analysis", Description: stringPtr("Analyzing financial data and trends"), Active: true},
		{Category: "Business", Name: "Marketing Strategy", Description: stringPtr("Developing marketing campaigns and strategies"), Active: true},
		{Category: "Business", Name: "Sales Management", Description: stringPtr("Leading and managing sales teams"), Active: true},
		{Category: "Business", Name: "Operations Management", Description: stringPtr("Optimizing business operations"), Active: true},
		
		{Category: "Design", Name: "UI/UX Design", Description: stringPtr("User interface and experience design"), Active: true},
		{Category: "Design", Name: "Graphic Design", Description: stringPtr("Visual communication and design"), Active: true},
		{Category: "Design", Name: "Product Design", Description: stringPtr("Designing physical and digital products"), Active: true},
		{Category: "Design", Name: "Brand Design", Description: stringPtr("Creating and maintaining brand identity"), Active: true},
		
		{Category: "Data & Analytics", Name: "Data Analysis", Description: stringPtr("Analyzing and interpreting data"), Active: true},
		{Category: "Data & Analytics", Name: "Machine Learning", Description: stringPtr("Building ML models and algorithms"), Active: true},
		{Category: "Data & Analytics", Name: "Business Intelligence", Description: stringPtr("Transforming data into business insights"), Active: true},
		{Category: "Data & Analytics", Name: "Statistics", Description: stringPtr("Statistical analysis and modeling"), Active: true},
	}
	
	for i := range skills {
		dg.db.Create(&skills[i])
	}
	
	fmt.Println("Generated skills")
	return skills
}

func (dg *DataGenerator) generateUsers() []models.User {
	users := []models.User{
		{
			FirstName:    "John",
			LastName:     stringPtr("Smith"),
			LoginEmail:   "john.smith@example.com",
			ContactEmail: stringPtr("john.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0101"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Sarah",
			LastName:     stringPtr("Johnson"),
			LoginEmail:   "sarah.johnson@example.com",
			ContactEmail: stringPtr("sarah.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0102"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Michael",
			LastName:     stringPtr("Brown"),
			LoginEmail:   "michael.brown@example.com",
			ContactEmail: stringPtr("mike.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0103"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Emily",
			LastName:     stringPtr("Davis"),
			LoginEmail:   "emily.davis@example.com",
			ContactEmail: stringPtr("emily.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0104"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "David",
			LastName:     stringPtr("Wilson"),
			LoginEmail:   "david.wilson@example.com",
			ContactEmail: stringPtr("david.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0105"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Lisa",
			LastName:     stringPtr("Anderson"),
			LoginEmail:   "lisa.anderson@example.com",
			ContactEmail: stringPtr("lisa.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0106"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Robert",
			LastName:     stringPtr("Taylor"),
			LoginEmail:   "robert.taylor@example.com",
			ContactEmail: stringPtr("rob.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0107"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Jennifer",
			LastName:     stringPtr("Thomas"),
			LoginEmail:   "jennifer.thomas@example.com",
			ContactEmail: stringPtr("jen.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0108"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Christopher",
			LastName:     stringPtr("Jackson"),
			LoginEmail:   "chris.jackson@example.com",
			ContactEmail: stringPtr("chris.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0109"),
			EmailVerified: true,
			Active:       true,
		},
		{
			FirstName:    "Amanda",
			LastName:     stringPtr("White"),
			LoginEmail:   "amanda.white@example.com",
			ContactEmail: stringPtr("amanda.contact@example.com"),
			ContactPhoneNo: stringPtr("+1-555-0110"),
			EmailVerified: true,
			Active:       true,
		},
	}
	
	for i := range users {
		hashedPassword, _ := utils.HashPassword("Password123!")
		users[i].PasswordHash = &hashedPassword
		dg.db.Create(&users[i])
	}
	
	fmt.Println("Generated users")
	return users
}

func (dg *DataGenerator) generateSubscriptions() []models.Subscription {
	subscriptions := []models.Subscription{
		{Name: "Basic Plan", Price: 9.99, ValidMonths: intPtr(1)},
		{Name: "Professional Plan", Price: 29.99, ValidMonths: intPtr(1)},
		{Name: "Enterprise Plan", Price: 99.99, ValidMonths: intPtr(1)},
		{Name: "Trial Plan", Price: 0.00, ValidDays: intPtr(14)},
	}
	
	for i := range subscriptions {
		dg.db.Create(&subscriptions[i])
	}
	
	fmt.Println("Generated subscriptions")
	return subscriptions
}

func (dg *DataGenerator) generateUserSubscriptions(users []models.User, subscriptions []models.Subscription) []models.UserSubscription {
	var userSubscriptions []models.UserSubscription
	
	for i, user := range users {
		subscription := subscriptions[i%len(subscriptions)]
		now := time.Now()
		
		userSub := models.UserSubscription{
			UserID:         user.ID,
			SubscriptionID: subscription.ID,
			DateFrom:       now,
			DateTo:         now.AddDate(0, 1, 0),
			IsTrial:        subscription.Price == 0.00,
		}
		
		dg.db.Create(&userSub)
		userSubscriptions = append(userSubscriptions, userSub)
	}
	
	fmt.Println("Generated user subscriptions")
	return userSubscriptions
}

func (dg *DataGenerator) generateBusinesses(users []models.User) []models.Business {
	businesses := []models.Business{
		{
			OperatorUserID:   users[0].ID,
			Name:             "TechCorp Solutions",
			Tagline:          stringPtr("Innovative technology solutions for modern businesses"),
			Website:          stringPtr("https://techcorp.example.com"),
			ContactName:      stringPtr("John Smith"),
			ContactPhoneNo:   stringPtr("+1-555-0101"),
			ContactEmail:     stringPtr("contact@techcorp.example.com"),
			Description:      stringPtr("A leading technology company specializing in cloud solutions and digital transformation."),
			Address:          stringPtr("123 Tech Street"),
			City:             stringPtr("San Francisco"),
			State:            stringPtr("California"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("94105"),
			Value:            float64Ptr(5000000.00),
			BusinessType:     models.BusinessTypeTechnology,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseGrowth,
			Active:           true,
		},
		{
			OperatorUserID:   users[1].ID,
			Name:             "ConsumerTech Apps",
			Tagline:          stringPtr("Mobile apps that simplify everyday life"),
			Website:          stringPtr("https://consumertech.example.com"),
			ContactName:      stringPtr("Sarah Johnson"),
			ContactPhoneNo:   stringPtr("+1-555-0102"),
			ContactEmail:     stringPtr("hello@consumertech.example.com"),
			Description:      stringPtr("Consumer-focused mobile application development company."),
			Address:          stringPtr("456 App Street"),
			City:             stringPtr("New York"),
			State:            stringPtr("New York"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("10001"),
			Value:            float64Ptr(1200000.00),
			BusinessType:     models.BusinessTypeTechnology,
			BusinessCategory: models.BusinessCategoryB2C,
			BusinessPhase:    models.BusinessPhaseGrowth,
			Active:           true,
		},
		{
			OperatorUserID:   users[2].ID,
			Name:             "Digital Marketing Pro",
			Tagline:          stringPtr("Data-driven marketing solutions"),
			Website:          stringPtr("https://digitalmarketing.example.com"),
			ContactName:      stringPtr("Michael Brown"),
			ContactPhoneNo:   stringPtr("+1-555-0103"),
			ContactEmail:     stringPtr("info@digitalmarketing.example.com"),
			Description:      stringPtr("Digital marketing agency specializing in SEO, PPC, and social media marketing."),
			Address:          stringPtr("789 Marketing Ave"),
			City:             stringPtr("Austin"),
			State:            stringPtr("Texas"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("73301"),
			Value:            float64Ptr(2500000.00),
			BusinessType:     models.BusinessTypeServices,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseGrowth,
			Active:           true,
		},
		{
			OperatorUserID:   users[3].ID,
			Name:             "HealthTech Innovations",
			Tagline:          stringPtr("Revolutionizing healthcare through technology"),
			Website:          stringPtr("https://healthtech.example.com"),
			ContactName:      stringPtr("Emily Davis"),
			ContactPhoneNo:   stringPtr("+1-555-0104"),
			ContactEmail:     stringPtr("contact@healthtech.example.com"),
			Description:      stringPtr("Healthcare technology company developing innovative medical devices and software."),
			Address:          stringPtr("321 Medical Plaza"),
			City:             stringPtr("Boston"),
			State:            stringPtr("Massachusetts"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("02101"),
			Value:            float64Ptr(8000000.00),
			BusinessType:     models.BusinessTypeTechnology,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseStartup,
			Active:           true,
		},
		{
			OperatorUserID:   users[4].ID,
			Name:             "Strategic Consulting Partners",
			Tagline:          stringPtr("Strategic consulting for business excellence"),
			Website:          stringPtr("https://consultingpartners.example.com"),
			ContactName:      stringPtr("David Wilson"),
			ContactPhoneNo:   stringPtr("+1-555-0105"),
			ContactEmail:     stringPtr("info@consultingpartners.example.com"),
			Description:      stringPtr("Management consulting firm specializing in operational efficiency and strategic planning."),
			Address:          stringPtr("654 Business Center"),
			City:             stringPtr("Chicago"),
			State:            stringPtr("Illinois"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("60601"),
			Value:            float64Ptr(3000000.00),
			BusinessType:     models.BusinessTypeConsulting,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseMature,
			Active:           true,
		},
		{
			OperatorUserID:   users[5].ID,
			Name:             "Creative Design Studio",
			Tagline:          stringPtr("Bringing ideas to life through exceptional design"),
			Website:          stringPtr("https://creativestudio.example.com"),
			ContactName:      stringPtr("Lisa Anderson"),
			ContactPhoneNo:   stringPtr("+1-555-0106"),
			ContactEmail:     stringPtr("hello@creativestudio.example.com"),
			Description:      stringPtr("A full-service design agency focused on branding, web design, and digital marketing."),
			Address:          stringPtr("456 Design Avenue"),
			City:             stringPtr("Los Angeles"),
			State:            stringPtr("California"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("90210"),
			Value:            float64Ptr(1800000.00),
			BusinessType:     models.BusinessTypeServices,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseMature,
			Active:           true,
		},
		{
			OperatorUserID:   users[6].ID,
			Name:             "Green Energy Co",
			Tagline:          stringPtr("Sustainable energy solutions for a better tomorrow"),
			Website:          stringPtr("https://greenenergy.example.com"),
			ContactName:      stringPtr("Robert Taylor"),
			ContactPhoneNo:   stringPtr("+1-555-0107"),
			ContactEmail:     stringPtr("info@greenenergy.example.com"),
			Description:      stringPtr("Renewable energy company focused on solar and wind power solutions."),
			Address:          stringPtr("789 Green Street"),
			City:             stringPtr("Denver"),
			State:            stringPtr("Colorado"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("80201"),
			Value:            float64Ptr(4200000.00),
			BusinessType:     models.BusinessTypeManufacturing,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseStartup,
			Active:           true,
		},
		{
			OperatorUserID:   users[7].ID,
			Name:             "FinTech Solutions",
			Tagline:          stringPtr("Next-generation financial technology"),
			Website:          stringPtr("https://fintech.example.com"),
			ContactName:      stringPtr("Jennifer Thomas"),
			ContactPhoneNo:   stringPtr("+1-555-0108"),
			ContactEmail:     stringPtr("contact@fintech.example.com"),
			Description:      stringPtr("Financial technology company providing payment processing and banking solutions."),
			Address:          stringPtr("321 Finance Plaza"),
			City:             stringPtr("Miami"),
			State:            stringPtr("Florida"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("33101"),
			Value:            float64Ptr(6500000.00),
			BusinessType:     models.BusinessTypeTechnology,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseGrowth,
			Active:           true,
		},
		{
			OperatorUserID:   users[8].ID,
			Name:             "RetailTech Innovations",
			Tagline:          stringPtr("Transforming retail through technology"),
			Website:          stringPtr("https://retailtech.example.com"),
			ContactName:      stringPtr("Christopher Jackson"),
			ContactPhoneNo:   stringPtr("+1-555-0109"),
			ContactEmail:     stringPtr("info@retailtech.example.com"),
			Description:      stringPtr("Retail technology company specializing in e-commerce platforms and inventory management."),
			Address:          stringPtr("654 Retail Center"),
			City:             stringPtr("Seattle"),
			State:            stringPtr("Washington"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("98101"),
			Value:            float64Ptr(3200000.00),
			BusinessType:     models.BusinessTypeRetail,
			BusinessCategory: models.BusinessCategoryB2B,
			BusinessPhase:    models.BusinessPhaseMature,
			Active:           true,
		},
		{
			OperatorUserID:   users[9].ID,
			Name:             "EdTech Learning",
			Tagline:          stringPtr("Educational technology for the future"),
			Website:          stringPtr("https://edtech.example.com"),
			ContactName:      stringPtr("Amanda White"),
			ContactPhoneNo:   stringPtr("+1-555-0110"),
			ContactEmail:     stringPtr("hello@edtech.example.com"),
			Description:      stringPtr("Educational technology company creating learning management systems and online courses."),
			Address:          stringPtr("987 Education Blvd"),
			City:             stringPtr("Portland"),
			State:            stringPtr("Oregon"),
			Country:          stringPtr("United States"),
			PostalCode:       stringPtr("97201"),
			Value:            float64Ptr(2100000.00),
			BusinessType:     models.BusinessTypeTechnology,
			BusinessCategory: models.BusinessCategoryB2C,
			BusinessPhase:    models.BusinessPhaseStartup,
			Active:           true,
		},
	}
	
	for i := range businesses {
		dg.db.Create(&businesses[i])
	}
	
	fmt.Println("Generated businesses")
	return businesses
}

func (dg *DataGenerator) generateProjects(users []models.User, businesses []models.Business) []models.Project {
	projects := []models.Project{
		{
			ManagedByUserID: users[0].ID,
			BusinessID:      &businesses[0].ID,
			Name:            "E-commerce Platform Redesign",
			Description:     stringPtr("Complete redesign of the company's e-commerce platform with modern UI/UX and improved performance."),
			ProjectStatus:   models.ProjectStatusActive,
			StartDate:       timePtr(time.Now().AddDate(0, 0, -30)),
			TargetEndDate:   timePtr(time.Now().AddDate(0, 0, 60)),
		},
		{
			ManagedByUserID: users[1].ID,
			BusinessID:      &businesses[1].ID,
			Name:            "Brand Identity Refresh",
			Description:     stringPtr("Comprehensive brand identity update including logo, color scheme, and marketing materials."),
			ProjectStatus:   models.ProjectStatusPlanning,
			StartDate:       timePtr(time.Now().AddDate(0, 0, -10)),
			TargetEndDate:   timePtr(time.Now().AddDate(0, 0, 45)),
		},
		{
			ManagedByUserID: users[2].ID,
			BusinessID:      &businesses[2].ID,
			Name:            "Solar Panel Installation System",
			Description:     stringPtr("Development of an automated system for solar panel installation and monitoring."),
			ProjectStatus:   models.ProjectStatusActive,
			StartDate:       timePtr(time.Now().AddDate(0, 0, -60)),
			TargetEndDate:   timePtr(time.Now().AddDate(0, 0, 90)),
		},
		{
			ManagedByUserID: users[3].ID,
			BusinessID:      &businesses[3].ID,
			Name:            "Patient Management App",
			Description:     stringPtr("Mobile application for healthcare providers to manage patient records and appointments."),
			ProjectStatus:   models.ProjectStatusOnHold,
			StartDate:       timePtr(time.Now().AddDate(0, 0, -90)),
			TargetEndDate:   timePtr(time.Now().AddDate(0, 0, 30)),
		},
		{
			ManagedByUserID: users[4].ID,
			BusinessID:      &businesses[4].ID,
			Name:            "Process Optimization Study",
			Description:     stringPtr("Comprehensive analysis and optimization of client's business processes for improved efficiency."),
			ProjectStatus:   models.ProjectStatusCompleted,
			StartDate:       timePtr(time.Now().AddDate(0, 0, -120)),
			TargetEndDate:   timePtr(time.Now().AddDate(0, 0, -30)),
			ActualEndDate:   timePtr(time.Now().AddDate(0, 0, -25)),
		},
		{
			ManagedByUserID: users[0].ID,
			Name:            "Open Source Library",
			Description:     stringPtr("Development of an open-source Go library for API rate limiting and caching."),
			ProjectStatus:   models.ProjectStatusActive,
			StartDate:       timePtr(time.Now().AddDate(0, 0, -15)),
			TargetEndDate:   timePtr(time.Now().AddDate(0, 0, 30)),
		},
	}
	
	for i := range projects {
		dg.db.Create(&projects[i])
	}
	
	fmt.Println("Generated projects")
	return projects
}

func (dg *DataGenerator) generateProjectMembers(users []models.User, projects []models.Project) []models.ProjectMember {
	var projectMembers []models.ProjectMember
	
	// Add members to projects
	for _, project := range projects {
		// Add 2-4 random members to each project
		numMembers := 2 + rand.Intn(3)
		usedUsers := make(map[uint]bool)
		usedUsers[project.ManagedByUserID] = true // Skip the manager
		
		for j := 0; j < numMembers && j < len(users)-1; j++ {
			var userID uint
			for {
				userID = users[rand.Intn(len(users))].ID
				if !usedUsers[userID] {
					usedUsers[userID] = true
					break
				}
			}
			
			roles := []models.ProjectMemberRole{
				models.ProjectMemberRoleContributor,
				models.ProjectMemberRoleReviewer,
			}
			
			member := models.ProjectMember{
				ProjectID: project.ID,
				UserID:    userID,
				Role:      roles[rand.Intn(len(roles))],
				JoinedAt:  time.Now().AddDate(0, 0, -rand.Intn(30)),
			}
			
			dg.db.Create(&member)
			projectMembers = append(projectMembers, member)
		}
	}
	
	fmt.Println("Generated project members")
	return projectMembers
}

func (dg *DataGenerator) generateProjectSkills(projects []models.Project, skills []models.Skill) []models.ProjectSkill {
	var projectSkills []models.ProjectSkill
	
	for _, project := range projects {
		// Add 3-6 random skills to each project
		numSkills := 3 + rand.Intn(4)
		usedSkills := make(map[uint]bool)
		
		for i := 0; i < numSkills && i < len(skills); i++ {
			var skillID uint
			for {
				skillID = skills[rand.Intn(len(skills))].ID
				if !usedSkills[skillID] {
					usedSkills[skillID] = true
					break
				}
			}
			
			importances := []models.ProjectSkillImportance{
				models.SkillImportanceRequired,
				models.SkillImportancePreferred,
				models.SkillImportanceOptional,
			}
			
			projectSkill := models.ProjectSkill{
				ProjectID:  project.ID,
				SkillID:    skillID,
				Importance: importances[rand.Intn(len(importances))],
			}
			
			dg.db.Create(&projectSkill)
			projectSkills = append(projectSkills, projectSkill)
		}
	}
	
	fmt.Println("Generated project skills")
	return projectSkills
}

func (dg *DataGenerator) generateProjectRegions(projects []models.Project, regions []models.Region) []models.ProjectRegion {
	var projectRegions []models.ProjectRegion
	
	for _, project := range projects {
		// Add 1-3 random regions to each project
		numRegions := 1 + rand.Intn(3)
		usedRegions := make(map[string]bool)
		
		for i := 0; i < numRegions && i < len(regions); i++ {
			var regionID string
			for {
				regionID = regions[rand.Intn(len(regions))].ID
				if !usedRegions[regionID] {
					usedRegions[regionID] = true
					break
				}
			}
			
			projectRegion := models.ProjectRegion{
				RegionID:  regionID,
				ProjectID: project.ID,
			}
			
			dg.db.Create(&projectRegion)
			projectRegions = append(projectRegions, projectRegion)
		}
	}
	
	fmt.Println("Generated project regions")
	return projectRegions
}

func (dg *DataGenerator) generateUserSkills(users []models.User, skills []models.Skill) []models.UserSkill {
	var userSkills []models.UserSkill
	
	// Create strategic skill distributions for better connection recommendations
	userSkillProfiles := [][]struct {
		skillName string
		proficiency models.UserSkillProficiency
	}{
		// User 0 (John - TechCorp) - Technology focused
		{
			{"Go Programming", models.ProficiencyExpert},
			{"Python", models.ProficiencyAdvanced},
			{"Docker", models.ProficiencyAdvanced},
			{"AWS", models.ProficiencyIntermediate},
			{"Project Management", models.ProficiencyIntermediate},
			{"Strategic Planning", models.ProficiencyBeginner},
		},
		// User 1 (Sarah - ConsumerTech) - Technology + Design
		{
			{"JavaScript", models.ProficiencyExpert},
			{"React", models.ProficiencyAdvanced},
			{"Node.js", models.ProficiencyAdvanced},
			{"UI/UX Design", models.ProficiencyIntermediate},
			{"Graphic Design", models.ProficiencyIntermediate},
			{"Marketing Strategy", models.ProficiencyBeginner},
		},
		// User 2 (Michael - Digital Marketing) - Business + Marketing
		{
			{"Marketing Strategy", models.ProficiencyExpert},
			{"Project Management", models.ProficiencyAdvanced},
			{"Financial Analysis", models.ProficiencyIntermediate},
			{"Sales Management", models.ProficiencyIntermediate},
			{"Data Analysis", models.ProficiencyAdvanced},
			{"Go Programming", models.ProficiencyBeginner},
		},
		// User 3 (Emily - HealthTech) - Technology + Data
		{
			{"Python", models.ProficiencyExpert},
			{"Machine Learning", models.ProficiencyAdvanced},
			{"Data Analysis", models.ProficiencyExpert},
			{"PostgreSQL", models.ProficiencyAdvanced},
			{"Project Management", models.ProficiencyIntermediate},
			{"UI/UX Design", models.ProficiencyBeginner},
		},
		// User 4 (David - Consulting) - Business focused
		{
			{"Strategic Planning", models.ProficiencyExpert},
			{"Financial Analysis", models.ProficiencyExpert},
			{"Operations Management", models.ProficiencyAdvanced},
			{"Project Management", models.ProficiencyAdvanced},
			{"Sales Management", models.ProficiencyIntermediate},
			{"JavaScript", models.ProficiencyBeginner},
		},
		// User 5 (Lisa - Creative Design) - Design focused
		{
			{"UI/UX Design", models.ProficiencyExpert},
			{"Graphic Design", models.ProficiencyExpert},
			{"Product Design", models.ProficiencyAdvanced},
			{"Brand Design", models.ProficiencyAdvanced},
			{"Marketing Strategy", models.ProficiencyIntermediate},
			{"Python", models.ProficiencyBeginner},
		},
		// User 6 (Robert - Green Energy) - Manufacturing + Business
		{
			{"Operations Management", models.ProficiencyExpert},
			{"Project Management", models.ProficiencyAdvanced},
			{"Financial Analysis", models.ProficiencyAdvanced},
			{"Strategic Planning", models.ProficiencyIntermediate},
			{"Sales Management", models.ProficiencyIntermediate},
			{"React", models.ProficiencyBeginner},
		},
		// User 7 (Jennifer - FinTech) - Technology + Finance
		{
			{"Python", models.ProficiencyExpert},
			{"Financial Analysis", models.ProficiencyExpert},
			{"Data Analysis", models.ProficiencyAdvanced},
			{"PostgreSQL", models.ProficiencyAdvanced},
			{"Machine Learning", models.ProficiencyIntermediate},
			{"Graphic Design", models.ProficiencyBeginner},
		},
		// User 8 (Christopher - RetailTech) - Technology + Business
		{
			{"JavaScript", models.ProficiencyExpert},
			{"React", models.ProficiencyAdvanced},
			{"Project Management", models.ProficiencyAdvanced},
			{"Operations Management", models.ProficiencyIntermediate},
			{"Sales Management", models.ProficiencyIntermediate},
			{"Machine Learning", models.ProficiencyBeginner},
		},
		// User 9 (Amanda - EdTech) - Technology + Education
		{
			{"Python", models.ProficiencyExpert},
			{"JavaScript", models.ProficiencyAdvanced},
			{"UI/UX Design", models.ProficiencyAdvanced},
			{"Project Management", models.ProficiencyIntermediate},
			{"Marketing Strategy", models.ProficiencyIntermediate},
			{"Financial Analysis", models.ProficiencyBeginner},
		},
	}
	
	for i, user := range users {
		if i < len(userSkillProfiles) {
			// Use predefined skill profile
			for _, skillProfile := range userSkillProfiles[i] {
				// Find skill by name
				for _, skill := range skills {
					if skill.Name == skillProfile.skillName {
						userSkill := models.UserSkill{
							SkillID:          skill.ID,
							UserID:           user.ID,
							ProficiencyLevel: skillProfile.proficiency,
						}
						dg.db.Create(&userSkill)
						userSkills = append(userSkills, userSkill)
						break
					}
				}
			}
		} else {
			// Fallback to random skills for additional users
			numSkills := 5 + rand.Intn(6)
			usedSkills := make(map[uint]bool)
			
			for j := 0; j < numSkills && j < len(skills); j++ {
				var skillID uint
				for {
					skillID = skills[rand.Intn(len(skills))].ID
					if !usedSkills[skillID] {
						usedSkills[skillID] = true
						break
					}
				}
				
				proficiencies := []models.UserSkillProficiency{
					models.ProficiencyBeginner,
					models.ProficiencyIntermediate,
					models.ProficiencyAdvanced,
					models.ProficiencyExpert,
				}
				
				userSkill := models.UserSkill{
					SkillID:          skillID,
					UserID:           user.ID,
					ProficiencyLevel: proficiencies[rand.Intn(len(proficiencies))],
				}
				
				dg.db.Create(&userSkill)
				userSkills = append(userSkills, userSkill)
			}
		}
	}
	
	fmt.Println("Generated user skills")
	return userSkills
}

func (dg *DataGenerator) generateBusinessTags(businesses []models.Business) []models.BusinessTag {
	var businessTags []models.BusinessTag
	
	tagTypes := []models.BusinessTagType{
		models.BusinessTagClient,
		models.BusinessTagService,
		models.BusinessTagSpecialty,
	}
	
	for _, business := range businesses {
		// Add 2-4 random tags to each business
		numTags := 2 + rand.Intn(3)
		
		for i := 0; i < numTags; i++ {
			descriptions := []string{
				"Cloud Computing", "AI/ML", "Web Development", "Mobile Apps",
				"Data Analytics", "Cybersecurity", "DevOps", "UI/UX Design",
				"Digital Marketing", "E-commerce", "SaaS", "Consulting",
				"Healthcare Tech", "FinTech", "EdTech", "Green Tech",
			}
			
			tag := models.BusinessTag{
				BusinessID:  business.ID,
				TagType:     tagTypes[rand.Intn(len(tagTypes))],
				Description: descriptions[rand.Intn(len(descriptions))],
			}
			
			dg.db.Create(&tag)
			businessTags = append(businessTags, tag)
		}
	}
	
	fmt.Println("Generated business tags")
	return businessTags
}

func (dg *DataGenerator) generateBusinessConnections(businesses []models.Business, users []models.User) []models.BusinessConnection {
	var businessConnections []models.BusinessConnection
	
	// Create strategic business connections that demonstrate different connection types
	strategicConnections := []struct {
		initiatingIdx int
		receivingIdx  int
		connectionType models.BusinessConnectionType
		status        models.BusinessConnectionStatus
		notes         string
	}{
		// TechCorp (0) -> ConsumerTech (1) - Partnership (Complementary Partners)
		{0, 1, models.ConnectionTypePartnership, models.ConnectionStatusActive, "Partnership for cross-selling B2B and B2C solutions"},
		
		// TechCorp (0) -> Digital Marketing Pro (2) - Collaboration (Alliance Partners)
		{0, 2, models.ConnectionTypeCollaboration, models.ConnectionStatusActive, "Collaboration on marketing campaigns for tech solutions"},
		
		// Digital Marketing Pro (2) -> Creative Design Studio (5) - Partnership (Alliance Partners)
		{2, 5, models.ConnectionTypePartnership, models.ConnectionStatusActive, "Partnership for full-service marketing and design"},
		
		// HealthTech (3) -> FinTech (7) - Referral (Complementary Partners)
		{3, 7, models.ConnectionTypeReferral, models.ConnectionStatusActive, "Referral partnership for healthcare payment solutions"},
		
		// Strategic Consulting (4) -> TechCorp (0) - Client (Complementary Partners)
		{4, 0, models.ConnectionTypeClient, models.ConnectionStatusActive, "Consulting services for technology strategy"},
		
		// FinTech (7) -> RetailTech (8) - Supplier (Complementary Partners)
		{7, 8, models.ConnectionTypeSupplier, models.ConnectionStatusActive, "Payment processing services for retail platform"},
		
		// EdTech (9) -> TechCorp (0) - Collaboration (Alliance Partners)
		{9, 0, models.ConnectionTypeCollaboration, models.ConnectionStatusPending, "Collaboration on educational technology solutions"},
		
		// Green Energy (6) -> Strategic Consulting (4) - Client (Complementary Partners)
		{6, 4, models.ConnectionTypeClient, models.ConnectionStatusActive, "Consulting services for renewable energy strategy"},
		
		// Creative Design Studio (5) -> ConsumerTech (1) - Supplier (Complementary Partners)
		{5, 1, models.ConnectionTypeSupplier, models.ConnectionStatusActive, "Design services for mobile applications"},
		
		// RetailTech (8) -> Digital Marketing Pro (2) - Client (Complementary Partners)
		{8, 2, models.ConnectionTypeClient, models.ConnectionStatusActive, "Digital marketing services for retail platform"},
	}
	
	for _, conn := range strategicConnections {
		if conn.initiatingIdx < len(businesses) && conn.receivingIdx < len(businesses) {
			connection := models.BusinessConnection{
				InitiatingBusinessID: businesses[conn.initiatingIdx].ID,
				ReceivingBusinessID:  businesses[conn.receivingIdx].ID,
				ConnectionType:       conn.connectionType,
				Status:               conn.status,
				InitiatedByUserID:    users[conn.initiatingIdx].ID,
				Notes:                stringPtr(conn.notes),
			}
			
			dg.db.Create(&connection)
			businessConnections = append(businessConnections, connection)
		}
	}
	
	fmt.Println("Generated business connections")
	return businessConnections
}

func (dg *DataGenerator) generatePublications(users []models.User, businesses []models.Business) []models.Publication {
	var publications []models.Publication
	
	publicationTypes := []models.PublicationType{
		models.PublicationPost,
		models.PublicationCaseStudy,
		models.PublicationTestimonial,
		models.PublicationArticle,
	}
	
	titles := []string{
		"The Future of Cloud Computing",
		"Building Scalable APIs with Go",
		"Design Thinking in Product Development",
		"Case Study: Digital Transformation Success",
		"Best Practices for Remote Team Management",
		"Understanding Machine Learning Algorithms",
		"Customer Success Story: E-commerce Platform",
		"The Art of Effective Communication",
		"Data-Driven Decision Making",
		"Agile Development Methodologies",
	}
	
	for i := 0; i < 15; i++ {
		user := users[rand.Intn(len(users))]
		var businessID *uint
		if rand.Float32() < 0.7 { // 70% chance of having a business
			businessID = &businesses[rand.Intn(len(businesses))].ID
		}
		
		title := titles[rand.Intn(len(titles))]
		slug := fmt.Sprintf("%s-%d", title, i+1)
		
		publication := models.Publication{
			UserID:          user.ID,
			BusinessID:      businessID,
			PublicationType: publicationTypes[rand.Intn(len(publicationTypes))],
			Title:           title,
			Slug:            slug,
			Excerpt:         stringPtr("This is a sample excerpt for the publication."),
			Content:         "This is the full content of the publication. It contains detailed information about the topic and provides valuable insights for readers.",
			Thumbnail:       stringPtr(fmt.Sprintf("https://example.com/images/%s.jpg", slug)),
			VideoURL:        stringPtr(fmt.Sprintf("https://example.com/videos/%s.mp4", slug)),
			Published:       rand.Float32() < 0.8, // 80% chance of being published
			PublishedAt:     timePtr(time.Now().AddDate(0, 0, -rand.Intn(90))),
		}
		
		dg.db.Create(&publication)
		publications = append(publications, publication)
	}
	
	fmt.Println("Generated publications")
	return publications
}

func (dg *DataGenerator) generateNotifications(users []models.User) []models.Notification {
	var notifications []models.Notification
	
	notificationTypes := []models.NotificationType{
		"connection_request",
		"project_invite",
		"message",
		"system",
	}
	
	relatedEntityTypes := []models.RelatedEntityType{
		"business",
		"project",
		"publication",
		"idea",
	}
	
	titles := []string{
		"New Connection Request",
		"Project Invitation",
		"New Message Received",
		"System Update Available",
		"Profile View Notification",
		"Skill Endorsement",
		"Project Status Update",
		"Payment Received",
		"Meeting Reminder",
		"Account Security Alert",
	}
	
	messages := []string{
		"You have received a new connection request.",
		"You've been invited to join a new project.",
		"You have a new message waiting for you.",
		"A new system update is available for download.",
		"Someone viewed your profile recently.",
		"Your skill has been endorsed by a colleague.",
		"The project status has been updated.",
		"Your payment has been processed successfully.",
		"Don't forget about your upcoming meeting.",
		"Please review your account security settings.",
	}
	
	for i := 0; i < 25; i++ {
		receiver := users[rand.Intn(len(users))]
		var senderID *uint
		if rand.Float32() < 0.8 { // 80% chance of having a sender
			senderID = &users[rand.Intn(len(users))].ID
		}
		
		notification := models.Notification{
			SenderUserID:      senderID,
			ReceiverUserID:    receiver.ID,
			NotificationType:  notificationTypes[rand.Intn(len(notificationTypes))],
			Title:             titles[rand.Intn(len(titles))],
			Message:           messages[rand.Intn(len(messages))],
			RelatedEntityType: &relatedEntityTypes[rand.Intn(len(relatedEntityTypes))],
			RelatedEntityID:   uintPtr(uint(rand.Intn(100) + 1)),
			Read:              rand.Float32() < 0.3, // 30% chance of being read
			ActionURL:         stringPtr(fmt.Sprintf("https://app.example.com/action/%d", i+1)),
		}
		
		dg.db.Create(&notification)
		notifications = append(notifications, notification)
	}
	
	fmt.Println("Generated notifications")
	return notifications
}

func (dg *DataGenerator) generateFeedback() []models.Feedback {
	feedback := []models.Feedback{
		{
			Name:    "Alice Johnson",
			Email:   "alice.johnson@example.com",
			Content: "Great platform! The user interface is intuitive and the features are exactly what we needed for our business.",
		},
		{
			Name:    "Bob Smith",
			Email:   "bob.smith@example.com",
			Content: "Excellent service and support. The team was very responsive to our questions and helped us get started quickly.",
		},
		{
			Name:    "Carol Davis",
			Email:   "carol.davis@example.com",
			Content: "The platform has significantly improved our workflow efficiency. Highly recommended for any growing business.",
		},
		{
			Name:    "David Wilson",
			Email:   "david.wilson@example.com",
			Content: "Some minor issues with the mobile app, but overall a solid product. Looking forward to future updates.",
		},
		{
			Name:    "Eva Brown",
			Email:   "eva.brown@example.com",
			Content: "Outstanding customer support and a feature-rich platform. It has transformed how we manage our projects.",
		},
	}
	
	for i := range feedback {
		dg.db.Create(&feedback[i])
	}
	
	fmt.Println("Generated feedback")
	return feedback
}

func (dg *DataGenerator) generateDailyActivities() []models.DailyActivity {
	activities := []models.DailyActivity{
		{
			Name:        "Morning Exercise",
			Description: "Start your day with 30 minutes of physical activity to boost energy and focus.",
		},
		{
			Name:        "Daily Standup",
			Description: "Participate in the team's daily standup meeting to stay aligned and informed.",
		},
		{
			Name:        "Code Review",
			Description: "Review at least one pull request or code submission to maintain code quality.",
		},
		{
			Name:        "Learning Session",
			Description: "Spend 30 minutes learning something new related to your field or interests.",
		},
		{
			Name:        "Client Communication",
			Description: "Reach out to at least one client or stakeholder to maintain relationships.",
		},
		{
			Name:        "Documentation Update",
			Description: "Update project documentation or create new documentation as needed.",
		},
		{
			Name:        "Team Collaboration",
			Description: "Engage in collaborative work with team members on shared projects.",
		},
		{
			Name:        "Reflection Time",
			Description: "Take 15 minutes to reflect on the day's accomplishments and plan for tomorrow.",
		},
	}
	
	for i := range activities {
		dg.db.Create(&activities[i])
	}
	
	fmt.Println("Generated daily activities")
	return activities
}

func (dg *DataGenerator) generateDailyActivityEnrolments(users []models.User, activities []models.DailyActivity) []models.DailyActivityEnrolment {
	var enrolments []models.DailyActivityEnrolment
	
	for _, user := range users {
		// Each user enrolls in 3-5 random activities
		numActivities := 3 + rand.Intn(3)
		usedActivities := make(map[uint]bool)
		
		for i := 0; i < numActivities && i < len(activities); i++ {
			var activityID uint
			for {
				activityID = activities[rand.Intn(len(activities))].ID
				if !usedActivities[activityID] {
					usedActivities[activityID] = true
					break
				}
			}
			
			enrolment := models.DailyActivityEnrolment{
				DailyActivityID: activityID,
				UserID:          user.ID,
			}
			
			dg.db.Create(&enrolment)
			enrolments = append(enrolments, enrolment)
		}
	}
	
	fmt.Println("Generated daily activity enrolments")
	return enrolments
}

func (dg *DataGenerator) generateUserDailyActivityProgress(users []models.User, activities []models.DailyActivity) []models.UserDailyActivityProgress {
	var progress []models.UserDailyActivityProgress
	
	statuses := []models.DailyActivityProgressStatus{
		models.ProgressStatusNotStarted,
		models.ProgressStatusInProgress,
		models.ProgressStatusCompleted,
	}
	
	// Generate progress for the last 30 days
	for _, user := range users {
		for _, activity := range activities {
			// 50% chance of having progress for each user-activity combination
			if rand.Float32() < 0.5 {
				date := time.Now().AddDate(0, 0, -rand.Intn(30))
				status := statuses[rand.Intn(len(statuses))]
				
				progressEntry := models.UserDailyActivityProgress{
					UserID:          user.ID,
					DailyActivityID: activity.ID,
					Date:            datatypes.Date(date),
					Status:          status,
					Progress:        rand.Intn(101), // 0-100%
				}
				
				dg.db.Create(&progressEntry)
				progress = append(progress, progressEntry)
			}
		}
	}
	
	fmt.Println("Generated user daily activity progress")
	return progress
}

func (dg *DataGenerator) generateEvents(users []models.User) []models.Event {
	var events []models.Event
	
	eventTypes := []string{
		"user_login",
		"user_logout",
		"project_created",
		"project_updated",
		"business_connection_requested",
		"skill_added",
		"publication_published",
		"notification_sent",
		"daily_activity_completed",
		"profile_updated",
	}
	
	for i := 0; i < 50; i++ {
		user := users[rand.Intn(len(users))]
		eventType := eventTypes[rand.Intn(len(eventTypes))]
		
		payload := map[string]interface{}{
			"event_id":   fmt.Sprintf("evt_%d", i+1),
			"timestamp":  time.Now().Unix(),
			"metadata": map[string]interface{}{
				"source": "data_generator",
				"version": "1.0",
			},
		}
		
		// Convert map to JSON bytes
		payloadBytes, _ := json.Marshal(payload)
		
		event := models.Event{
			EventType: eventType,
			Payload:   datatypes.JSON(payloadBytes),
			Timestamp: time.Now().AddDate(0, 0, -rand.Intn(30)),
			UserID:    &user.ID,
		}
		
		dg.db.Create(&event)
		events = append(events, event)
	}
	
	fmt.Println("Generated events")
	return events
}

func (dg *DataGenerator) generateL2EResponses(users []models.User) []models.L2EResponse {
	var responses []models.L2EResponse
	
	for i := 0; i < 20; i++ {
		user := users[rand.Intn(len(users))]
		
		response := map[string]interface{}{
			"question_id": fmt.Sprintf("q_%d", i+1),
			"answer":      fmt.Sprintf("Sample response %d", i+1),
			"confidence":  rand.Float64(),
			"category":    "learning",
			"timestamp":   time.Now().Unix(),
		}
		
		// Convert map to JSON bytes
		responseBytes, _ := json.Marshal(response)
		
		l2eResponse := models.L2EResponse{
			UserID:    user.ID,
			Response:  datatypes.JSON(responseBytes),
			DateAdded: time.Now().AddDate(0, 0, -rand.Intn(30)),
		}
		
		dg.db.Create(&l2eResponse)
		responses = append(responses, l2eResponse)
	}
	
	fmt.Println("Generated L2E responses")
	return responses
}

func (dg *DataGenerator) generateUserConfigs(users []models.User) []models.UserConfig {
	var configs []models.UserConfig
	
	configTypes := []string{
		"preferences",
		"notifications",
		"privacy",
		"display",
		"security",
	}
	
	for _, user := range users {
		for _, configType := range configTypes {
			config := map[string]interface{}{
				"enabled":     rand.Float32() < 0.8,
				"frequency":   "daily",
				"theme":       "light",
				"language":    "en",
				"timezone":    "UTC",
				"created_at":  time.Now().Unix(),
			}
			
			// Convert map to JSON bytes
			configBytes, _ := json.Marshal(config)
			
			userConfig := models.UserConfig{
				UserID:     user.ID,
				ConfigType: configType,
				Config:     datatypes.JSON(configBytes),
			}
			
			dg.db.Create(&userConfig)
			configs = append(configs, userConfig)
		}
	}
	
	fmt.Println("Generated user configs")
	return configs
}

func (dg *DataGenerator) generateProjectApplicants(users []models.User, projects []models.Project) []models.ProjectApplicant {
	var applicants []models.ProjectApplicant
	
	// 20% of users apply to 1-3 random projects
	for _, user := range users {
		if rand.Float32() < 0.2 {
			numApplications := 1 + rand.Intn(3)
			usedProjects := make(map[uint]bool)
			
			for i := 0; i < numApplications && i < len(projects); i++ {
				var projectID uint
				for {
					projectID = projects[rand.Intn(len(projects))].ID
					if !usedProjects[projectID] {
						usedProjects[projectID] = true
						break
					}
				}
				
				applicant := models.ProjectApplicant{
					ProjectID: projectID,
					UserID:    user.ID,
				}
				
				dg.db.Create(&applicant)
				applicants = append(applicants, applicant)
			}
		}
	}
	
	fmt.Println("Generated project applicants")
	return applicants
}

func (dg *DataGenerator) generateInferredConnections(businesses []models.Business, projects []models.Project, users []models.User) []models.InferredConnection {
	var connections []models.InferredConnection
	
	entityTypes := []string{"business", "project", "user"}
	connectionTypes := []string{"similar_skills", "geographic_proximity", "industry_match", "collaboration_potential"}
	
	for i := 0; i < 30; i++ {
		connection := models.InferredConnection{
			SourceEntityType: entityTypes[rand.Intn(len(entityTypes))],
			SourceEntityID:   uint(rand.Intn(20) + 1),
			TargetEntityType: entityTypes[rand.Intn(len(entityTypes))],
			TargetEntityID:   uint(rand.Intn(20) + 1),
			ConnectionType:   connectionTypes[rand.Intn(len(connectionTypes))],
			ConfidenceScore:  rand.Float64(),
			ModelVersion:     "v1.0",
		}
		
		dg.db.Create(&connection)
		connections = append(connections, connection)
	}
	
	fmt.Println("Generated inferred connections")
	return connections
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func uintPtr(u uint) *uint {
	return &u
}
