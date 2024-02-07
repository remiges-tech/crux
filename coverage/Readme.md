TO RUN TEST COVERAGE run below commands:

Step-1. Add your test file path as per package ans also modify name according to package name:

go test -coverprofile=./coverage/coverage_schema.out -coverpkg=./... server/schema/schemaTest/schemaList_test.go server/schema/schemaTest/schemaNew_test.go server/schema/schemaTest/schemaUpdate_test.go server/schema/schemaTest/schema_get_test.go server/schema/schemaTest/main_test.go server/schema/schemaTest/schema_delete_test.go

Step-2. After step-1 run this command to generate html file of coverage  & modify your file path within below command accordingly:

go tool cover -html=./coverage/coverage_schema.out -o ./coverage/coverage_schema.html

_________________________________________________________________________________________

1. coverage_schema :
go test -coverprofile=./coverage/coverage_schema.out -coverpkg=./... server/schema/schemaTest/schemaList_test.go server/schema/schemaTest/schemaNew_test.go server/schema/schemaTest/schemaUpdate_test.go server/schema/schemaTest/schema_get_test.go server/schema/schemaTest/main_test.go server/schema/schemaTest/schema_delete_test.go

go tool cover -html=./coverage/coverage_schema.out -o ./coverage/coverage_schema.html


2. coverage_workflow :
go test -coverprofile=./coverage/coverage_workflow.out -coverpkg=./... server/workflow/workflowTest/workflow_get_test.go server/workflow/workflowTest/main_test.go

go tool cover -html=./coverage/coverage_workflow.out -o ./coverage/coverage_workflow.html
