package publichost_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pjvds/publichost"
)

var _ = Describe("ReadRequest", func() {
	Context("When reading a valid request", func() {
		buffer := new(bytes.Buffer)

		TheId := uint16(1)
		TheType := uint8(2)
		TheLength := uint16(0)

		Request{
			Id:     TheId,
			Type:   TheType,
			Length: TheLength,
		}.Write(buffer)

		request, err := ReadRequest(buffer)

		It("should not error", func() {
			Expect(err).ToNot(HaveOccured())
		})

		It("should have read the id", func() {
			Expect(request.Id).To(Equal(TheId))
		})
		It("should have read the type", func() {
			Expect(request.Type).To(Equal(TheType))
		})
		It("should have read the length", func() {
			Expect(request.Length).To(Equal(TheLength))
		})
	})
})
