/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package securitygroups

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/Azure/go-autorest/autorest"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/cloud/azure"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/cloud/azure/actuators"
)

// Service provides operations on resource groups
type Service struct {
	Client network.SecurityGroupsClient
	Scope  *actuators.MachineScope
}

// getGroupsClient creates a new groups client from subscriptionid.
func getSecurityGroupsClient(resourceManagerEndpoint, subscriptionID string, authorizer autorest.Authorizer) network.SecurityGroupsClient {
	securityGroupsClient := network.NewSecurityGroupsClientWithBaseURI(resourceManagerEndpoint, subscriptionID)
	securityGroupsClient.Authorizer = authorizer
	securityGroupsClient.AddToUserAgent(azure.UserAgent)
	return securityGroupsClient
}

// NewService creates a new groups service.
func NewService(scope *actuators.MachineScope) azure.Service {
	if scope.IsStackHub() {
		return NewStackHubService(scope)
	}

	return &Service{
		Client: getSecurityGroupsClient(scope.ResourceManagerEndpoint, scope.SubscriptionID, scope.Authorizer),
		Scope:  scope,
	}
}
