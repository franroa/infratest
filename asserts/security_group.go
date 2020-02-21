package asserts

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/franroa/infratest/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AssertSecurityGroup struct {
	AwsSecurityGroup *ec2.SecurityGroup
	egressRulesList  map[int64]*AssertSecurityGroupRule
	igressRulesList  map[int64]*AssertSecurityGroupRule
	GroupId          string
	VpcId            string
	Description      string
	t                *testing.T
	region           string
}

func GetSecurityGroupsByName(t *testing.T, securityGroupName string, region string) *AssertSecurityGroup {
	awsSecurityGroup := aws.GetSecurityGroupsByName(t, securityGroupName, region)

	return &AssertSecurityGroup{
		AwsSecurityGroup: awsSecurityGroup,
		VpcId:            StringValue(awsSecurityGroup.VpcId),
		Description:      StringValue(awsSecurityGroup.Description),
		GroupId:          StringValue(awsSecurityGroup.GroupId),
		t:                t,
		region:           region,
	}
}

func (group *AssertSecurityGroup) IsCalled(name string) *AssertSecurityGroup {
	assert.Equal(group.t, name, group.Tag("Name"))

	return group
}

func (group *AssertSecurityGroup) IsAttachedToVpc(vpcId string) *AssertSecurityGroup {
	assert.Equal(group.t, vpcId, group.VpcId)

	return group
}

func (group *AssertSecurityGroup) IsUsedFor(description string) *AssertSecurityGroup {
	assert.Equal(group.t, description, group.Description)

	return group
}

func (group AssertSecurityGroup) IngressRulesWithPortFrom(port int64) *AssertSecurityGroupRule {
	if group.igressRulesList == nil || group.igressRulesList[port] == nil {
		group.igressRulesList = group.MapIngressRules(group.AwsSecurityGroup.IpPermissions)
	}

	return group.igressRulesList[port]
}

func (group AssertSecurityGroup) EgressRulesWithPortFrom(port int64) *AssertSecurityGroupRule {
	if group.egressRulesList == nil || group.egressRulesList[port] == nil {
		group.egressRulesList = group.MapIngressRules(group.AwsSecurityGroup.IpPermissionsEgress)
	}

	return group.egressRulesList[port]
}

func (group AssertSecurityGroup) MapIngressRules(permissions []*ec2.IpPermission) map[int64]*AssertSecurityGroupRule {
	rules := make(map[int64]*AssertSecurityGroupRule)

	for _, rule := range permissions {
		rules[Int64Value(rule.FromPort)] = &AssertSecurityGroupRule{
			assertRule:            rule,
			ToPort:                Int64Value(rule.ToPort),
			IpProtocol:            StringValue(rule.IpProtocol),
			Description:           group.GetRuleDescription(rule),
			SourceSecurityGroupId: group.GetRuleSourceSecurityGroupId(rule),
			t:                     group.t,
		}
	}

	return rules
}

func (group AssertSecurityGroup) GetRuleDescription(rule *ec2.IpPermission) string {
	if rule.UserIdGroupPairs == nil {
		return ""
	}

	return StringValue(rule.UserIdGroupPairs[0].Description)
}

func (group AssertSecurityGroup) GetRuleSourceSecurityGroupId(rule *ec2.IpPermission) string {
	if rule.UserIdGroupPairs == nil {
		return ""
	}

	return StringValue(rule.UserIdGroupPairs[0].GroupId)
}

func (group AssertSecurityGroup) Tag(tagKey string) string {
	for _, tag := range group.AwsSecurityGroup.Tags {
		if *tag.Key == tagKey {
			return *tag.Value
		}
	}

	return ""
}

func (group *AssertSecurityGroup) AllowsAllIngressTrafficFromItself(description string) *AssertSecurityGroup {
	group.IngressRulesWithPortFrom(0).
		HasSourceSecurityGroup(group).
		HasIpProtocol("-1").
		HasDescription(description).
		HasPortTo(0)

	return group
}

func (group *AssertSecurityGroup) AllowsAllOutboundTrafficToEverywhere() *AssertSecurityGroup {
	group.EgressRulesWithPortFrom(0).
		HasIpProtocol("-1").
		HasCidrIp("0.0.0.0/0").
		HasPortTo(0)

	return group
}

func (group *AssertSecurityGroup) AllowSshIngressTrafficFrom(securityGroup *AssertSecurityGroup, description string) *AssertSecurityGroup {
	group.IngressRulesWithPortFrom(443).
		HasSourceSecurityGroup(securityGroup).
		HasPortTo(443).
		HasIpProtocol("tcp").
		HasDescription(description)

	return group
}

func (group *AssertSecurityGroup) AllowsTcpIngressTrafficFromEphemeralPortsCommingFrom(sourceGroup *AssertSecurityGroup, description string) *AssertSecurityGroup {
	group.IngressRulesWithPortFrom(1025).
		HasSourceSecurityGroup(sourceGroup).
		HasIpProtocol("tcp").
		HasDescription(description).
		HasPortTo(65535)

	return group
}

type AssertSecurityGroupRule struct {
	assertRule            *ec2.IpPermission
	ToPort                int64
	IpProtocol            string
	Description           string
	SourceSecurityGroupId string
	t                     *testing.T
}

func (rule AssertSecurityGroupRule) HasPortTo(port int64) AssertSecurityGroupRule {
	assert.Equal(rule.t, port, rule.ToPort)

	return rule
}

func (rule AssertSecurityGroupRule) HasIpProtocol(protocol string) AssertSecurityGroupRule {
	assert.Equal(rule.t, protocol, rule.IpProtocol)

	return rule
}

func (rule AssertSecurityGroupRule) HasDescription(description string) AssertSecurityGroupRule {
	assert.Equal(rule.t, description, rule.Description)

	return rule
}

func (rule AssertSecurityGroupRule) HasCidrIp(cidrRange string) AssertSecurityGroupRule {
	for _, ipRange := range rule.assertRule.IpRanges {
		if StringValue(ipRange.CidrIp) == cidrRange {
			return rule
		}
	}

	assert.Fail(rule.t, "The cidr range "+cidrRange+" was expected, but not found")
	return rule
}

func (rule AssertSecurityGroupRule) HasSourceSecurityGroup(group *AssertSecurityGroup) AssertSecurityGroupRule {
	assert.Equal(rule.t, group.GroupId, rule.SourceSecurityGroupId)

	return rule
}
