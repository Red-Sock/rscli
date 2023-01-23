# SERVICE_NAME - имя сервиса который хотим сбилдить
.PHONY: build-app
build-app:
	cd ../${SERVICE_NAME}/ && HOSTOS=linux HOSTARCH=arm64 make --makefile scratch.mk build
	mv ../${SERVICE_NAME}/bin/${SERVICE_NAME} ./${SERVICE_NAME}
	cp ../${SERVICE_NAME}/.o3/k8s/values.yaml ./${SERVICE_NAME}/.o3/k8s/values.yaml
	cp ../${SERVICE_NAME}/.o3/k8s/platform.yaml ./${SERVICE_NAME}/.o3/k8s/platform.yaml
	cp ../${SERVICE_NAME}/.o3/k8s/values_development.yaml ./${SERVICE_NAME}/.o3/k8s/values_local.yaml

.PHONY: start-app
start-app:
	rm -rf ${SERVICE_NAME}config


# ENVIRONMENTS - перечисленные через запятую окружения, которые хотим поднять
.PHONY: validate-env
validate-env:
	go run v