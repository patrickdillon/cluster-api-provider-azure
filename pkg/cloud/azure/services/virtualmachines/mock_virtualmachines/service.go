// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_virtualmachines is a generated GoMock package.
package mock_virtualmachines

import (
	context "context"
	reflect "reflect"

	compute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	gomock "github.com/golang/mock/gomock"
)

// MockVMImagesClient is a mock of VMImagesClient interface.
type MockVMImagesClient struct {
	ctrl     *gomock.Controller
	recorder *MockVMImagesClientMockRecorder
}

// MockVMImagesClientMockRecorder is the mock recorder for MockVMImagesClient.
type MockVMImagesClientMockRecorder struct {
	mock *MockVMImagesClient
}

// NewMockVMImagesClient creates a new mock instance.
func NewMockVMImagesClient(ctrl *gomock.Controller) *MockVMImagesClient {
	mock := &MockVMImagesClient{ctrl: ctrl}
	mock.recorder = &MockVMImagesClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVMImagesClient) EXPECT() *MockVMImagesClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockVMImagesClient) Get(ctx context.Context, location, publisherName, offer, skus, version string) (compute.VirtualMachineImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, location, publisherName, offer, skus, version)
	ret0, _ := ret[0].(compute.VirtualMachineImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockVMImagesClientMockRecorder) Get(ctx, location, publisherName, offer, skus, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVMImagesClient)(nil).Get), ctx, location, publisherName, offer, skus, version)
}
