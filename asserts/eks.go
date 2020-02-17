package asserts

import (
	"github.com/aws/aws-sdk-go/service/eks"
	aws "github.com/franroa/infratest/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AssertEksCluster struct {
	awsEksCluster     *eks.Cluster
	Name              string
	Status            string
	SubnetIds         []string
	VpcId             string
	vpc               *AssertVpc
	roles             map[string]*AssertRole
	nodeGroups        map[string]*AssertNodeGroup
	securityGroupList map[string]*AssertSecurityGroup
	region            string
	subnetIds         []string
	t                 *testing.T
}

func GetEksClusterByName(t *testing.T, clusterName string, region string) AssertEksCluster {
	awsCluster := aws.GetEksClusterByName(t, clusterName, region)
	return AssertEksCluster{
		awsEksCluster:     awsCluster,
		Name:              StringValue(awsCluster.Name),
		Status:            StringValue(awsCluster.Status),
		SubnetIds:         StringValueSlice(awsCluster.ResourcesVpcConfig.SubnetIds),
		VpcId:             StringValue(awsCluster.ResourcesVpcConfig.VpcId),
		roles:             make(map[string]*AssertRole),
		nodeGroups:        make(map[string]*AssertNodeGroup),
		securityGroupList: make(map[string]*AssertSecurityGroup),
		t:                 t,
		region:            region,
	}
}

func (assertCluster AssertEksCluster) HasStatus(status string) {
	assert.Equal(assertCluster.t, status, assertCluster.Status)
}

func (assertCluster AssertEksCluster) HasNodegroup(groupNode string) {
	assert.NotEmpty(assertCluster.t, assertCluster.NodeGroup(groupNode))
}

func (assertCluster AssertEksCluster) HasAttachedAllSubnetsFromItsVpc() {
	assertCluster.Vpc().Subnets().HaveExactlyTheIds(assertCluster.SubnetIds)
}

func (assertCluster AssertEksCluster) HasRole(roleName string) {
	assert.NotEmpty(assertCluster.t, assertCluster.Role(roleName))
}

func (assertCluster AssertEksCluster) Role(roleName string) *AssertRole {
	if assertCluster.roles[roleName] == nil {
		assertCluster.roles[roleName] = GetRoleByName(assertCluster.t, roleName, assertCluster.region)
	}

	return assertCluster.roles[roleName]
}

func (assertCluster AssertEksCluster) Vpc() *AssertVpc {
	if assertCluster.vpc == nil {
		assertCluster.vpc = GetAwsVpcById(assertCluster.t, assertCluster.VpcId, assertCluster.region)
	}

	return assertCluster.vpc
}

func (assertCluster AssertEksCluster) NodeGroup(nodeGroupName string) *AssertNodeGroup {
	if assertCluster.nodeGroups[nodeGroupName] == nil {
		assertCluster.nodeGroups[nodeGroupName] = GetEksGroupNode(assertCluster.t, assertCluster.Name, nodeGroupName, assertCluster.region) // TODO can be better
		assertCluster.nodeGroups[nodeGroupName].Cluster = &assertCluster
	}

	return assertCluster.nodeGroups[nodeGroupName]
}

func (assertCluster AssertEksCluster) SecurityGroup(securityGroupName string) *AssertSecurityGroup {
	if assertCluster.securityGroupList[securityGroupName] == nil {
		assertCluster.securityGroupList[securityGroupName] = GetSecurityGroupsByName(assertCluster.t, securityGroupName, assertCluster.region)
	}

	return assertCluster.securityGroupList[securityGroupName]
}
