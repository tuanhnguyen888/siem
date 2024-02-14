package aws_s3sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type canceler interface {
	Done() <-chan struct{}
	Err() error
}

// sqs interface

type ISQSClient interface {
	ReceiveMessageRequest(input *sqs.ReceiveMessageInput) IReceiveMessageRequest
	DeleteMessageRequest(input *sqs.DeleteMessageInput) IDeleteMessageRequest
	ChangeMessageVisibilityRequest(input *sqs.ChangeMessageVisibilityInput) IChangeMessageVisibilityRequest
}

type SQSClient struct {
	sqsClient *sqs.Client
}

type IReceiveMessageRequest interface {
	Send(ctx context.Context) (*sqs.ReceiveMessageResponse, error)
}
type ReceiveMessageRequest struct {
	req sqs.ReceiveMessageRequest
}

type IDeleteMessageRequest interface {
	Send(ctx context.Context) (*sqs.DeleteMessageResponse, error)
}
type DeleteMessageRequest struct {
	req sqs.DeleteMessageRequest
}

type IChangeMessageVisibilityRequest interface {
	Send(ctx context.Context) (*sqs.ChangeMessageVisibilityResponse, error)
}
type ChangeMessageVisibilityRequest struct {
	req sqs.ChangeMessageVisibilityRequest
}

func (s *SQSClient) ReceiveMessageRequest(input *sqs.ReceiveMessageInput) IReceiveMessageRequest {
	req := s.sqsClient.ReceiveMessageRequest(input)
	return &ReceiveMessageRequest{req: req}
}

func (s *SQSClient) DeleteMessageRequest(input *sqs.DeleteMessageInput) IDeleteMessageRequest {
	req := s.sqsClient.DeleteMessageRequest(input)
	return &DeleteMessageRequest{req: req}
}

func (s *SQSClient) ChangeMessageVisibilityRequest(input *sqs.ChangeMessageVisibilityInput) IChangeMessageVisibilityRequest {
	req := s.sqsClient.ChangeMessageVisibilityRequest(input)
	return &ChangeMessageVisibilityRequest{req: req}
}

func (r *ReceiveMessageRequest) Send(ctx context.Context) (*sqs.ReceiveMessageResponse, error) {
	return r.req.Send(ctx)
}

func (r *ChangeMessageVisibilityRequest) Send(ctx context.Context) (*sqs.ChangeMessageVisibilityResponse, error) {
	return r.req.Send(ctx)

}

func (r *DeleteMessageRequest) Send(ctx context.Context) (*sqs.DeleteMessageResponse, error) {
	return r.req.Send(ctx)

}

//s3 client interface

type IS3Client interface {
	GetObjectRequest(input *s3.GetObjectInput) IGetObject
}

type IGetObject interface {
	Send(ctx context.Context) (*s3.GetObjectResponse, error)
}

type S3Client struct {
	s3Client *s3.Client
}

type GetObject struct {
	getObject s3.GetObjectRequest
}

func (s *S3Client) GetObjectRequest(input *s3.GetObjectInput) IGetObject {
	resp := s.s3Client.GetObjectRequest(input)
	return &GetObject{getObject: resp}
}

func (i *GetObject) Send(ctx context.Context) (*s3.GetObjectResponse, error) {
	return i.getObject.Send(ctx)
}
