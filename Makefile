.PHONY: prepare generate frontend backend build

prepare:
	cd backend && go install -v github.com/searKing/golang/tools/go-enum

generate:
	cd backend && go generate -v -x ./...

frontend:
	cd frontend && npm run build
	-rm -r backend/resources
	cp -r frontend/dist backend/resources

backend:
	cd backend && go build .

build: frontend generate backend
