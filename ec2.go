package main

import (
	"context"
	"errors"

	"github.com/st-phuongvu/st-aws-slack-bot/model"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Handler struct {
	UserID   string
	Action   string
	Resource model.AWSResource
}

type EC2Response struct {
	InstancesState map[string]bool
}

const (
	EC2_REBOOT_ACTION = "reboot"
)

var EC2_ALLOWED_ACTION = map[string]bool{
	"reboot":   true,
	"describe": true,
	"list":     true,
	"start":    true,
}

func NewEC2Handler(id, action string, resource *model.AWSResource) (*EC2Response, error) {
	if !EC2_ALLOWED_ACTION[action] {
		return nil, errors.New("No action")
	}
	handler := &EC2Handler{
		UserID:   id,
		Action:   action,
		Resource: *resource,
	}
	return handler.Handle()
	// return nil
}

func (handler *EC2Handler) Handle() (*EC2Response, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch handler.Action {
	case "describe":
		return handler.DescribeInstance(ctx)
	case "start":
		return handler.StartInstance(ctx)
	case "reboot":
		return handler.RebootInstance(ctx)
	default:
		return nil, errors.New("NOT SUPPORT EC2 ACTION")
	}
	return nil, nil
}

func (handler *EC2Handler) StartInstance(c context.Context) (*EC2Response, error) {
	creds := credentials.NewStaticCredentialsProvider(handler.Resource.AccessKey, handler.Resource.SecretKey, "")

	cfg, err := config.LoadDefaultConfig(c, config.WithCredentialsProvider(creds), config.WithRegion("us-east-1"))

	client := ec2.NewFromConfig(cfg)
	ids := []string{handler.Resource.ID}
	input := &ec2.StartInstancesInput{
		InstanceIds: ids,
	}
	resp := &EC2Response{
		InstancesState: map[string]bool{},
	}

	result, err := client.StartInstances(c, input)
	if err != nil {
		return resp, err
	}

	for _, r := range result.StartingInstances {
		resp.InstancesState[*r.InstanceId] = true
	}
	for _, v := range ids {
		if !resp.InstancesState[v] {
			resp.InstancesState[v] = false
		}
	}
	return resp, nil
}

func (handler *EC2Handler) RebootInstance(c context.Context) (*EC2Response, error) {
	creds := credentials.NewStaticCredentialsProvider(handler.Resource.AccessKey, handler.Resource.SecretKey, "")

	cfg, err := config.LoadDefaultConfig(c, config.WithCredentialsProvider(creds), config.WithRegion("us-east-1"))

	client := ec2.NewFromConfig(cfg)
	ids := []string{handler.Resource.ID}
	input := &ec2.RebootInstancesInput{
		InstanceIds: ids,
	}
	resp := &EC2Response{
		InstancesState: map[string]bool{},
	}

	_, err = client.RebootInstances(c, input)
	if err != nil {
		return resp, err
	}

	resp.InstancesState[handler.Resource.ID] = true

	return resp, nil
}

func (handler *EC2Handler) DescribeInstance(c context.Context) (*EC2Response, error) {
	creds := credentials.NewStaticCredentialsProvider(handler.Resource.AccessKey, handler.Resource.SecretKey, "")

	cfg, err := config.LoadDefaultConfig(c, config.WithCredentialsProvider(creds), config.WithRegion("us-east-1"))

	client := ec2.NewFromConfig(cfg)
	ids := []string{handler.Resource.ID}
	input := &ec2.DescribeInstancesInput{
		InstanceIds: ids,
	}
	resp := &EC2Response{
		InstancesState: map[string]bool{},
	}

	result, err := client.DescribeInstances(c, input)
	if err != nil {
		return resp, err
	}

	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			resp.InstancesState[*i.InstanceId] = true
		}
	}
	for _, v := range ids {
		if !resp.InstancesState[v] {
			resp.InstancesState[v] = false
		}
	}
	return resp, nil
}
