.PHONY: husky format migrate

husky:
	chmod +x formatter.sh; chmod +x pre-commit.sh; cp pre-commit.sh .git/hooks/pre-commit

format:
	./formatter.sh .

migrate_up:
	migrate -path db/migration -database "postgresql://root:root@localhost:5433/superindo?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:root@localhost:5433/superindo?sslmode=disable" -verbose down