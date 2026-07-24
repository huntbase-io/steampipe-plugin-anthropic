STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/local/anthropic/anthropic.plugin -tags "$(BUILD_TAGS)" *.go

# Build the plugin as a SQLite loadable extension (steampipe_sqlite_anthropic.so)
# using turbot/steampipe-sqlite with a replace directive pointing at this repo.
sqlite:
	rm -rf build/steampipe-sqlite && mkdir -p build
	git clone --quiet --depth 1 https://github.com/turbot/steampipe-sqlite.git build/steampipe-sqlite
	cd build/steampipe-sqlite && rm -rf work && mkdir -p work && rsync -a --exclude='.git' . work/
	cd build/steampipe-sqlite/work && \
		go run generate/generator.go templates . anthropic "" github.com/huntbase-io/steampipe-plugin-anthropic && \
		echo 'replace github.com/huntbase-io/steampipe-plugin-anthropic => $(CURDIR)' >> go.mod && \
		go mod tidy && \
		$(MAKE) -f out/Makefile build
	cp build/steampipe-sqlite/work/steampipe_sqlite_anthropic.so .
	@echo "Built ./steampipe_sqlite_anthropic.so"

# Build the plugin as a standalone Postgres FDW extension
# (steampipe_postgres_anthropic.so + .control + .sql in build/postgres/) using
# turbot/steampipe-postgres-fdw. Requires pg_config on PATH.
postgres:
	rm -rf build/steampipe-postgres-fdw && mkdir -p build
	git clone --quiet --depth 1 https://github.com/turbot/steampipe-postgres-fdw.git build/steampipe-postgres-fdw
	cd build/steampipe-postgres-fdw && \
		$(MAKE) prebuild.go && \
		rm -rf work && mkdir -p work && rsync -a --exclude='.git' . work/
	cd build/steampipe-postgres-fdw/work && \
		go run generate/generator.go templates . anthropic "" github.com/huntbase-io/steampipe-plugin-anthropic && \
		echo 'replace github.com/huntbase-io/steampipe-plugin-anthropic => $(CURDIR)' >> go.mod && \
		go mod tidy && \
		$(MAKE) -C ./fdw clean plugin=anthropic && \
		$(MAKE) -C ./fdw go plugin=anthropic && \
		$(MAKE) -C ./fdw plugin=anthropic && \
		$(MAKE) -C ./fdw standalone plugin=anthropic
	rm -rf build/postgres && mkdir -p build/postgres
	cp -a build/steampipe-postgres-fdw/work/build-$(shell uname)/* build/postgres/
	@echo "Built build/postgres/ (steampipe_postgres_anthropic extension)"

clean:
	rm -rf build steampipe_sqlite_anthropic.so
