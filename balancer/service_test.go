package balancer

import (
	"time"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("service", func() {
	var subject *service

	BeforeEach(func() {
		var err error
		subject, err = newService("svcname", mockDiscovery{backendA.Address(), backendB.Address()}, time.Minute, time.Minute)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		subject.Close()
	})

	It("should report servers", func() {
		Expect(subject.Servers()).To(ConsistOf([]*balancerpb.Server{
			{Address: backendA.Address()},
			{Address: backendB.Address()},
		}))

		Eventually(func() []*balancerpb.Server {
			return subject.Servers()
		}).Should(ConsistOf([]*balancerpb.Server{
			{Address: backendA.Address(), Score: 10},
			{Address: backendB.Address(), Score: 40},
		}))
	})

	It("should update backends", func() {
		subject.discovery = mockDiscovery{backendA.Address()}
		Expect(subject.updateBackends()).NotTo(HaveOccurred())

		Eventually(func() []*balancerpb.Server {
			return subject.Servers()
		}).Should(ConsistOf([]*balancerpb.Server{
			{Address: backendA.Address(), Score: 10},
		}))
	})

})
