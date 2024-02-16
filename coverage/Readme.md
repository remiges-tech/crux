TO RUN TEST COVERAGE run below commands:

Step-1. Add your test file path as per package ans also modify name according to package name:

go test -coverprofile=./coverage/coverage_schema.out -coverpkg=./... server/schema/test/schemaList_test.go server/schema/test/schemaNew_test.go server/schema/test/schemaUpdate_test.go server/schema/test/schema_get_test.go server/schema/test/main_test.go server/schema/test/schema_delete_test.go

Step-2. After step-1 run this command to generate html file of coverage  & modify your file path within below command accordingly:

go tool cover -html=./coverage/coverage_schema.out -o ./coverage/coverage_schema.html

_________________________________________________________________________________________

1. coverage_schema :
go test -coverprofile=./coverage/coverage_schema.out -coverpkg=./... server/schema/test/schemaList_test.go server/schema/test/schemaNew_test.go server/schema/test/schemaUpdate_test.go server/schema/test/schema_get_test.go server/schema/test/main_test.go server/schema/test/schema_delete_test.go

<!-- go tool cover -html=./coverage/coverage_schema.out -o ./coverage/coverage_schema.html -->


2. coverage_workflow :
go test -coverprofile=./coverage/coverage_workflow.out -coverpkg=./... server/workflow/test/workflow_delete_test.go server/workflow/test/workflow_get_test.go server/workflow/test/workflow_list_test.go server/workflow/test/workflow_new_test.go server/workflow/test/workflow_update_test.go  server/workflow/test/main_test.go

go tool cover -html=./coverage/coverage_workflow.out -o ./coverage/coverage_workflow.html
