package asserts

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	aws "github.com/franroa/infratest/aws"
	"testing"
)

type AssertRouteTable struct {
	awsRouteTable []*ec2.RouteTable
	t             *testing.T
	region        string
}

func GetRouteTableFromSubnet(t *testing.T, id string, region string) *AssertRouteTable {
	return &AssertRouteTable{
		awsRouteTable: aws.GetRouteTablesFromSubnets(t, id, region),
		t:             t,
		region:        region,
	}
}

func (routeTable AssertRouteTable) IsPublicWithIp4Range(cidrRange string) {
	for _, rt := range routeTable.awsRouteTable {
		for _, r := range rt.Routes {
			if HasPrefix(StringValue(r.GatewayId), "igw-") && StringValue(r.DestinationCidrBlock) == cidrRange {
				return
			}
		}
	}

	assert.Fail(routeTable.t, "Ip was expected to be public")
}
