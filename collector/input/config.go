package aws_s3sqs

import (

	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tuanhnguyen888/siem/collector/input/credentails"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	ModuleName                          = "aws_s3sqs"
	sqsRetryDelay                       = 10 * time.Second
	sqsApproximateReceiveCountAttribute = "ApproximateReceiveCount"
)

type IInputPlugin interface {
	Start(ctx context.Context,  msgChan chan<- []map[string]interface{}) error
}

// ConfigInput holds the configuration json fields and internal objects
type ConfigInput struct {
	Name string
	QueueUrl            string `mapstructure:"queue_url"`
	LogType             string `mapstructure:"log_type"`
	VisibilityTimeout   int64  `mapstructure:"visibility_timeout"`
	SQSWaitTime         int64  `mapstructure:"wait_time"`
	MaxNumberOfMessages int    `mapstructure:"max_number_of_messages"`
	ApiTimeout          int64  `mapstructure:"api_timeout"`
	SQSMaxReceiveCount  int    `mapstructure:"max_receive_count"`
	AccessKeyId         string `mapstructure:"access_key_id"`
	SecretAccessKey     string `mapstructure:"secret_access_key"`
	ProxyUrl            string `mapstructure:"proxy_url"`
	Endpoint            string `mapstructure:"endpoint"`
	DefaultRegion       string `mapstructure:"default_region"`
	SqsClient           ISQSClient
	S3Client            IS3Client
}

//InitHandler initialize the aws s3 sqs input plugin
func InitHandler(ctx context.Context, input map[string]interface{} ) (IInputPlugin, error) {
	return ProcessHandler(ctx, input)
}

// DefaultInputConfig returns an ConfigInput struct with default values
func DefaultInputConfig() ConfigInput {
	return ConfigInput{
		Name:         ModuleName,
		ApiTimeout:          120,
		VisibilityTimeout:   300,
		SQSWaitTime:         20,
		MaxNumberOfMessages: 5,
		Endpoint:            "amazonaws.com",
		DefaultRegion:       "us-east-1",
	}
}

func ProcessHandler(ctx context.Context, input map[string]interface{}) (IInputPlugin, error) {
	conf := DefaultInputConfig()
	err := mapstructure.Decode(input, &conf)
	if err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		return nil, errors.New(err.Error() + ", sqsUrl: " + conf.QueueUrl)
	}

	logrus.Infof("Successfully initialize aws s3 sqs input. (sqs_queue: %v) ", conf.QueueUrl)
	return &conf, nil
}

func (c *ConfigInput) Validate() error {
	if c.QueueUrl == "" {
		return fmt.Errorf("queue_url is require, input aws will stop")
	}

	if c.LogType == "" {
		return fmt.Errorf("log_type is require, input aws will stop")
	}

	if c.VisibilityTimeout <= int64(0) || c.VisibilityTimeout > int64(43200) {
		return fmt.Errorf("visibility_timeout <%v> must be greater than 0 and "+
			"less than or equal to 12h", c.VisibilityTimeout)
	}

	if c.SQSWaitTime <= int64(0) || c.SQSWaitTime > int64(20) {
		return fmt.Errorf("wait_time <%v> must be greater than 0 and "+
			"less than or equal to 20s", c.SQSWaitTime)
	}

	if c.MaxNumberOfMessages <= 0 {
		return fmt.Errorf("max_number_of_messages <%v> must be greater than 0",
			c.MaxNumberOfMessages)
	}

	if c.ApiTimeout < c.SQSWaitTime {
		return fmt.Errorf("api_timeout <%v> must be greater than the sqs.wait_time <%v ",
			c.ApiTimeout, c.SQSWaitTime)
	}

	return nil
}

func (c *ConfigInput) Run(ctx context.Context, msgChan chan<- []map[string]interface{}) error {

	if err := c.CreateSQSReceiver(); err != nil {
		return errors.New(err.Error() + ", sqsQueue: " + c.QueueUrl)
	}

	if err := c.Receive(ctx, msgChan); err != nil {
		return err
	}

	return nil
}

func (c *ConfigInput) CreateSQSReceiver() error {
	awsConfig, err := credentails.InitializeAWSConfig(c.AccessKeyId, c.SecretAccessKey, c.DefaultRegion, c.ProxyUrl)
	if err != nil {
		return fmt.Errorf("initialize AWS credentials failed: %w", err)
	}

	regionName, err := getRegionFromQueueURL(c.QueueUrl, c.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to get AWS region from queue_url: %w", err)
	}
	awsConfig.Region = regionName

	sqsClient := sqs.New(c.EnrichAWSConfigWithEndpoint(c.Endpoint, "sqs", awsConfig.Region, awsConfig))
	c.SqsClient = &SQSClient{sqsClient: sqsClient}

	s3Client := s3.New(c.EnrichAWSConfigWithEndpoint(c.Endpoint, "s3", awsConfig.Region, awsConfig))
	c.S3Client = &S3Client{s3Client: s3Client}

	return nil
}

