package stream_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pjvds/publichost/stream"
)

var _ = Describe("ThreadSafeMap", func() {
	var theMap Map

	BeforeEach(func() {
		theMap = NewThreadSafeMap()
	})

	Context("Adding a stream", func() {
		var id Id
		var stream Stream

		BeforeEach(func() {
			id = Id(1)
			stream = NewNullStream()

			theMap.Add(id, stream)
		})

		It("should be available", func() {
			getStream, getErr := theMap.Get(id)
			Expect(stream).To(Equal(getStream))
			Expect(getErr).ToNot(HaveOccured())
		})

		Context("Deleting the stream", func() {
			var delErr error

			BeforeEach(func() {
				delErr = theMap.Delete(id)
			})

			It("should not error", func() {
				Expect(delErr).ToNot(HaveOccured())
			})

			It("should not be available anymore", func() {
				_, getErr := theMap.Get(id)
				Expect(getErr).To(HaveOccured())
			})
		})

		Context("Adding another with the same id", func() {
			var otherStream Stream
			var addErr error

			BeforeEach(func() {
				otherStream = NewNullStream()
				addErr = theMap.Add(id, otherStream)
			})

			It("should error", func() {
				Expect(addErr).To(HaveOccured())
			})
		})
	})

	Context("Getting a non existing stream", func() {
		var nonId Id
		var getVal Stream
		var getErr error

		BeforeEach(func() {
			nonId = Id(999)
			getVal, getErr = theMap.Get(nonId)
		})

		It("should return no stream", func() {
			Expect(getVal).To(BeNil())
		})
		It("should error", func() {
			Expect(getErr).To(HaveOccured())
		})
	})
})
