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

.PHONY: newSchema