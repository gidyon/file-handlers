package static

import (
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Static File Server with http.NotFound as notFound handler", func() {

	var res *httptest.ResponseRecorder

	BeforeEach(func() {
		res = httptest.NewRecorder()
		Expect(res).ShouldNot(BeNil())
	})

	Context("Sending Request", func() {
		It("should return StatusBadRequest when method is not GET", func() {
			req := httptest.NewRequest(http.MethodPost, Server.URL()+"/", nil)
			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusBadRequest))
		})

		It("should return file resource when requested file is present in the server", func() {
			url := Server.URL() + "/js/about.b5d251bd.js"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("application/javascript"))
		})

		It("should return index page when requested path is /index.html", func() {
			url := Server.URL() + "/index.html"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("text/html"))
		})

		It("should return index page when requested path is /", func() {
			url := Server.URL() + "/"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("text/html"))
		})
	})

	Context("Receiving Response", func() {
		It("should return StatusBadRequest when method is not GET", func() {
			req := httptest.NewRequest(http.MethodPost, Server.URL()+"/", nil)
			Expect(req).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusBadRequest))
		})

		It("should return StatusNotFound when the requested file resource is not in the server", func() {
			url := Server.URL() + "/notfound"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusNotFound))
		})

		It("should return StatusOK when requested file is present in the server", func() {
			url := Server.URL() + "/js/about.b5d251bd.js"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			Handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("application/javascript"))
		})
	})
})
