migrate.up: pogo
	@ migrate up $(DATABASE_URL)

migrate.reset:
	@ migrate down $(DATABASE_URL)
	@ migrate up $(DATABASE_URL)
	@ $(MAKE) pogo

pogo:
	@ pogo --db $(DATABASE_URL) --dir bud/pkg/table
