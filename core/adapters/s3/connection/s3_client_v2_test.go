package connection

import (
	"net"
	"os"
	"strings"
	"testing"
)

func TestConnectionV2(t *testing.T) {
	os.Setenv("environment", "test")
	s3LocalHost := os.Getenv("AWS_SERVERLESS_S3_LOCAL_HOST")

	s3LocalHostWithHttp := strings.Replace(s3LocalHost, "http://", "", 1)
	_, err := net.Dial("tcp", s3LocalHostWithHttp)

	if err != nil {
		t.Errorf("Expected server %s listen but got %v", s3LocalHost, err)
	}
}
