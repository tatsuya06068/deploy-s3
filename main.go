package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type S3Event struct {
	Records []struct {
		S3 struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func handler(ctx context.Context, event S3Event) (string, error) {
	distributionID := "YOUR_DISTRIBUTION_ID"
	newOriginPath := "/new/directory/path"

	// Initialize CloudFront client
	sess := session.Must(session.NewSession())
	svc := cloudfront.New(sess)

	// Get current distribution config
	resp, err := svc.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: aws.String(distributionID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get distribution config: %v", err)
	}

	// Update origin path
	for _, origin := range resp.DistributionConfig.Origins.Items {
		origin.OriginPath = aws.String(newOriginPath)
	}

	// Update distribution config
	_, err = svc.UpdateDistribution(&cloudfront.UpdateDistributionInput{
		Id:                 aws.String(distributionID),
		IfMatch:            resp.ETag,
		DistributionConfig: resp.DistributionConfig,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update distribution: %v", err)
	}

	return fmt.Sprintf("CloudFront distribution updated successfully: %s", newOriginPath), nil
}

func main() {
	lambda.Start(handler)
}
