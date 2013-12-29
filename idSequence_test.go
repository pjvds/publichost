package publichost

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("idSequence", func() {
	var lastId uint32
	var sequence *idSequence

	BeforeEach(func() {
		sequence = newIdSequence()
		lastId = sequence.Next()

		Expect(lastId).To(Equal(uint32(1)))
	})

	Context("when advancing to next", func() {
		var nextId uint32

		BeforeEach(func() {
			nextId = sequence.Next()
		})

		It("it should increment the id with one", func() {
			Expect(nextId).To(Equal(uint32(lastId + 1)))
		})
	})
})
