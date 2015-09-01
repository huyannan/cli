package formatters_test

import (
	. "github.com/cloudfoundry/cli/cf/formatters"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("formatting bits to / from strings", func() {
	It("converts kilobits to a human readable description", func() {
		Expect(BitSize(100 * KILOBIT)).To(Equal("100Kb"))
		Expect(BitSize(int64(100.5 * KILOBIT))).To(Equal("100.5Kb"))
	})

	It("parses bit amounts with short units (e.g. K, M, G)", func() {
		var (
			kilobits int64
			err      error
		)

		kilobits, err = ToKilobits("5K")
		Expect(kilobits).To(Equal(int64(5)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("5k")
		Expect(kilobits).To(Equal(int64(5)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("5M")
		Expect(kilobits).To(Equal(int64(5 * 1024)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("5m")
		Expect(kilobits).To(Equal(int64(5 * 1024)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("2G")
		Expect(kilobits).To(Equal(int64(2 * 1024 * 1024)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("3T")
		Expect(kilobits).To(Equal(int64(3 * 1024 * 1024 * 1024)))
		Expect(err).NotTo(HaveOccurred())
	})

	It("parses bit amounts with long units (e.g Kb, Mb, Gb)", func() {
		var (
			kilobits int64
			err      error
		)

		kilobits, err = ToKilobits("5Kb")
		Expect(kilobits).To(Equal(int64(5)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("5kb")
		Expect(kilobits).To(Equal(int64(5)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("5Mb")
		Expect(kilobits).To(Equal(int64(5 * 1024)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("5mb")
		Expect(kilobits).To(Equal(int64(5 * 1024)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("2Gb")
		Expect(kilobits).To(Equal(int64(2 * 1024 * 1024)))
		Expect(err).NotTo(HaveOccurred())

		kilobits, err = ToKilobits("3Tb")
		Expect(kilobits).To(Equal(int64(3 * 1024 * 1024 * 1024)))
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns an error when the unit is missing", func() {
		_, err := ToKilobits("5")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unit of measurement"))
	})

	It("returns an error when the unit is unrecognized", func() {
		_, err := ToKilobits("5Kbb")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unit of measurement"))
	})

	It("allows whitespace before and after the value", func() {
		kilobits, err := ToKilobits("\t\n\r 5kb ")
		Expect(kilobits).To(Equal(int64(5)))
		Expect(err).NotTo(HaveOccurred())
	})
})
