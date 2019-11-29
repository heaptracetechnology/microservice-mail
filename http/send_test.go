package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/oms-services/email"
	. "github.com/oms-services/email/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	password string
	to       string
	from     string
)

var _ = BeforeSuite(func() {
	password = getEnvOrError("EMAIL_PASSWORD")
	to = getEnvOrError("EMAIL_TO")
	from = getEnvOrError("EMAIL_FROM")
})

var _ = Describe("Sending Emails", func() {

	var (
		recorder    *httptest.ResponseRecorder
		requestBody *bytes.Buffer
	)

	BeforeEach(func() {
		recorder = nil
		requestBody = &bytes.Buffer{}

		os.Unsetenv("PASSWORD")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("IMAP_HOST")
		os.Unsetenv("IMAP_PORT")
	})

	JustBeforeEach(func() {
		request, err := http.NewRequest("POST", "/send", requestBody)
		Expect(err).NotTo(HaveOccurred())

		recorder = httptest.NewRecorder()
		handler := SendHandler{}
		handler.ServeHTTP(recorder, request)
	})

	When("all env vars are set correctly", func() {
		BeforeEach(func() {
			os.Setenv("PASSWORD", password)
			os.Setenv("SMTP_HOST", "smtp.gmail.com")
			os.Setenv("SMTP_PORT", "587")
		})

		When("a valid body is sent in the request", func() {
			BeforeEach(func() {
				email := email.Email{
					From:    from,
					To:      []string{to},
					Subject: "Testing microservice",
					Body:    "Any body message to test"}

				Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())
			})

			It("should result in a successful SMTP response", func() {
				Expect(recorder.Code).To(Equal(250))
			})
		})

		When("an invalid body is sent in the request", func() {
			BeforeEach(func() {
				email := []byte(`{"invalid":body}`)
				Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())
			})

			It("should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("the body does not contain required details", func() {
			BeforeEach(func() {
				email := email.Email{}
				Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())
			})

			It("should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	When("not all env vars are set", func() {
		When("no env vars are set", func() {
			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("no smtp host is set", func() {
			BeforeEach(func() {
				os.Setenv("PASSWORD", password)
				os.Setenv("SMTP_PORT", "587")
			})

			It("Should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})