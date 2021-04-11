scripts_dir = $(shell pwd)/scripts
run_script = /bin/bash $(scripts_dir)

help:
	@echo ---------------------labsystem---------------------
	@echo "generate_rsa_key      ---- repeatedly generate rsa_key"
	@echo "db_migrate            ---- migrate database"
	@echo "test                  ---- execute go test to all test"
	@echo "generate_verify_code  ---- repeatedly generate verify image"
	@echo "run                   ---- run main"

generate_rsa_key:
	$(run_script)/generate_rsa_key.sh
db_migrate:
	$(run_script)/db_migrate.sh
generate_verify_code:
	$(run_script)/generate_verify_code.sh
test:
	$(run_script)/run_test.sh
run:
	$(run_script)/run.sh