package file

import (
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
)

var _ = Describe("Save File", func() {
	var res *httptest.ResponseRecorder

	BeforeEach(func() {
		res = httptest.NewRecorder()
		Expect(res).ShouldNot(BeNil())
	})

	Context("Saving file", func() {

		It("should fail with StatusBadRequest when content-type in request is incorrect", func() {
			filename := filepath.Join(DataDir, "leo.jpg")
			Expect(filename).Should(BeARegularFile())
			Expect(filename).Should(BeAnExistingFile())

			body, ctype, err := createFormFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ctype).ShouldNot(BeZero())
			Expect(body).ShouldNot(BeNil())
			Expect(ctype).ShouldNot(BeZero())

			url := path.Join(Server.URL(), "/image1")

			req := httptest.NewRequest(http.MethodPost, url, body)

			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusBadRequest))
		})

		It("should fail with StatusBadRequest when form file size is larger than cap", func() {
			filename := filepath.Join(DataDir, "big file.mp4")
			Expect(filename).Should(BeARegularFile())
			Expect(filename).Should(BeAnExistingFile())

			body, ctype, err := createFormFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ctype).ShouldNot(BeZero())
			Expect(body).ShouldNot(BeNil())
			Expect(ctype).ShouldNot(BeZero())

			url := path.Join(Server.URL(), "/image1")

			req := httptest.NewRequest(http.MethodPost, url, body)

			req.Header.Set("content-type", ctype)

			Expect(req).ShouldNot(BeNil())
			Expect(req.Header.Get("content-type")).Should(BeEquivalentTo(ctype))

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusBadRequest))
		})

		It("should fail with StatusBadRequest when the specified directory is not allowed access", func() {
			filename := filepath.Join(DataDir, "sala.webp")
			Expect(filename).Should(BeARegularFile())
			Expect(filename).Should(BeAnExistingFile())

			body, ctype, err := createFormFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ctype).ShouldNot(BeZero())
			Expect(body).ShouldNot(BeNil())
			Expect(ctype).ShouldNot(BeZero())

			url := path.Join(Server.URL(), "/image2?"+urlQueryKeyDirectory+"="+"classified")

			req := httptest.NewRequest(http.MethodPost, url, body)

			req.Header.Set("content-type", ctype)

			Expect(req).ShouldNot(BeNil())
			Expect(req.Header.Get("content-type")).Should(BeEquivalentTo(ctype))

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusBadRequest))
		})

		It("should succeed with StatusCreated when the image file is sent as part of request", func() {
			filename := filepath.Join(DataDir, "leo.jpg")
			Expect(filename).Should(BeARegularFile())
			Expect(filename).Should(BeAnExistingFile())

			body, ctype, err := createFormFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ctype).ShouldNot(BeZero())
			Expect(body).ShouldNot(BeNil())
			Expect(ctype).ShouldNot(BeZero())

			url := path.Join(Server.URL(), "/image1")

			req := httptest.NewRequest(http.MethodPost, url, body)

			req.Header.Set("content-type", ctype)

			Expect(req).ShouldNot(BeNil())
			Expect(req.Header.Get("content-type")).Should(BeEquivalentTo(ctype))

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusCreated))
		})
	})
})
