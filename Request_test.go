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
		theLength := uint16(0)
		NewRequest(theId, theType, nil).Write(buffer)

		request, err := ReadRequest(buffer)

		It("should not error", func() {
			Expect(err).ToNot(HaveOccured())
		})

		It("should have read the data", func() {
			Expect(request.Id).To(Equal(theId))
			Expect(request.Type).To(Equal(theType))
			Expect(request.Length).To(Equal(theLength))
		})
	})

	Context("When reading a valid request", func() {
		buffer := new(bytes.Buffer)

		theId := uint16(1)
		theType := uint8(2)
		theBody := bytes.NewBufferString("hello world")
		NewRequest(theId, theType, theBody).Write(buffer)

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
			Expect(request.Length).To(Equal(uint16(theBody.Len())))
		})
		It("should have read the body", func() {
			actual, _ := ioutil.ReadAll(request.Body)

			Expect(actual).To(Equal(theBody.Bytes()))
		})
	})
})
