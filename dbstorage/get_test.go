package dbstorage

import (
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
)

var _ = Describe("Get File resource #get", func() {
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

		It("should succeed with StatusOK when the file resource is available on the server", func() {
			url := path.Join(Server.URL(), DefaultFileURL)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("application/pdf"))
		})
	})

	// Create owned file + caching for get
	Context("Creating owned file and caching it for subsequent Get", func() {
		It("should succeed with StatusCreated", func() {
			RedisClient.FlushAll()

			filename := filepath.Join(DataDir, CachedFile)
			Expect(filename).Should(BeARegularFile())
			Expect(filename).Should(BeAnExistingFile())

			body, ctype, err := createFormFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ctype).ShouldNot(BeZero())
			Expect(body).ShouldNot(BeNil())
			Expect(ctype).ShouldNot(BeZero())

			url := path.Join(Server.URL(), CachedFileURL)
			url += "?" + urlQueryCacheKey + "=true&" + urlQueryKeyOwnerID + "=" + OwnerID

			req := httptest.NewRequest(http.MethodPost, url, body)

			req.Header.Set("content-type", ctype)

			Expect(req).ShouldNot(BeNil())
			Expect(req.Header.Get("content-type")).Should(BeEquivalentTo(ctype))

			res := httptest.NewRecorder()
			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusCreated))
		})

		Describe("Getting owner file", func() {
			It("should succeed with StatusOK", func() {
				url := path.Join(Server.URL(), CachedFileURL)
				req := httptest.NewRequest(http.MethodGet, url, nil)
				Expect(req).ShouldNot(BeNil())

				Handler.ServeHTTP(res, req)

				Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			})
		})

		Describe("Getting owner file from cache", func() {
			It("should succeed with StatusOK", func() {
				url := path.Join(Server.URL(), CachedFileURL)
				url += "?" + urlQueryCacheKey + "=true"

				req := httptest.NewRequest(http.MethodGet, url, nil)
				Expect(req).ShouldNot(BeNil())

				Handler.ServeHTTP(res, req)

				Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			})
		})
	})
})
