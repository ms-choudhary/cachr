package cachr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var awsSession *Session

// TODO: read from env
func Bucket() string {
	return ""
}

// TODO: read from env
func AWSRegion() aws.String {
	return ""
}

func AWSSession() (*Session, error) {
	if awsSession != nil {
		return awsSession, nil
	}

	awsSession = session.NewSession(&aws.Config{
		Region: AWSRegion(),
	})

	return awsSession, nil
}
