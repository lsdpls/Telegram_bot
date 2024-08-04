include .env

# Yandex container creation
create:
	yc serverless container create --name $(SERVERLESS_CONTAINER_NAME)
	yc serverless container allow-unauthenticated-invoke --name  $(SERVERLESS_CONTAINER_NAME)

# Yandex gateway api
create_gw_spec:
	$(shell sed "s/SERVERLESS_CONTAINER_ID/${SERVERLESS_CONTAINER_ID}/;s/SERVICE_ACCOUNT_ID/${SERVICE_ACCOUNT_ID}/" api-gw.yaml.example > api-gw.yaml)
create_gw: create_gw_spec
	yc serverless api-gateway create --name $(SERVERLESS_CONTAINER_NAME) --spec api-gw.yaml

# Telegram bot webhook
webhook_info:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/getMe"

webhook_delete:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/deleteWebhook"

webhook_create: webhook_delete
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/setWebhook" --header 'content-type: application/json' --data '{"url": "$(SERVERLESS_APIGW_URL)"}'

# Yandex container deployment
build: 
	docker build -t cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME) .

push: build
	docker push cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME)

deploy: push
	$(shell sed 's/=.*/=/' .env > .env.example)
	yc serverless container revision deploy --container-name $(SERVERLESS_CONTAINER_NAME) --image 'cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME):latest' --service-account-id $(SERVICE_ACCOUNT_ID)  --environment='$(shell tr '\n' ',' < .env | sed 's/\r//g')' --core-fraction 5 --execution-timeout $(SERVERLESS_CONTAINER_EXEC_TIMEOUT)

# fast login
auth:
	docker login \
	--username oauth \
	--password $(YC_PASSWORD) \
	cr.yandex