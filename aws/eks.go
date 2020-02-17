package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewEksClientE(region string) (*eks.EKS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return eks.New(sess), nil
}

func GetEksClusterByName(t *testing.T, clusterName string, region string) *eks.Cluster {
	clusterOutput, err := GetEksClusterByNameE(clusterName, region)
	require.NoError(t, err)
	return clusterOutput.Cluster
}

func GetEksClusterByNameE(clusterName string, region string) (*eks.DescribeClusterOutput, error) {
	client, err := NewEksClientE(region)
	if err != nil {
		return nil, err
	}

	clusterOutput, err := client.DescribeCluster(&eks.DescribeClusterInput{Name: aws.String(clusterName)})
	if err != nil {
		return nil, err
	}

	return clusterOutput, nil
}

func GetEksGroupNodeByName(t *testing.T, clusterName string, groupNodeName string, region string) *eks.Nodegroup {
	outputNodeGroup, err := GetEksGroupNodeByNameE(clusterName, groupNodeName, region)
	require.NoError(t, err)
	return outputNodeGroup.Nodegroup
}

func GetEksGroupNodeByNameE(clusterName string, groupNodeName string, region string) (*eks.DescribeNodegroupOutput, error) {
	client, err := NewEksClientE(region)
	if err != nil {
		return nil, err
	}

	inputNodeGroup := eks.DescribeNodegroupInput{ClusterName: aws.String(clusterName), NodegroupName: aws.String(groupNodeName)}
	outputNodeGroup, err := client.DescribeNodegroup(&inputNodeGroup)
	if err != nil {
		return nil, err
	}

	return outputNodeGroup, err
}
