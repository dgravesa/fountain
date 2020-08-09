#
build:
	go build ./...

deploy:
	gcloud functions deploy waterlogs \
		--region us-east1 \
		--runtime go113 \
		--entry-point FountainPost \
		--trigger-http \
		--allow-unauthenticated \
		--timeout 10 \
		--memory 128MB \

undeploy:
	gcloud functions delete --quiet --region us-east1 waterlogs