func (c *ConfigInput) Receive(ctx context.Context,  msgChan chan<- []map[string]interface{}) error {
	var workerWg sync.WaitGroup
	for ctx.Err() == nil {
		msgs, err := c.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() == nil {
				logrus.Warn("SQS ReceiveMessage returned an error. Will retry after a short delay.", "error", err, "\n sqsUrl: ", c.QueueUrl)
				// Throttle retries.
				Wait(ctx, sqsRetryDelay)
			}
			continue
		}

		if len(msgs) > 0 {
			logrus.Infof(" Received %v SQS messages, sqsUrl: %v", len(msgs), c.QueueUrl)
			workerWg.Add(len(msgs))
			for _, msg := range msgs {
				go func(msg sqs.Message) {
					defer func() {
						workerWg.Done()
					}()
					if err := c.ProcessSQS(ctx, &msg, msgChan); err != nil {
						logrus.Warn("Failed processing SQS message.", "error", err, "\n", "sqsQueue: ", c.QueueUrl, ", message_id", *msg.MessageId)
					}
				}(msg)
			}
		} else {
			logrus.Infof("----  Will retry receive api ---- , sqsUrl: %v", c.QueueUrl)
		}
	}
	workerWg.Wait()

	if errors.Is(ctx.Err(), context.Canceled) {
		// A canceled context is a normal shutdown.
		return nil
	}
	return ctx.Err()
}

func (c *ConfigInput) ReceiveMessage(ctx context.Context) ([]sqs.Message, error) {
	const sqsMaxNumberOfMessagesLimit = 10

	req := c.SqsClient.ReceiveMessageRequest(
		&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.QueueUrl),
			MaxNumberOfMessages: aws.Int64(int64(min(c.MaxNumberOfMessages, sqsMaxNumberOfMessagesLimit))),
			VisibilityTimeout:   aws.Int64(c.VisibilityTimeout),
			WaitTimeSeconds:     aws.Int64(c.SQSWaitTime),
			AttributeNames:      []sqs.QueueAttributeName{sqsApproximateReceiveCountAttribute},
		})
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.ApiTimeout)*time.Second)
	defer cancel()

	resp, err := req.Send(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			err = fmt.Errorf("api_timeout exceeded: %w", err)
		}
		return nil, fmt.Errorf(" sqs ReceiveMessage failed: %w", err)
	}
	return resp.Messages, nil
}

func (c *ConfigInput) Start(ctx context.Context,  msgChan chan<- []map[string]interface{}) error {
	var err error


	go func() {
		defer ctx.Done() // Ensure cancellation function is called when the goroutine exits

		if err := c.Run(ctx, msgChan); err != nil {
			logrus.Error(err, ", sqsUrl: ", c.QueueUrl)
		}
	}()
	// ---------
	<-ctx.Done()
	logrus.Error("CyM_Agent input aws stopped, sqsUrl: ", c.QueueUrl)
	return err
}

func getRegionFromQueueURL(queueURL string, endpoint string) (string, error) {
	// get region from queueURL
	// Example: https://sqs.us-east-1.amazonaws.com/627959692251/test-s3-logs
	urlName, err := url.Parse(queueURL)
	if err != nil {
		return "", fmt.Errorf(queueURL + " is not a valid URL")
	}
	if urlName.Scheme == "https" && urlName.Host != "" {
		queueHostSplit := strings.Split(urlName.Host, ".")
		if len(queueHostSplit) > 2 && (strings.Join(queueHostSplit[2:], ".") == endpoint || (endpoint == "" && queueHostSplit[2] == "amazonaws")) {
			return queueHostSplit[1], nil
		}
	}
	return "", fmt.Errorf("QueueURL is not in format: https://sqs.{REGION_ENDPOINT}.{ENDPOINT}/{ACCOUNT_NUMBER}/{QUEUE_NAME}")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Wait(ctx canceler, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *ConfigInput) EnrichAWSConfigWithEndpoint(endpoint string, serviceName string, regionName string, awsConfig aws.Config) aws.Config {
	if endpoint != "" {
		awsConfig.EndpointResolver = aws.ResolveWithEndpointURL("https://" + serviceName + "." + regionName + "." + endpoint)
	}
	return awsConfig
}
