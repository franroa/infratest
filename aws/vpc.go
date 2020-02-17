package aws

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/require"
)

var vpcIDFilterName = "vpc-id"

func GetVpcById(t *testing.T, vpcId string, region string) *ec2.Vpc {
	vpc, err := GetVpcByIdE(t, vpcId, region)
	require.NoError(t, err)
	return vpc
}

func GetVpcByIdE(t *testing.T, vpcId string, region string) (*ec2.Vpc, error) {
	vpcIdFilter := ec2.Filter{Name: &vpcIDFilterName, Values: []*string{&vpcId}}
	vpcs, err := GetVpcsE(t, []*ec2.Filter{&vpcIdFilter}, region)

	numVpcs := len(vpcs)
	if numVpcs != 1 {
		return nil, fmt.Errorf("Expected to find one VPC with ID %s in region %s but found %s", vpcId, region, strconv.Itoa(numVpcs))
	}

	return vpcs[0], err
}

func GetVpcsE(t *testing.T, filters []*ec2.Filter, region string) ([]*ec2.Vpc, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	describeVpcsOutput, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	retVal := make([]*ec2.Vpc, len(describeVpcsOutput.Vpcs))

	for i, vpc := range describeVpcsOutput.Vpcs {
		retVal[i] = vpc
	}

	return retVal, nil
}
