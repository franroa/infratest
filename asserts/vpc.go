package asserts

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	aws "github.com/franroa/infratest/aws"
	"testing"
)

type AssertVpc struct {
	awsVpc           *ec2.Vpc
	Id               string
	Cidr             string
	assertSubnetList *AssertSubnetList
	t                *testing.T
	region           string
}

func GetAwsVpcById(t *testing.T, vpcId string, region string) *AssertVpc {
	awsVpc := aws.GetVpcById(t, vpcId, region)

	return &AssertVpc{
		awsVpc: awsVpc,
		Id:     StringValue(awsVpc.VpcId),
		Cidr:   StringValue(awsVpc.CidrBlock),
		t:      t,
		region: region,
	}
}

func GetTestVpcByOutput(t *testing.T, opts *terraform.Options, vpcId string, awsRegion string) *AssertVpc {
	return GetAwsVpcById(t, terraform.Output(t, opts, vpcId), awsRegion)
}

func (assertVpc *AssertVpc) HasCidrRange(cidr string) *AssertVpc {
	assert.Equal(assertVpc.t, cidr, assertVpc.Cidr)

	return assertVpc
}

func (assertVpc *AssertVpc) IsCalled(name string) *AssertVpc {
	assert.Equal(assertVpc.t, name, assertVpc.Tag("Name"))

	return assertVpc
}

func (assertVpc *AssertVpc) HasKubernetesTagWithClusterName(clusterName string) *AssertVpc {
	assertVpc.HasKeyValueTag("kubernetes.io/cluster/"+clusterName, "shared")

	return assertVpc
}

func (assertVpc *AssertVpc) HasKeyValueTag(key string, value string) {
	if assertVpc.Tag(key) != value {
		assert.Fail(assertVpc.t, "Key-Value pair '"+key+"-"+value+"' was expected, but was not found")
	}
}

func (assertVpc *AssertVpc) Subnets() *AssertSubnetList {
	if assertVpc.assertSubnetList == nil {
		assertVpc.assertSubnetList = GetSortedTestSubnetsFromVpc(assertVpc.t, assertVpc.Id, assertVpc.region)
	}

	return assertVpc.assertSubnetList
}

func (assertVpc *AssertVpc) Tag(tagKey string) string {
	for _, tag := range assertVpc.awsVpc.Tags {
		if *tag.Key == tagKey {
			return *tag.Value
		}
	}

	return ""
}
