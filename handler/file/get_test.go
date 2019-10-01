package file

import (
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
)

var _ = Describe("Get File", func() {
	var res *httptest.ResponseRecorder

	BeforeEach(func() {
		res = httptest.NewRecorder()
		Expect(res).ShouldNot(BeNil())
	})

	// Create default file for get
	Context("Creating default file for subsequent Get", func() {
		It("should succeed with StatusCreated", func() {
			filename := filepath.Join(DataDir, "output.pdf")
			Expect(filename).Should(BeARegularFile())
			Expect(filename).Should(BeAnExistingFile())

			body, ctype, err := createFormFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ctype).ShouldNot(BeZero())
			Expect(body).ShouldNot(BeNil())
			Expect(ctype).ShouldNot(BeZero())

			url := path.Join(Server.URL(), DefaultFileURL)

			req := httptest.NewRequest(http.MethodPost, url, body)

			req.Header.Set("content-type", ctype)

			Expect(req).ShouldNot(BeNil())
			Expect(req.Header.Get("content-type")).Should(BeEquivalentTo(ctype))

			res := httptest.NewRecorder()
			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusCreated))
		})
	})

	Context("Retrieving a file", func() {
		It("should fail with StatusNotFound when file is not available on the server", func() {
			url := path.Join(Server.URL(), "/not/yet/uploaded/file")
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusNotFound))
		})

		It("should fail with StatusBadRequest when the specified directory is not allowed access", func() {
			url := path.Join(Server.URL(), "/not/allowed/?"+urlQueryKeyDirectory+"="+"notallowed")
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusBadRequest))
		})

		It("should succeed with StatusOK when the file resource is available on the server", func() {
			url := path.Join(Server.URL(), DefaultFileURL)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("application/pdf"))
		})
	})
})
