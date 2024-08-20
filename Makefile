CURRENT_DIR=$(shell pwd)
DBURL := postgres://postgres:1918@localhost:5432/car_auth?sslmode=disable


proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

swag-init:
	swag init -g api/router.go -o api/docs


mig-up:
	migrate -path databases/migrations -database '${DBURL}' -verbose up

mig-down:
	migrate -path databases/migrations -database '${DBURL}' -verbose down

mig-force:
	migrate -path databases/migrations -database '${DBURL}' -verbose force 1


mig-create-user:
	migrate create -ext sql -dir databases/migrations -seq create_refreshtokens_table

