package static

import (
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

var _ = Describe("Server push with custom notFoundhandler", func() {

	var (
		handler http.Handler
		server  *ghttp.Server
		res     *httptest.ResponseRecorder
	)

	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(RootDir, "/index.html"))
	})

	pushFiles := []string{
		"/css/app.f1c93db5.css",
		"/css/chunk-vendors.4ca44aa7.css",
		"/js/app.2a385984.js",
		"/js/chunk-vendors.fe8e2aad.js",
	}
	// push options
	pushContent := map[string][]string{
		"/":           pushFiles,
		"/index.html": pushFiles,
	}

	// Setup handler
	handler, err := NewHandler(&ServerOptions{
		RootDir:         RootDir,
		Index:           "index.html",
		URLPathPrefix:   URLPrefix,
		AllowedDirs:     nil,
		NotFoundHandler: notFound,
		PushContent:     pushContent,
	})
	Context("setup handler", func() {
		It("should setup handler without an error", func() {
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
		It("should push resources to the user", func() {
			url := server.URL() + "/"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusOK))
			Expect(res.Header().Get("content-type")).Should(ContainSubstring("text/html"))
		})
	})
})
