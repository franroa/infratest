package asserts

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"

	aws "github.com/franroa/infratest/aws"
	"testing"
)

type AssertSubnetList struct {
	list   []AssertSubnet
	t      *testing.T
	region string
}

func GetSortedTestSubnetsFromVpc(t *testing.T, vpcId string, region string) *AssertSubnetList {
	subnetsList := AssertSubnetList{t: t, region: region}

	for _, subnet := range aws.GetSubnetsForVpc(t, vpcId, region) {
		subnetsList.list = append(subnetsList.list, *NewSubnet(subnet, t, region))
	}

	sort.Sort(subnetsList)
	return &subnetsList
}

func (assertSubnetList AssertSubnetList) Len() int { return len(assertSubnetList.list) }
func (assertSubnetList AssertSubnetList) Less(i, j int) bool {
	return assertSubnetList.GetCidrBlock(i) < assertSubnetList.GetCidrBlock(j)
}
func (assertSubnetList AssertSubnetList) Swap(i, j int) {
	assertSubnetList.list[i], assertSubnetList.list[j] = assertSubnetList.list[j], assertSubnetList.list[i]
}

func (assertSubnetList AssertSubnetList) ExistInEveryAvailabilityZones() {
	availabilityZones := aws.GetAvailabilityZones(assertSubnetList.t, assertSubnetList.region)

	for _, zone := range availabilityZones {
		var wasFound = false

		for _, assertSubnet := range assertSubnetList.list {
			if assertSubnet.AvailabilityZone == zone {
				wasFound = true
				break
			}
		}

		if wasFound == false {
			assert.Fail(assertSubnetList.t, "Subnet was expected to be on the availability zone '"+zone+"', but was not")
		}
	}
}

func (assertSubnetList AssertSubnetList) ExistOnlyOnceInEachAvailabilityZone() {
	assertSubnetList.ExistInEveryAvailabilityZones()

	if len(assertSubnetList.list) != len(aws.GetAvailabilityZones(assertSubnetList.t, assertSubnetList.region)) {
		assert.Fail(assertSubnetList.t, "Subnet expected to exists only once per availability zone")
	}
}

func (assertSubnetList AssertSubnetList) IdsAreReturnedAsOutputs(t *testing.T, opts *terraform.Options, outputName string) {
	terraformSubnets := terraform.Output(t, opts, outputName)

	for _, assertSubnet := range assertSubnetList.list {
		if !strings.Contains(terraformSubnets, assertSubnet.SubnetId) {
			assert.Fail(t, "Subnet "+assertSubnet.SubnetId+" was expected to be found in the terraform output, but it was not")
		}
	}
}

func (assertSubnetList AssertSubnetList) HaveExactlyTheIds(subnetIds []string) {
	for _, subnetInConfigId := range subnetIds {
		wasFound := false

		for _, subnetInVpc := range assertSubnetList.list {
			if subnetInVpc.SubnetId == subnetInConfigId {
				wasFound = true
			}
		}

		if wasFound == false {
			assert.Fail(assertSubnetList.t, "Subnet with id "+subnetInConfigId+" in the cluster configuration was expected to be found in the given Vpc")
		}
	}

	if len(subnetIds) > assertSubnetList.Len() {
		assert.Fail(assertSubnetList.t, "The subnet config for the cluster has more subnet assigned to it than the total count of subnetIds in the vpc")
	}
}

func (assertSubnetList AssertSubnetList) GetList() []AssertSubnet {
	return assertSubnetList.list
}

func (assertSubnetList AssertSubnetList) GetCidrBlock(i int) string {
	return StringValue(assertSubnetList.list[i].awsSubnet.CidrBlock)
}
