// Code generated by mockery v1.0.0
package poolmocks

import mock "github.com/stretchr/testify/mock"
import pool "github.com/cameronwp/glacier/pool"

// Drainer is an autogenerated mock type for the Drainer type
type Drainer struct {
	mock.Mock
}

// Drain provides a mock function with given fields: _a0
func (_m *Drainer) Drain(_a0 *pool.JobQueue) {
	_m.Called(_a0)
}
