PLUGIN_NAME=samplefs
PLUGIN_TAG=:latest


all: clean docker rootfs create enable

clean:
	@echo "### rm ./plugin"
	@rm -rf ./plugin

docker:
	@echo "### docker build: builder image"
	@docker build -t ${PLUGIN_NAME}:rootfs .

rootfs:
	@echo "### create rootfs directory in ./plugin/rootfs"
	@mkdir -p ./plugin/rootfs
	@docker create --name tmp ${PLUGIN_NAME}:rootfs
	@docker export tmp | tar -x -C ./plugin/rootfs
	@echo "### copy config.json to ./plugin/"
	@cp config.json ./plugin/
	@docker rm -vf tmp

create:
	@echo "### remove existing plugin ${PLUGIN_NAME}${PLUGIN_TAG} if exists"
	@docker plugin disable ${PLUGIN_NAME}${PLUGIN_TAG} || true
	@docker plugin rm ${PLUGIN_NAME}${PLUGIN_TAG} || true
	@echo "### create new plugin ${PLUGIN_NAME}${PLUGIN_TAG} from ./plugin"
	@docker plugin create ${PLUGIN_NAME}${PLUGIN_TAG} ./plugin

enable:
	@echo "### enable plugin ${PLUGIN_NAME}${PLUGIN_TAG}"
	@docker plugin enable ${PLUGIN_NAME}${PLUGIN_TAG}