package aws

import (
	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
	"testing"
)

var accountAuthorizationDetailsFilter []string = []string{"Role", "LocalManagedPolicy", "AWSManagedPolicy"}

func GetRole(t *testing.T, roleName string, region string) *iam.Role {
	role, err := GetRoleE(roleName, region)
	require.NoError(t, err)
	return role
}

func GetRoleE(roleName string, region string) (*iam.Role, error) {
	client, err := NewIamClientE(region)
	if err != nil {
		return nil, err
	}

	output, err := client.GetRole(&iam.GetRoleInput{RoleName: awsSdk.String(roleName)})
	if err != nil {
		return nil, err
	}

	return output.Role, nil
}

func GetRoleDetails(t *testing.T, roleName string, region string) *iam.RoleDetail {
	roleDetails, err := GetRoleDetailsE(roleName, region)
	require.NoError(t, err)
	return roleDetails
}

func GetRoleDetailsE(roleName string, region string) (*iam.RoleDetail, error) {
	client, err := NewIamClientE(region)
	if err != nil {
		return nil, err
	}

	filter := awsSdk.StringSlice(accountAuthorizationDetailsFilter)
	accountAuthorizationDetails, err := client.GetAccountAuthorizationDetails(&iam.GetAccountAuthorizationDetailsInput{Filter: filter})
	if err != nil {
		return nil, err
	}

	for _, roleDetails := range accountAuthorizationDetails.RoleDetailList {
		if awsSdk.StringValue(roleDetails.RoleName) == roleName {
			return roleDetails, nil
		}
	}

	return nil, aws.NewNotFoundError("RoleDetail", "Account-Id-Will-Not-Be-Shown", region)
}

func NewIamClientE(region string) (*iam.IAM, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}
