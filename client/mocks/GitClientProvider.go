// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	client "guard_rails/client"

	mock "github.com/stretchr/testify/mock"
)

// GitClientProvider is an autogenerated mock type for the GitClientProvider type
type GitClientProvider struct {
	mock.Mock
}

// NewGitClient provides a mock function with given fields:
func (_m *GitClientProvider) NewGitClient() client.GitClient {
	ret := _m.Called()

	var r0 client.GitClient
	if rf, ok := ret.Get(0).(func() client.GitClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.GitClient)
		}
	}

	return r0
}

type mockConstructorTestingTNewGitClientProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewGitClientProvider creates a new instance of GitClientProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGitClientProvider(t mockConstructorTestingTNewGitClientProvider) *GitClientProvider {
	mock := &GitClientProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
