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

jaadu: generate tern

pg-drop-all:
	cd db/migrations/; tern migrate --destination 0

db-migrate-generate: pg-drop-all tern generate

.PHONY: newSchema