package taskvault

import (
	"net/http"
	"strconv"

	"github.com/danluki/taskvault/pkg/types"
	"github.com/gin-contrib/expvar"
	"github.com/hashicorp/go-uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	pretty        = "pretty"
	apiPathPrefix = "v1"
)

// Transport is the interface that wraps the ServeHTTP method.
type Transport interface {
	ServeHTTP()
}

// HTTPTransport stores pointers to an agent and a gin Engine.
type HTTPTransport struct {
	Engine *gin.Engine

	agent  *Agent
	logger *logrus.Entry
}

// NewTransport creates an HTTPTransport with a bound agent.
func NewTransport(a *Agent, log *logrus.Entry) *HTTPTransport {
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
	rootPath.Use(h.MetaMiddleware())

	h.APIRoutes(rootPath)
	if h.agent.config.UI {
		h.UI(rootPath)
	}

	h.logger.WithFields(
		logrus.Fields{
			"address": h.agent.config.HTTPAddr,
		},
	).Info("api: Running HTTP server")

	go func() {
		if err := h.Engine.Run(h.agent.config.HTTPAddr); err != nil {
			h.logger.WithError(err).Error("api: Error starting HTTP server")
		}
	}()
}

// APIRoutes registers the api routes on the gin RouterGroup.
func (h *HTTPTransport) APIRoutes(
	r *gin.RouterGroup, middleware ...gin.HandlerFunc,
) {
	r.GET("/debug/vars", expvar.Handler())

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
		// Prometheus metrics scrape endpoint
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

// MetaMiddleware adds middleware to the gin Context.
func (h *HTTPTransport) MetaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Whom", h.agent.config.NodeName)
		c.Next()
	}
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
		renderJSON(c, http.StatusOK, "I am a leader")
	} else {
		renderJSON(c, http.StatusNotFound, "I am a follower")
	}
}

func (h *HTTPTransport) leaveHandler(c *gin.Context) {
	if err := h.agent.Stop(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	renderJSON(c, http.StatusOK, h.agent.peers)
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
	// metadata := c.QueryMap("metadata")
	sort := c.DefaultQuery("_sort", "id")
	if sort == "id" {
		sort = "name"
	}
	// order := c.DefaultQuery("_order", "ASC")
	// q := c.Query("q")

	pairs, err := h.agent.Store.GetAllValues()
	if err != nil {
		h.logger.WithError(err).Error("api: Unable to get pairs, store not reachable.")
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

	// Call gRPC Deletepair
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
