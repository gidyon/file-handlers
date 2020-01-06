package static

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

var _ = Describe("Testing server with allowed all directory #allowed", func() {
	var (
		handler  http.Handler
		err      error
		notFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(RootDir, "/index.html"))
		})
		pushFiles = []string{
			"/css/app.f1c93db5.css",
			"/css/chunk-vendors.4ca44aa7.css",
		}
		// push options
		pushContent = map[string][]string{
			"/":           pushFiles,
			"/index.html": pushFiles,
		}
	)

	Context("Creating server", func() {
		It("should create the server", func() {
			handler, err = NewHandler(&ServerOptions{
				RootDir:         RootDir,
				Index:           "index.html",
				URLPathPrefix:   URLPrefix,
				AllowedDirs:     []string{"css", "img"},
				NotFoundHandler: notFound,
				PushContent:     pushContent,
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Sending Request", func() {
		It("should return file resource when requested file is present in the server", func() {
			url := Server.URL() + "/js/about.b5d251bd.js"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			Expect(req).ShouldNot(BeNil())

			res := httptest.NewRecorder()
			Expect(res).ShouldNot(BeNil())

			handler.ServeHTTP(res, req)

			Expect(res.Code).Should(BeEquivalentTo(http.StatusInternalServerError))
		})
	})
})
