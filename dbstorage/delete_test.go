package dbstorage

import (
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
)

var _ = Describe("Deleting s File resource #delete", func() {
	var res *httptest.ResponseRecorder

	BeforeEach(func() {
		res = httptest.NewRecorder()
		Expect(res).ShouldNot(BeNil())
	})

	Context("Deleting file", func() {

		It("should succeed with StatusOK when the file resource is not present on the server", func() {
			url := path.Join(Server.URL(), "/image1/not/exist/oops")

			req := httptest.NewRequest(http.MethodDelete, url, nil)

			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
		})

		const DeleteFileURL = "/myfile/2"

		// Create default file for delete
		Context("Delete a created file", func() {

			It("should succeed with StatusCreated", func() {
				filename := filepath.Join(DataDir, "output.pdf")
				Expect(filename).Should(BeARegularFile())
				Expect(filename).Should(BeAnExistingFile())

				body, ctype, err := createFormFile(filename)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ctype).ShouldNot(BeZero())
				Expect(body).ShouldNot(BeNil())
				Expect(ctype).ShouldNot(BeZero())

				url := path.Join(Server.URL(), DeleteFileURL)

				req := httptest.NewRequest(http.MethodPost, url, body)

				req.Header.Set("content-type", ctype)

				Expect(req).ShouldNot(BeNil())
				Expect(req.Header.Get("content-type")).Should(BeEquivalentTo(ctype))

				res := httptest.NewRecorder()
				Handler.ServeHTTP(res, req)

				Expect(res.Code).Should(BeEquivalentTo(http.StatusCreated))
			})

			It("should succeed with StatusOK when the file resource is available on the server", func() {
				url := path.Join(Server.URL(), DeleteFileURL)

				req := httptest.NewRequest(http.MethodDelete, url, nil)

				Expect(req).ShouldNot(BeNil())

				Handler.ServeHTTP(res, req)

				Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			})
		})

	})
})
