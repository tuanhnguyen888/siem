package aws_s3sqs

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tuanhnguyen888/siem/collector/common"
	"io"
	"net/http"
	"strings"
)

func (c *ConfigInput) ProcessS3Object(ctx context.Context, obj S3EventV2, msgChan chan<- []map[string]interface{}) error {
	// Request object (download).
	body, err := c.download(ctx, obj)
	if err != nil {
		return errors.Wrap(err, "failed to get s3 object")
	}
	defer body.Close()
	//p.s3Metadata = meta

	reader, err := c.AddGzipDecoderIfNeeded(body)
	if err != nil {
		return errors.Wrap(err, "failed checking for gzip content")
	}

	switch c.LogType {
	case "cloudtrail", "route53":
		if err = c.ReadJSON(ctx, reader, msgChan, obj); err != nil {
			return err
		}
	case "vpc_flow_logs", "cloudfront":
		if err = c.ReadFileLog(ctx, reader, msgChan, obj); err != nil {
			return err
		}
	default:
		return errors.New("Module name is incorrect")
	}
	return nil
}

func (c *ConfigInput) download(ctx context.Context, obj S3EventV2) (body io.ReadCloser, err error) {
	resp, err := c.GetObject(ctx, obj.S3.Bucket.Name, obj.S3.Object.Key)
	if err != nil {
		return nil, err
	}
	//meta := s3Metadata(resp, p.readerConfig.IncludeS3Metadata...)
	return resp.Body, nil
}

func (c *ConfigInput) GetObject(ctx context.Context, bucket, key string) (*s3.GetObjectResponse, error) {
	req := c.S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	resp, err := req.Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject failed: %w", err)
	}

	return resp, nil
}

func (c *ConfigInput) AddGzipDecoderIfNeeded(body io.Reader) (io.Reader, error) {
	bufReader := bufio.NewReader(body)

	gzipped, err := isStreamGzipped(bufReader)
	if err != nil {
		return nil, err
	}
	if !gzipped {
		return bufReader, nil
	}

	return gzip.NewReader(bufReader)
}

func (c *ConfigInput) ReadJSON(ctx context.Context, r io.Reader, msgChan chan<- []map[string]interface{}, obj S3EventV2) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	// encode messages to bytes
	for dec.More() && ctx.Err() == nil {
		var item json.RawMessage
		if err := dec.Decode(&item); err != nil {
			return fmt.Errorf("failed to decode json: %w", err)
		}

		data, _ := item.MarshalJSON()
		if err := c.CreateEvent(string(data), msgChan, obj); err != nil {
			return err
		}
	}

	return nil
}

func (c *ConfigInput) ReadFileLog(ctx context.Context, reader io.Reader, msgChan chan<- []map[string]interface{}, obj S3EventV2) error {
	scanner := bufio.NewScanner(reader)
	if scanner.Scan() {
		// Skip first line of file log

		var batchEvents []map[string]interface{}
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") {
				evt := createEvent(line, c.QueueUrl, c.LogType)
				batchEvents = append(batchEvents, evt)
				//fmt.Println("data ", ": ", line)
			}
		}
		if len(batchEvents) > 0 {
			 common.PushEvtToChan(batchEvents, msgChan)
			logrus.Infof("Push success %v event AWS to CyM chan, key: %v, url: %v", len(batchEvents), obj.S3.Object.Key, c.QueueUrl)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(errors.New(err.Error()))
	}
	return nil
}

func isStreamGzipped(r *bufio.Reader) (bool, error) {
	// Why 512? See https://godoc.org/net/http#DetectContentType
	buf, err := r.Peek(512)
	if err != nil && err != io.EOF {
		return false, err
	}

	switch http.DetectContentType(buf) {
	case "application/x-gzip", "application/zip":
		return true, nil
	default:
		return false, nil
	}
}

//func s3ObjectHash(obj s3EventV2) string {
//	h := sha256.New()
//	h.Write([]byte(obj.S3.Bucket.ARN))
//	h.Write([]byte(obj.S3.Object.Key))
//	prefix := hex.EncodeToString(h.Sum(nil))
//	return prefix[:10]
//}

func (c *ConfigInput) CreateEvent(message string, msgChan chan<- []map[string]interface{},  obj S3EventV2) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		logrus.Error("Error parsing JSON:", err)
		return err
	}
	var batchEvents []map[string]interface{}
	records, ok := data["Records"].([]interface{})
	if !ok {
		event := createEvent(message, c.QueueUrl, c.LogType)
		batchEvents = append(batchEvents, event)
	} else {
		for i, record := range records {
			eventJSON, err := json.Marshal(record)
			if err != nil {
				logrus.Errorf("Error converting event %d, %s to JSON string: %v\n", i+1, string(eventJSON), err)
				continue
			}
			event := createEvent(string(eventJSON), c.QueueUrl, c.LogType)
			batchEvents = append(batchEvents, event)
		}
	}

	if len(batchEvents) > 0 {
		common.PushEvtToChan(batchEvents, msgChan)
		logrus.Infof("Push success %v event AWS to CyM chan, key: %v, url: %v", int64(len(batchEvents)), obj.S3.Object.Key, c.QueueUrl)
		batchEvents = []map[string]interface{}{}
	}
	return nil
}

func createEvent(message string, urlQueue string, provider string) map[string]interface{} {
	id, _ := uuid.GenerateUUID()
	return map[string]interface{}{
		"id": id,
		"log_raw":   message,
		"queue_url": urlQueue,
		"provider":  provider,
	}
}


