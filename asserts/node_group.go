package asserts

import (
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/stretchr/testify/assert"
	aws "github.com/franroa/infratest/aws"
	"testing"
)

type AssertNodeGroup struct {
	awsNodeGroup      *eks.Nodegroup
	roles             map[string]*AssertRole
	securityGroupList map[string]*AssertSecurityGroup
	Cluster           *AssertEksCluster
	subnetIds         []string
	t                 *testing.T
	region            string
}

func GetEksGroupNode(t *testing.T, clusterName string, groupNodeName string, region string) *AssertNodeGroup {
	awsNodeGroup := aws.GetEksGroupNodeByName(t, clusterName, groupNodeName, region)

	return &AssertNodeGroup{
		awsNodeGroup:      awsNodeGroup,
		roles:             make(map[string]*AssertRole),
		securityGroupList: make(map[string]*AssertSecurityGroup),
		subnetIds:         StringValueSlice(awsNodeGroup.Subnets),
		t:                 t,
		region:            region,
	}
}

func (nodeGroup AssertNodeGroup) HasAmiType(amiType string) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn(amiType, nodeGroup.GetAmiType())
}

func (nodeGroup AssertNodeGroup) HasDiskSize(diskSize int64) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn(diskSize, nodeGroup.GetDiskSize())
}

func (nodeGroup AssertNodeGroup) HasInstanceType(instanceType string) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn([]string{instanceType}, nodeGroup.GetInstanceTypeList())
}

func (nodeGroup AssertNodeGroup) HasExactlyLabel(key string, value string) {
	nodeGroup.HasExactlyLabels(map[string]string{key: value})
}

func (nodeGroup AssertNodeGroup) HasExactlyLabels(labels map[string]string) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn(labels, nodeGroup.GetLabels())
}

func (nodeGroup AssertNodeGroup) HasDesiredSize(size int64) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn(size, nodeGroup.GetDesiredSize())
}

func (nodeGroup AssertNodeGroup) HasMaxSize(size int64) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn(size, nodeGroup.GetMaxSize())
}

func (nodeGroup AssertNodeGroup) HasMinSize(size int64) AssertNodeGroup {
	return nodeGroup.assertEqualAndReturn(size, nodeGroup.GetMinSize())
}

func (nodeGroup AssertNodeGroup) HasRole(roleName string) AssertNodeGroup {
	assert.NotEmpty(nodeGroup.t, nodeGroup.Role(roleName))
	return nodeGroup
}

func (nodeGroup AssertNodeGroup) HasAttachedAllSubnetsFromTheClusterVpc() AssertNodeGroup {
	nodeGroup.Cluster.Vpc().Subnets().HaveExactlyTheIds(nodeGroup.subnetIds)
	return nodeGroup
}

func (nodeGroup AssertNodeGroup) Role(roleName string) *AssertRole {
	if nodeGroup.roles[roleName] == nil {
		nodeGroup.roles[roleName] = GetRoleByName(nodeGroup.t, roleName, nodeGroup.region)
	}

	return nodeGroup.roles[roleName]
}

func (nodeGroup AssertNodeGroup) SecurityGroup(securityGroupName string) *AssertSecurityGroup {
	if nodeGroup.securityGroupList[securityGroupName] == nil {
		nodeGroup.securityGroupList[securityGroupName] = GetSecurityGroupsByName(nodeGroup.t, securityGroupName, nodeGroup.region)
	}

	return nodeGroup.securityGroupList[securityGroupName]
}

// GETTERS

func (nodeGroup AssertNodeGroup) GetAmiType() string {
	return StringValue(nodeGroup.awsNodeGroup.AmiType)
}

func (nodeGroup AssertNodeGroup) GetDiskSize() int64 {
	return Int64Value(nodeGroup.awsNodeGroup.DiskSize)
}

func (nodeGroup AssertNodeGroup) GetInstanceTypeList() []string {
	return StringValueSlice(nodeGroup.awsNodeGroup.InstanceTypes)
}

func (nodeGroup AssertNodeGroup) GetLabels() map[string]string {
	return StringValueMap(nodeGroup.awsNodeGroup.Labels)
}

func (nodeGroup AssertNodeGroup) GetDesiredSize() int64 {
	return Int64Value(nodeGroup.awsNodeGroup.ScalingConfig.DesiredSize)
}

func (nodeGroup AssertNodeGroup) GetMaxSize() int64 {
	return Int64Value(nodeGroup.awsNodeGroup.ScalingConfig.MaxSize)
}

func (nodeGroup AssertNodeGroup) GetMinSize() int64 {
	return Int64Value(nodeGroup.awsNodeGroup.ScalingConfig.MinSize)
}

func (nodeGroup AssertNodeGroup) assertEqualAndReturn(expected, actual interface{}) AssertNodeGroup {
	assert.Equal(nodeGroup.t, expected, actual)

	return nodeGroup
}
