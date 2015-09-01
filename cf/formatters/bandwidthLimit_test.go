package formatters_test

import (
	. "github.com/cloudfoundry/cli/cf/formatters"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("bandwidthLimit formatting", func() {
	It("returns 'Unlimited' when limit is -1", func() {
		Expect(InstanceBandwidthLimit(-1)).To(Equal("Unlimited"))
	})

	It("formats original value to <val>Kb when limit is not -1", func() {
		Expect(InstanceBandwidthLimit(100)).To(Equal("100Kb"))
	})
})
