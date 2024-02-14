package credentails

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func InitializeAWSConfig(accessKeyId, secretAccessKey, region, proxy string) (aws.Config, error) {
	AWSConfig := getAccessKeys(accessKeyId, secretAccessKey, region)
	//Add ProxyUrl
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			logrus.Errorf("Error parsing proxy URL: %s", err)
			return AWSConfig, err
		}
		logrus.Println("Run AWS Cloudwatch via httpProxy: ", proxyURL)
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
		AWSConfig.HTTPClient = httpClient
	}
	return AWSConfig, nil
}

func getAccessKeys(accessKeyId, secretAccessKey, region string) aws.Config {
	awsConfig := defaults.Config()
	awsCredentials := aws.Credentials{
		AccessKeyID:     accessKeyId,
		SecretAccessKey: secretAccessKey,
	}

	awsConfig.Credentials = aws.StaticCredentialsProvider{
		Value: awsCredentials,
	}

	// Set default region if empty to make initial aws api call
	awsConfig.Region = region

	return awsConfig
}

// EnrichAWSConfigWithEndpoint function enabled endpoint resolver for AWS
// service clients when endpoint is given in config.
func EnrichAWSConfigWithEndpoint(endpoint string, serviceName string, regionName string, awsConfig aws.Config) aws.Config {
	if endpoint != "" {
		if regionName == "" {
			awsConfig.EndpointResolver = aws.ResolveWithEndpointURL("https://" + serviceName + "." + endpoint)
		} else {
			awsConfig.EndpointResolver = aws.ResolveWithEndpointURL("https://" + serviceName + "." + regionName + "." + endpoint)
		}
	}
	return awsConfig
}
