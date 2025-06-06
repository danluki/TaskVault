package taskvault

import (
	"net/http"
	"strconv"

	"github.com/danluki/taskvault/pkg/types"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	pretty        = "pretty"
	apiPathPrefix = "v1"
)

type Transport interface {
	ServeHTTP()
}

type HTTPTransport struct {
	Engine *gin.Engine

	agent  *Agent
	logger *zap.SugaredLogger
}

func NewTransport(a *Agent, log *zap.SugaredLogger) *HTTPTransport {
	return &HTTPTransport{
		agent:  a,
		logger: log,
	}
}

func (h *HTTPTransport) ServeHTTP() {
	h.Engine = gin.Default()

	rootPath := h.Engine.Group("/")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"*"}
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = []string{"*"}

	rootPath.Use(cors.New(config))

	h.APIRoutes(rootPath)
	if h.agent.config.UI {
		h.UI(rootPath)
	}

	h.logger.Info("api: Running HTTP server", zap.String("address", h.agent.config.HTTPAddr))

	go func() {
		if err := h.Engine.Run(h.agent.config.HTTPAddr); err != nil {
			panic(err)
		}
	}()
}

func (h *HTTPTransport) APIRoutes(
	r *gin.RouterGroup, middleware ...gin.HandlerFunc,
) {
	h.Engine.GET(
		"/health", func(c *gin.Context) {
			c.JSON(
				http.StatusOK, gin.H{
					"status": "healthy",
				},
			)
		},
	)

	if h.agent.config.EnablePrometheus {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	r.GET("/v1", h.indexHandler)
	v1 := r.Group("/v1")
	v1.Use(middleware...)
	v1.GET("/", h.indexHandler)
	v1.GET("/members", h.membersHandler)
	v1.GET("/leader", h.leaderHandler)
	v1.GET("/isleader", h.isLeaderHandler)
	v1.POST("/leave", h.leaveHandler)

	pairs := v1.Group("/storage")
	pairs.GET("", h.pairsHandler)
	pairs.GET("/:key", h.pairGetHandler)
	pairs.POST("", h.pairPostHandler)
	pairs.DELETE("/:key", h.pairDeleteHandler)
	pairs.PATCH("/", h.pairDeleteHandler)
}

func renderJSON(c *gin.Context, status int, v interface{}) {
	if _, ok := c.GetQuery(pretty); ok {
		c.IndentedJSON(status, v)
	} else {
		c.JSON(status, v)
	}
}

func (h *HTTPTransport) membersHandler(c *gin.Context) {
	mems := []*types.Member{}
	for _, m := range h.agent.serf.Members() {
		id, _ := uuid.GenerateUUID()
		mid := &types.Member{m, id, m.Status.String()}
		mems = append(mems, mid)
	}
	c.Header("X-Total-Count", strconv.Itoa(len(mems)))
	renderJSON(c, http.StatusOK, mems)
}

func (h *HTTPTransport) leaderHandler(c *gin.Context) {
	member, err := h.agent.leaderMember()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	if member == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	renderJSON(c, http.StatusOK, member)
}

func (h *HTTPTransport) isLeaderHandler(c *gin.Context) {
	isleader := h.agent.IsLeader()
	if isleader {
		renderJSON(c, http.StatusOK, "leader")
		return
	}

	renderJSON(c, http.StatusNotFound, "follower")
}

func (h *HTTPTransport) leaveHandler(c *gin.Context) {
	if err := h.agent.Stop(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	renderJSON(c, http.StatusOK, h.agent.serf.Memberlist())
}

func (h *HTTPTransport) indexHandler(c *gin.Context) {
	local := h.agent.serf.LocalMember()

	stats := map[string]map[string]string{
		"agent": {
			"name":    local.Name,
			"version": Version,
		},
		"serf": h.agent.serf.Stats(),
		"tags": local.Tags,
	}

	renderJSON(c, http.StatusOK, stats)
}

func (h *HTTPTransport) pairsHandler(c *gin.Context) {
	pairs, err := h.agent.Store.GetAllValues()
	if err != nil {
		return
	}

	start, ok := c.GetQuery("_start")
	if !ok {
		start = "0"
	}
	s, _ := strconv.Atoi(start)

	end, ok := c.GetQuery("_end")
	e := 0
	if !ok {
		e = len(pairs)
	} else {
		e, _ = strconv.Atoi(end)
		if e > len(pairs) {
			e = len(pairs)
		}
	}

	c.Header("X-Total-Count", strconv.Itoa(len(pairs)))
	renderJSON(c, http.StatusOK, pairs[s:e])
}

func (h *HTTPTransport) pairGetHandler(c *gin.Context) {
	pairName := c.Param("key")

	pair, err := h.agent.Store.GetValue(pairName)
	if err != nil {
		h.logger.Error(err)
		c.Status(http.StatusNotFound)
		return
	}

	renderJSON(c, http.StatusOK, pair)
}

func (h *HTTPTransport) pairDeleteHandler(c *gin.Context) {
	keyName := c.Param("key")

	err := h.agent.GRPCClient.DeleteValue(keyName)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *HTTPTransport) pairPostHandler(c *gin.Context) {
	pair := &Pair{}
	if err := c.ShouldBindJSON(pair); err != nil {
		h.logger.Error(err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, err := h.agent.GRPCClient.CreateValue(
		pair.Key, pair.Value,
	); err != nil {
		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Status(http.StatusCreated)
}
