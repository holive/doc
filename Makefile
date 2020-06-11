PROJECT_NAME = doc
PROJECT_TAG = dev
DOCKER_IMAGE = hbliveira/doc

build:
	mkdir -p $(PWD)/dist/$(PROJECT_NAME) \
    && cd $(PWD)/app/cmd \
    && go build -i -mod=vendor -o $(PROJECT_NAME) \
    && mv $(PROJECT_NAME) $(PWD)/dist/$(PROJECT_NAME) \
    && cp -r $(PWD)/config $(PWD)/dist/$(PROJECT_NAME) \
    && cp -r $(PWD)/templates $(PWD)/dist/$(PROJECT_NAME)\
    && cp -r $(PWD)/static $(PWD)/dist/$(PROJECT_NAME)

docker:
	cd $(PWD)/dist \
	&& docker build -t $(DOCKER_IMAGE):$(PROJECT_TAG) -f $(PWD)/infra/docker/$(PROJECT_NAME)/Dockerfile . \
	&& docker push $(DOCKER_IMAGE):$(PROJECT_TAG) \
	&& docker rmi ${DOCKER_REGISTRY}/${SERVICE}:${TAG}

run:
	docker run --network host hbliveira/doc:dev

mongo:
	docker run -d --network host mvertes/alpine-mongo:4.0.6-1