.PHONY: env-up env-down into-cli cleanfiles env-clean



all: env-clean env-up into-cli
# start networks
env-up:
	@echo "start environment"
	@docker-compose -f fixtures/docker-compose-cli.yaml up -d
	@echo "environment has been started"
	@docker ps

env-down:
	@echo "stop environment"
	@docker-compose -f fixtures/docker-compose-cli.yaml down
	@echo "environment has been stopped"
	@docker ps


env-clean: env-down
	@echo "clean up environment"
	@docker volume prune   -f   # 清理挂载卷
	@docker network prune  -f    # 来清理没有再被任何容器引用的networks
	@docker rm -f -v `docker ps -a --no-trunc | grep "mycc" | cut -d ' ' -f 1` 2>/dev/null || true  #清除不截断输出的容器
	@docker rmi `docker images --no-trunc | grep "mycc" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "environment has been cleaned"


into-cli:
	@docker exec -it cli bash


cleanfiles:
	@rm -rf ./fixtures/crypto-config
	@rm -f ./fixtures/channel-artifacts/*
