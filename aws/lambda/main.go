package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	msg := fmt.Sprintf("Hello %s!", name.Name)

	if err := outputCredential(); err != nil {
		return "", err
	}

	if existsCredential() {
		if err := publishPubSub(ctx, msg); err != nil {
			return "", err
		}
	} else {
		log.Print("Google application credential does not exist")
	}

	return msg, nil
}

func outputCredential() error {
	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file; %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(configJSON); err != nil {
		return fmt.Errorf("failed to write string; %w", err)
	}
	return nil
}

func existsCredential() bool {
	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	_, err := os.Stat(path)
	if err == nil {
		log.Printf("Found: %s", path)
		return true
	} else {
		log.Printf("Not found: %s", path)
		return false
	}
}

func publishPubSub(ctx context.Context, msg string) error {
	projectID := os.Getenv("PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	clt, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to create Pub/Sub client; %w", err)
	}
	topic := clt.Topic(topicName)
	resp := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	messageID, err := resp.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message; %w", err)
	}
	log.Printf("message ID: %s", messageID)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

const configJSON = `{
	"type": "external_account",
	"audience": "//iam.googleapis.com/projects/366817278570/locations/global/workloadIdentityPools/awspool/providers/aws0",
	"subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
	"service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/workloaduser@ideal-pancake-380204.iam.gserviceaccount.com:generateAccessToken",
	"token_url": "https://sts.googleapis.com/v1/token",
	"credential_source": {
	  "environment_id": "aws1",
	  "region_url": "http://169.254.169.254/latest/meta-data/placement/availability-zone",
	  "url": "http://169.254.169.254/latest/meta-data/iam/security-credentials",
	  "regional_cred_verification_url": "https://sts.{region}.amazonaws.com?Action=GetCallerIdentity&Version=2011-06-15"
	}
  }`
