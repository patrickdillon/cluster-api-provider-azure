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

package virtualmachineextensions

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/cloud/azure"
)

// Spec input specification for Get/CreateOrUpdate/Delete calls
type Spec struct {
	Name       string
	VMName     string
	ScriptData string
}

// Get provides information about a virtual network.
func (s *Service) Get(ctx context.Context, spec azure.Spec) (interface{}, error) {
	vmExtSpec, ok := spec.(*Spec)
	if !ok {
		return compute.VirtualMachineExtension{}, errors.New("invalid vm specification")
	}
	vmExt, err := s.Client.Get(ctx, s.Scope.MachineConfig.ResourceGroup, vmExtSpec.VMName, vmExtSpec.Name, "")
	if err != nil && azure.ResourceNotFound(err) {
		return nil, fmt.Errorf("vm extension %s not found: %w", vmExtSpec.Name, err)
	} else if err != nil {
		return vmExt, err
	}
	return vmExt, nil
}

// CreateOrUpdate creates or updates a virtual network.
func (s *Service) CreateOrUpdate(ctx context.Context, spec azure.Spec) error {
	vmExtSpec, ok := spec.(*Spec)
	if !ok {
		return errors.New("invalid vm specification")
	}

	klog.V(2).Infof("creating vm extension %s ", vmExtSpec.Name)

	future, err := s.Client.CreateOrUpdate(
		ctx,
		s.Scope.MachineConfig.ResourceGroup,
		vmExtSpec.VMName,
		vmExtSpec.Name,
		compute.VirtualMachineExtension{
			Name:     to.StringPtr(vmExtSpec.Name),
			Location: to.StringPtr(s.Scope.MachineConfig.Location),
			VirtualMachineExtensionProperties: &compute.VirtualMachineExtensionProperties{
				Type:                    to.StringPtr("CustomScript"),
				TypeHandlerVersion:      to.StringPtr("2.0"),
				AutoUpgradeMinorVersion: to.BoolPtr(true),
				Settings:                map[string]bool{"skipDos2Unix": true},
				Publisher:               to.StringPtr("Microsoft.Azure.Extensions"),
				ProtectedSettings:       map[string]string{"script": vmExtSpec.ScriptData},
			},
		})
	if err != nil {
		return fmt.Errorf("cannot create vm extension: %w", err)
	}

	err = future.WaitForCompletionRef(ctx, s.Client.Client)
	if err != nil {
		return fmt.Errorf("cannot get the extension create or update future response: %w", err)
	}

	_, err = future.Result(s.Client)
	if err != nil {
		return fmt.Errorf("cannot create vm: %w", err)
	}

	// if *vmExt.ProvisioningState != string(compute.ProvisioningStateSucceeded) {
	// 	// If the script failed delete it so it can be retried
	// 	s.Delete(ctx, vmExtSpec)
	// }

	klog.V(2).Infof("successfully created vm extension %s ", vmExtSpec.Name)
	return err
}

// Delete deletes the virtual network with the provided name.
func (s *Service) Delete(ctx context.Context, spec azure.Spec) error {
	vmExtSpec, ok := spec.(*Spec)
	if !ok {
		return errors.New("Invalid VNET Specification")
	}
	klog.V(2).Infof("deleting vm extension %s ", vmExtSpec.Name)
	future, err := s.Client.Delete(ctx, s.Scope.MachineConfig.ResourceGroup, vmExtSpec.VMName, vmExtSpec.Name)
	if err != nil && azure.ResourceNotFound(err) {
		// already deleted
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to delete vm extension %s in resource group %s: %w", vmExtSpec.Name, s.Scope.MachineConfig.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, s.Client.Client)
	if err != nil {
		return fmt.Errorf("cannot delete, future response: %w", err)
	}

	_, err = future.Result(s.Client)

	klog.V(2).Infof("successfully deleted vm %s ", vmExtSpec.Name)
	return err
}
