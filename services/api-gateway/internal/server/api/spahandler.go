package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(c *gin.Context) {
	reqPath := c.Request.URL.Path
	ext := filepath.Ext(reqPath)

	if reqPath == "/" {
		if ext != "" {
			c.Status(http.StatusNotFound)
			return
		}
	} else {
		if !(ext == ".css" || ext == ".js" || ext == ".json" || ext == ".ico" ||
			ext == ".png" || ext == ".geojson" || ext == ".svg" ||
			ext == ".otf" || ext == ".eot" || ext == ".woff" || ext == ".txt") {
			c.Status(http.StatusNotFound)
			return
		}
	}

	// 절대 경로 생성 (디렉토리 트래버설 방지)
	path, err := filepath.Abs(reqPath)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	path = filepath.Join(h.staticPath, path)

	// 파일 존재 확인
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// fallback → index.html
		c.File(filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 파일 직접 서빙
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(c.Writer, c.Request)
}
