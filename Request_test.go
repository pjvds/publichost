package publichost_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pjvds/publichost"
	"io/ioutil"
)

var _ = Describe("ReadRequest", func() {
	Context("When reading a valid request", func() {
		buffer := new(bytes.Buffer)

		theId := uint16(1)
		theType := uint8(2)
		theBody := []byte("hello world")
		NewRequest(theId, theType, bytes.NewBuffer(theBody)).Write(buffer)

		request, err := ReadRequest(buffer)

		It("should not error", func() {
			Expect(err).ToNot(HaveOccured())
		})

		It("should have read the id", func() {
			Expect(request.Id).To(Equal(theId))
		})
		It("should have read the type", func() {
			Expect(request.Type).To(Equal(theType))
		})
		It("should have read the length", func() {
			Expect(request.Length).To(Equal(uint16(len(theBody))))
		})
		It("should have read the body", func() {
			actual, _ := ioutil.ReadAll(request.Body)

			Expect(actual).To(Equal(theBody))
		})
	})
})

var _ = Describe("WriteRequest", func() {
	Context("When writing a valid request", func() {
		buffer := new(bytes.Buffer)

		theId := uint16(1)
		theType := uint8(2)
		theBody := []byte("hello world")

		err := NewRequest(theId, theType, bytes.NewBuffer(theBody)).Write(buffer)

		It("writing should not error", func() {
			Expect(err).ToNot(HaveOccured())
		})

		request, err := ReadRequest(buffer)

		It("reading should not error", func() {
			Expect(err).ToNot(HaveOccured())
		})

		It("should have written the id", func() {
			Expect(request.Id).To(Equal(theId))
		})
		It("should have read the type", func() {
			Expect(request.Type).To(Equal(theType))
		})
		It("should have read the length", func() {
			Expect(request.Length).To(Equal(uint16(len(theBody))))
		})
		It("should have read the body", func() {
			actual, _ := ioutil.ReadAll(request.Body)

			Expect(actual).To(Equal(theBody))
		})
	})
})
