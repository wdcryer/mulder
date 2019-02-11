# Mulder Chart

CHART_NAME=mulder
OUTPUT_DIR=../output

all: clean style test dry-run build

style: lint-helm lint-yaml

lint-helm: dependencies
	helm lint .

lint-yaml: render-templates
	yamllint -c ../.yamllint *.yaml
	yamllint -c ../.yamllint tests/*.yaml
	yamllint -c ../.yamllint ${OUTPUT_DIR}/${CHART_NAME}/templates/*.yaml

test: test-unit

test-unit:
	helm unittest .

dry-run: render-templates
	kubectl apply --dry-run -f ${OUTPUT_DIR}/${CHART_NAME}/templates

build: ensure-output-dir-exists dependencies
	helm package --destination ${OUTPUT_DIR} .

deploy: ensure-myvalues-file-exists
	helm template --name ${CHART_NAME} --values myvalues.yaml . | kubectl apply -f -

render-templates: ensure-output-dir-exists
	helm template --name ${CHART_NAME} --output-dir ${OUTPUT_DIR} .

dependencies: ensure-helm-repos
	@helm dependency build

ensure-output-dir-exists:
	@mkdir -p ${OUTPUT_DIR}

ensure-myvalues-file-exists:
	@touch myvalues.yaml

ensure-helm-repos:
	@helm repo add stable https://kubernetes-charts.storage.googleapis.com

clean:
	@rm -f requirements.lock
	@rm -rf charts
	@rm -rf ${OUTPUT_DIR}/${CHART_NAME}*
