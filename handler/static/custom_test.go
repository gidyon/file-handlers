package static

import (
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

var _ = Describe("Custom notFoundhandler", func() {

	var (
		handler http.Handler
		server  *ghttp.Server
		res     *httptest.ResponseRecorder
	)

	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(RootDir, "/index.html"))
	})

	// Setup handler
	handler, err := NewHandler(RootDir, URLPrefix, notFound, nil)
	Context("Setup Expectations", func() {
		It("should set up handler without error", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(handler).ShouldNot(BeNil())
		})
	})

	// setup server
	server = ghttp.NewServer()
	server.AppendHandlers(http.HandlerFunc(handler.ServeHTTP))

	BeforeEach(func() {
		res = httptest.NewRecorder()
		Expect(res).ShouldNot(BeNil())
	})

	Context("Sending Request", func() {

		It("should return index page when requested file resource is not in the server", func() {
			url := server.URL() + "/notfound"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("text/html"))
		})

		It("should return file resource when requested file is present in the server", func() {
			url := server.URL() + "/css/app.f1c93db5.css"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("text/css"))
		})
	})
})
