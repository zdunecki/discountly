up_local:
	docker-compose up
restart_local:
	docker-compose restart app
clear_e2e:
	docker-compose -f docker-compose.e2e.yml rm -v -f
e2e:
	docker-compose -f docker-compose.e2e.yml rm -v -f
	docker-compose -f docker-compose.e2e.yml up --exit-code-from e2e
	docker-compose -f docker-compose.e2e.yml rm -v -f
unit-test:
	go test ./features/...