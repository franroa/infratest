package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/require"
	"testing"
)

//func IsPublicSubnetIp4(t *testing.T, subnetId string, region string, cidrRange string) bool {
//	isPublic, err := IsPublicSubnetIp4E(t, subnetId, region, cidrRange)
//	require.NoError(t, err)
//	return isPublic
//}
//
//func IsPublicSubnetIp4E(t *testing.T, subnetId string, region string, cidrRange string) (bool, error) {
//	subnetIdFilterName := "association.subnet-id"
//
//	subnetIdFilter := ec2.Filter{
//		Name:   &subnetIdFilterName,
//		Values: []*string{&subnetId},
//	}
//
//	client, err := NewEc2ClientE(t, region)
//	if err != nil {
//		return false, err
//	}
//
//	rts, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{Filters: []*ec2.Filter{&subnetIdFilter}})
//	if err != nil {
//		return false, err
//	}
//
//	for _, rt := range rts.RouteTables {
//		for _, r := range rt.Routes {
//			if strings.HasPrefix(aws.StringValue(r.GatewayId), "igw-") && aws.StringValue(r.DestinationCidrBlock) == cidrRange {
//				return true, nil
//			}
//		}
//	}
//
//	return false, nil
//}

func GetRouteTablesFromSubnets(t *testing.T, subnetId string, region string) []*ec2.RouteTable {
	routeTables, err := GetRouteTablesFromSubnetsE(t, subnetId, region)
	require.NoError(t, err)
	return routeTables
}

func GetRouteTablesFromSubnetsE(t *testing.T, subnetId string, region string) ([]*ec2.RouteTable, error) {
	subnetIdFilterName := "association.subnet-id"

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	subnetIdFilter := ec2.Filter{Name: &subnetIdFilterName, Values: []*string{&subnetId}}

	routeTables, err := GetRoutesTablesFromFilter(err, client, subnetIdFilter)
	require.NoError(t, err)

	return routeTables, nil
}

func GetRoutesTablesFromFilter(err error, client *ec2.EC2, filter ec2.Filter) ([]*ec2.RouteTable, error) {
	rts, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{Filters: []*ec2.Filter{&filter}})
	if err != nil {
		return nil, err
	}
	return rts.RouteTables, nil
}
