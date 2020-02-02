package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type IntegrationTestSuite struct {
	suite.Suite
	cc *grpc.ClientConn
}

func (i *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	cc, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithBlock(), grpc.WithInsecure())
	i.Require().NoError(err)
	i.cc = cc
}

func (i *IntegrationTestSuite) TearDownSuite() {
	i.Require().NoError(i.cc.Close())
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
