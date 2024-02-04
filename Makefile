# -- Variables --

_DEV=deploy/dev
_PROD=deploy/prod

_COMPOSE=docker-compose.yaml
COMPOSE_DEV=${_DEV}/${_COMPOSE}
COMPOSE_PROD=${_PROD}/${_COMPOSE}

_CONFIG=config.yaml
CONFIG_DEV=${_DEV}/${_CONFIG}
CONFIG_PROD=${_PROD}/${_CONFIG}

_ENV=.env
ENV_DEV=${_DEV}/${_ENV}
ENV_PROD=${_PROD}/${_ENV}

_ENV_EXAMPLE=.env.example
ENV_EXAMPLE_DEV=${_DEV}/${_ENV_EXAMPLE}
ENV_EXAMPLE_PROD=${_PROD}/${_ENV_EXAMPLE}

BUILD_DIR=build
MAIN=main.go

# -- Commads --

prepare:
	cp ${ENV_EXAMPLE_DEV} ${ENV_DEV} \
	&& cp ${ENV_EXAMPLE_PROD} ${ENV_PROD} \
	&& go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	
run-dev: SHELL:=/bin/bash
run-dev:
	docker-compose -f ${COMPOSE_DEV} up -d --no-recreate \
	&& set -a \
	&& source ${ENV_DEV} \
	&& go run ${MAIN} --config ${CONFIG_DEV}
	
run-prod:
	docker-compose -f ${COMPOSE_PROD} up -d

build-binary:
	go build -o ${BUILD_DIR}/app ${MAIN}

build-docker:
	docker build -t tcaty/update-watcher .
	
clean:
	docker-compose -f ${COMPOSE_DEV} down --remove-orphans --volumes; \
	docker-compose -f ${COMPOSE_PROD} down --remove-orphans --volumes; \
	rm -rf ${BUILD_DIR}