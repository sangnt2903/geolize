protoc:
	docker build -t protoc -f bootstrap/protoc/Dockerfile bootstrap/protoc

swagger:
	sudo rm -rf /var/lib/swagger-ui && \
	sudo mkdir -p /var/lib/swagger-ui && \
	sudo chmod -R 0775 /var/lib/swagger-ui && \
	sudo cp -r bootstrap/public/swagger-ui/* /var/lib/swagger-ui/

bootstrap: protoc swagger

gen:
ifdef service
	@docker run -it --rm \
	-e SERVICE=$(service) \
	-v ./service-protos/services/$(service):/defs/$(service) \
	-v ./service-protos/includes:/defs/includes \
	-v ./service-protos/generated/:/generated \
	protoc:latest /entrypoint.sh

	@cp service-protos/generated/$(service)/service.swagger.json services/$(service)/swagger.json
else
	@echo "service is required, make gen service=core"
endif