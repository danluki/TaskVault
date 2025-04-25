package taskvault

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const uiPathPrefix = "ui/"

//go:embed ui-dist
var uiDist embed.FS

func (h *HTTPTransport) UI(r *gin.RouterGroup) {
	r.GET(
		"/", func(c *gin.Context) {
			switch c.NegotiateFormat(gin.MIMEHTML) {
			case gin.MIMEHTML:
				c.Redirect(http.StatusSeeOther, "/ui/")
			default:
				c.AbortWithStatus(http.StatusNotFound)
			}
		},
	)

	r.GET(
		"/login", func(c *gin.Context) {
			switch c.NegotiateFormat(gin.MIMEHTML) {
			case gin.MIMEHTML:
				c.Redirect(http.StatusSeeOther, "/ui/login")
			default:
				c.AbortWithStatus(http.StatusNotFound)
			}
		},
	)

	r.GET(
		"/storage", func(c *gin.Context) {
			switch c.NegotiateFormat(gin.MIMEHTML) {
			case gin.MIMEHTML:
				c.Redirect(http.StatusSeeOther, "/ui/storage")
			default:
				c.AbortWithStatus(http.StatusNotFound)
			}
		},
	)

	r.GET(
		"/dashboard", func(c *gin.Context) {
			switch c.NegotiateFormat(gin.MIMEHTML) {
			case gin.MIMEHTML:
				c.Redirect(http.StatusSeeOther, "/ui/dashboard")
			default:
				c.AbortWithStatus(http.StatusNotFound)
			}
		},
	)

	ui := r.Group("/" + uiPathPrefix)

	assets, err := fs.Sub(uiDist, "ui-dist")
	if err != nil {
		h.logger.Fatal(err)
	}
	a, err := assets.Open("index.html")
	if err != nil {
		h.logger.Fatal(err)
	}
	b, err := io.ReadAll(a)
	if err != nil {
		h.logger.Fatal(err)
	}
	t, err := template.New("index.html").Parse(string(b))
	if err != nil {
		h.logger.Fatal(err)
	}
	h.Engine.SetHTMLTemplate(t)

	ui.GET(
		"/*filepath", func(ctx *gin.Context) {
			p := ctx.Param("filepath")
			f := strings.TrimPrefix(p, "/")
			_, err := assets.Open(f)
			if err == nil && p != "/" && p != "/index.html" {
				ctx.FileFromFS(p, http.FS(assets))
			} else {
				pairs, err := h.agent.Store.GetAllValues()
				if err != nil {
					h.logger.Error(err)
				}
				var (
					totalPairs                             = len(pairs)
					pairsAdded, pairsUpdated, pairsDeleted = 0, 0, 0
				)
				l, err := h.agent.leaderMember()
				ln := "no leader"
				if err != nil {
					h.logger.Error(err)
				} else {
					ln = l.Name
				}
				ctx.HTML(
					http.StatusOK, "index.html", gin.H{
						"TASKVAULT_API_URL": fmt.Sprintf(
							"../%s", apiPathPrefix,
						),
						"TASKVAULT_LEADER":        ln,
						"TASKVAULT_TOTAL_PAIRS":   totalPairs,
						"TASKVAULT_PAIRS_ADDED":   pairsAdded,
						"TASKVAULT_PAIRS_UPDATED": pairsUpdated,
						"TASKVAULT_PAIRS_DELETED": pairsDeleted,
					},
				)
			}
		},
	)
}
