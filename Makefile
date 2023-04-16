.PHONY: build build-backend build-frontend init-db run

build: build-frontend build-backend

build-backend:
	cd backend && go build

build-frontend:
	cd frontend && bash -l -c "nvm exec node yarn build"

init-db:
	edgedb project init --non-interactive && edgedb migrate

run:
	./backend/backend