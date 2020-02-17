package asserts

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AssertSubnet struct {
	awsSubnet        *ec2.Subnet
	AvailabilityZone string
	SubnetId         string
	CidrBlock        string
	assertRouteTable *AssertRouteTable
	t                *testing.T
	region           string
}

func NewSubnet(subnet *ec2.Subnet, t *testing.T, region string) *AssertSubnet {
	return &AssertSubnet{
		awsSubnet:        subnet,
		AvailabilityZone: StringValue(subnet.AvailabilityZone),
		SubnetId:         StringValue(subnet.SubnetId),
		CidrBlock:        StringValue(subnet.CidrBlock),
		t:                t,
		region:           region,
	}
}

func (subnet AssertSubnet) IsCalled(name string) {
	assert.Equal(subnet.t, name, subnet.Tag("Name"))
}

func (subnet AssertSubnet) HasCidrRange(cidr string) {
	assert.Equal(subnet.t, cidr, subnet.CidrBlock)
}

func (subnet AssertSubnet) HasKubernetesTagWithClusterName(clusterName string) {
	subnet.HasKeyValueTag("kubernetes.io/cluster/"+clusterName, "shared")
}

func (subnet *AssertSubnet) HasKeyValueTag(key string, value string) {
	for _, tag := range subnet.awsSubnet.Tags {
		if *tag.Key == key && *tag.Value == value {
			return
		}
	}

	assert.Fail(subnet.t, "Key-Value pair '"+key+"-"+value+"' was expected, but was not found")
}

func (subnet AssertSubnet) IsPublic() {
	subnet.RouteTable().IsPublicWithIp4Range("0.0.0.0/0")
}

func (subnet *AssertSubnet) Tag(tagKey string) string {
	for _, tag := range subnet.awsSubnet.Tags {
		if *tag.Key == tagKey {
			return *tag.Value
		}
	}

	return ""
}

func (subnet AssertSubnet) RouteTable() *AssertRouteTable {
	if subnet.assertRouteTable == nil {
		subnet.assertRouteTable = GetRouteTableFromSubnet(subnet.t, subnet.SubnetId, subnet.region)
	}

	return subnet.assertRouteTable
}
