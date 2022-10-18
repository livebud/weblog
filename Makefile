migrate := go run github.com/matthewmueller/migrate/cmd/migrate
pogo := go run github.com/matthewmueller/pogo/cmd/pogo

migrate.up: pogo
	@ $(migrate) up $(DATABASE_URL)

migrate.reset:
	@ $(migrate) down $(DATABASE_URL)
	@ $(migrate) up $(DATABASE_URL)
	@ $(MAKE) pogo

pogo:
	@ $(pogo) --db $(DATABASE_URL) --dir bud/pkg/table
