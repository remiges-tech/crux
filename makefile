# create new schema 
newSchema:
	cd db/migrations; tern new crux

# tern migration
tern:
	cd db/migrations; tern migrate 

# sqlc generate
generate:
	cd db; sqlc generate

# start an etcd server
etcd:
	cd; cd etcd/bin; ./etcd


pg-drop-all:
	cd db/migrations/; tern migrate --destination 0

jaadu: pg-drop-all generate tern

run-dev: go run -tags dev main.go middleware_dev.go 

db-migrate-generate: pg-drop-all tern generate

test-all:
	go test -coverprofile=./coverage/coverage_schema.out -coverpkg=./... server/schema/test/schemaList_test.go server/schema/test/schemaNew_test.go server/schema/test/schemaUpdate_test.go server/schema/test/schema_get_test.go server/schema/test/main_test.go server/schema/test/schema_delete_test.go; go tool cover -html=./coverage/coverage_schema.out -o ./coverage/coverage_schema.html; go test -coverprofile=./coverage/coverage_workflow.out -coverpkg=./... server/workflow/test/workflow_delete_test.go server/workflow/test/workflow_get_test.go server/workflow/test/workflow_list_test.go server/workflow/test/workflow_new_test.go server/workflow/test/workflow_update_test.go  server/workflow/test/main_test.go; go tool cover -html=./coverage/coverage_workflow.out -o ./coverage/coverage_workflow.html -html=./coverage/coverage_schema.out -o ./coverage/coverage_schema.html;

.PHONY: newSchema