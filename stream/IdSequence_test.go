package stream_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pjvds/publichost/stream"
)

var _ = Describe("IdSequence", func() {
	var lastId Id
	var sequence IdSequence

	BeforeEach(func() {
		sequence = NewIdSequence()
		lastId = sequence.Next()

		Expect(lastId).To(Equal(Id(1)))
	})

	Context("when advancing to next", func() {
		var nextId Id

		BeforeEach(func() {
			nextId = sequence.Next()
		})

		It("it should increment the id with one", func() {
			Expect(nextId).To(Equal(Id(lastId + 1)))
		})
	})
})
