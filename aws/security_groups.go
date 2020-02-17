package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

// TODO lazyLoad on aws sdk

var securityGroupNameFilterName = "group-name"

func GetSecurityGroupsByName(t *testing.T, securityGroupName string, region string) *ec2.SecurityGroup {
	securityGroup, err := GetSecurityGroupsByNameE(t, securityGroupName, region)
	require.NoError(t, err)
	return securityGroup
}

func GetSecurityGroupsByNameE(t *testing.T, securityGroupName string, region string) (*ec2.SecurityGroup, error) {
	groupNameFilter := ec2.Filter{Name: &securityGroupNameFilterName, Values: []*string{&securityGroupName}}

	securityGroups, err := GetSecurityGroupsE(t, []*ec2.Filter{&groupNameFilter}, region)

	numSecurityGroups := len(securityGroups)
	if numSecurityGroups != 1 {
		return nil, fmt.Errorf("Expected to find one SECURITY GROUP with name %s in region %s but found %s", securityGroupName, region, strconv.Itoa(numSecurityGroups))
	}

	return securityGroups[0], err
}

func GetSecurityGroupsE(t *testing.T, filters []*ec2.Filter, region string) ([]*ec2.SecurityGroup, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	describeSecurityGroupsOutput, err := client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	retVal := make([]*ec2.SecurityGroup, len(describeSecurityGroupsOutput.SecurityGroups))

	for i, securityGroup := range describeSecurityGroupsOutput.SecurityGroups {
		retVal[i] = securityGroup
	}

	return retVal, nil
}
