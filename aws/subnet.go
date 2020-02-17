package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func GetSubnetsForVpc(t *testing.T, vpcID string, region string) []*ec2.Subnet {
	subnets, err := GetSubnetsForVpcE(t, vpcID, region)
	if err != nil {
		t.Fatal(err)
	}
	return subnets
}

func GetSubnetsForVpcE(t *testing.T, vpcID string, region string) ([]*ec2.Subnet, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	vpcIDFilter := ec2.Filter{Name: &vpcIDFilterName, Values: []*string{&vpcID}}
	subnetOutput, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{Filters: []*ec2.Filter{&vpcIDFilter}})
	if err != nil {
		return nil, err
	}

	subnets := []*ec2.Subnet{}

	for _, ec2Subnet := range subnetOutput.Subnets {
		subnets = append(subnets, ec2Subnet)
	}

	return subnets, nil
}
