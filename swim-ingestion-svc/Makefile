.PHONY: predix predix-test
predix-dev:
	godep save
	godep go install
	cf push -n aviation-ia-swim-ingestion-svc-dev

predix-test:
	curl https://aviation-ia-swim-ingestion-svc-dev.run.asv-pr.ice.predix.io/api/v1/info
	curl https://aviation-ia-swim-ingestion-svc-dev.run.asv-pr.ice.predix.io/api/v1/ping
	# go run $(GOPATH)/src/github.build.ge.com/aviation-intelligent-airport/swim-ingestion-svc/test/up-down.go -basepath=https://aviation-ia-configuration-manager-svc-dev.run.asv-pr.ice.predix.io/api/v1
