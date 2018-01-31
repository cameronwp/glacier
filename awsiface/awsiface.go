package awsiface

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
)

// https://github.com/aws/aws-sdk-go/blob/master/service/glacier/glacieriface/interface.go#L28

func InitiateMultiPartUpload(svc glacieriface.GlacierAPI, fileName string, partSize int, vault string) error {
	initResult, err := svc.InitiateMultipartUpload(&glacier.InitiateMultipartUploadInput{
		AccountId:          aws.String("-"),
		ArchiveDescription: aws.String(fileName),
		PartSize:           aws.String(fmt.Sprintf("%d", partSize)),
		VaultName:          aws.String(vault),
	})
	if err != nil {
		return formatAWSError(err)
	}
}

func formatAWSError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case glacier.ErrCodeResourceNotFoundException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeResourceNotFoundException, aerr.Error())
		case glacier.ErrCodeInvalidParameterValueException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeInvalidParameterValueException, aerr.Error())
		case glacier.ErrCodeMissingParameterValueException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeMissingParameterValueException, aerr.Error())
		case glacier.ErrCodeRequestTimeoutException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeRequestTimeoutException, aerr.Error())
		case glacier.ErrCodeServiceUnavailableException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeServiceUnavailableException, aerr.Error())
		default:
			return fmt.Errorf("%s", aerr.Error())
		}
	}
	return fmt.Errorf(err.Error())
}
