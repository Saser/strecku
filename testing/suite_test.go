package testing

import (
	"context"
	"fmt"
	"testing"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/strecku/v1"
	testingv1 "github.com/Saser/strecku/backend/gen/api/testing/v1"
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

func (i *IntegrationTestSuite) SetupTest() {
	ctx := context.Background()
	{
		c := streckuv1.NewUserAPIClient(i.cc)
		user := &streckuv1.User{
			DisplayName:  "Saser",
			EmailAddress: "saser@saser.com",
		}
		_, err := c.CreateUser(ctx, &streckuv1.CreateUserRequest{
			User: user,
		})
		i.Require().NoError(err)
	}
}

func (i *IntegrationTestSuite) AfterTest(suiteName, testName string) {
	ctx := context.Background()
	c := testingv1.NewResetAPIClient(i.cc)
	_, err := c.Reset(ctx, &testingv1.ResetRequest{
		Reason: fmt.Sprintf("%s/%s", suiteName, testName),
	})
	i.Require().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
