# tia api

## Runnning Scripts
go install github.com/roblaszczak/vgt@latest
go install github.com/goptics/vizb@latest

scripts
├── benchmark
├── benchmark.sh
├── code_coverage
├── code_coverage.sh
└── tests
    ├── run_integration.sh
    ├── run_unit.sh
    └── vis_integration.sh

## Tests
test
├── integration
│   ├── auth.service.integration_test.go
│   ├── business_connection.service.integration_test.go
│   ├── business.service.integration_test.go
│   ├── business_tag.service.integration_test.go
│   ├── daily_activity_enrolment.service.integration_test.go
│   ├── daily_activity.service.integration_test.go
│   ├── event.service.integration_test.go
│   ├── feedback.service.integration_test.go
│   ├── idea.service.integration_test.go
│   ├── idea_vote.service.integration_test.go
│   ├── inferred_connection.service.integration_test.go
│   ├── l2e.service.integration_test.go
│   ├── main_test.go
│   ├── migrate.integration_test.go
│   ├── notification.service.integration_test.go
│   ├── project_applicant.service.integration_test.go
│   ├── project_member.service.integration_test.go
│   ├── project_region.service.integration_test.go
│   ├── project.service.integration_test.go
│   ├── project_skill.service.integration_test.go
│   ├── publication.service.integration_test.go
│   ├── region.service.integration_test.go
│   ├── skill.service.integration_test.go
│   ├── subscription.service.integration_test.go
│   ├── user_config.service.integration_test.go
│   ├── user.service.integration_test.go
│   ├── user_session.service.integration_test.go
│   ├── user_skill.service.integration_test.go
│   └── user_subscription.service.integration_test.go
├── mutation
└── unit
