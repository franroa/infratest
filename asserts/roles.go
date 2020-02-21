package asserts

import (
	"github.com/aws/aws-sdk-go/service/iam"
	aws "github.com/franroa/infratest/aws"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

type AssertRole struct {
	role        *iam.Role
	t           *testing.T
	region      string
	roleDetails *iam.RoleDetail
}

func GetRoleByName(t *testing.T, roleName string, region string) *AssertRole {
	return &AssertRole{
		role:   aws.GetRole(t, roleName, region),
		t:      t,
		region: region,
	}
}

func (role AssertRole) AssumesRolePolicyDocument(assumeRoleString string) AssertRole {
	assert.Equal(role.t, url.QueryEscape(assumeRoleString), role.GetAssumedPolicy())

	return role
}

func (role AssertRole) IsAttachedToPolicy(policy string) AssertRole {
	for _, managedPolicy := range role.RoleDetails().AttachedManagedPolicies {
		if StringValue(managedPolicy.PolicyName) == policy {
			return role
		}
	}

	assert.Fail(role.t, "Policy "+policy+" was not found")
	return role
}

func (role AssertRole) GetName() string {
	return StringValue(role.role.RoleName)
}

func (role AssertRole) GetAssumedPolicy() string {
	return StringValue(role.RoleDetails().AssumeRolePolicyDocument)
}

func (role AssertRole) RoleDetails() *iam.RoleDetail {
	if role.roleDetails == nil {
		role.roleDetails = aws.GetRoleDetails(role.t, role.GetName(), role.region)
	}

	return role.roleDetails
}
