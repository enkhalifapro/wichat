package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WhisperHelper", func() {
	whisperHelper := &WhisperHelper{}
	Describe("SendAsymMsg", func() {
		It("should return error when pass an empty sender key", func() {
			err := whisperHelper.SendAsymMsg(&Message{})
			Expect(err.Error()).To(Equal("sender key is empty"))
		})
		It("should return error when pass an invalid hex key", func() {
			err := whisperHelper.SendAsymMsg(&Message{From: "nonHex"})
			Expect(err.Error()).To(Equal("invalid sender key"))
		})
		It("should return error when pass an empty recipient key", func() {
			err := whisperHelper.SendAsymMsg(&Message{From: "46342e3c48021fc8e6d5a704ba84c56945037ef151ba38b21663767c0abb3c63"})
			Expect(err.Error()).To(Equal("recipient key is empty"))
		})
		It("should return error when pass an invalid recipient key", func() {
			err := whisperHelper.SendAsymMsg(&Message{From: "46342e3c48021fc8e6d5a704ba84c56945037ef151ba38b21663767c0abb3c63", To: "nonHex"})
			Expect(err.Error()).To(Equal("invalid recipient key"))
		})
	})

})
