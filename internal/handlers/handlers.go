package handlers

import (
	"context"
	"m-db-ui/internal/config"
	"m-db-ui/internal/database"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	dbService         *database.Service
	connectionManager *config.ConnectionManager
}

func New(dbService *database.Service, connectionManager *config.ConnectionManager) *Handlers {
	return &Handlers{
		dbService:         dbService,
		connectionManager: connectionManager,
	}
}

// Index 首页
func (h *Handlers) Index(c *gin.Context) {
	databases, err := h.dbService.GetDatabases()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	currentConnection := h.connectionManager.GetCurrentConnection()

	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":            "MongoDB管理工具",
		"databases":        databases,
		"currentConnection": currentConnection,
	})
}

// ConnectionsPage 连接管理页面
func (h *Handlers) ConnectionsPage(c *gin.Context) {
	connections := h.connectionManager.GetConnections()
	currentId := h.connectionManager.GetCurrentID()

	c.HTML(http.StatusOK, "connections.html", gin.H{
		"title":       "连接管理",
		"connections": connections,
		"currentId":   currentId,
	})
}

// DatabasePage 数据库页面
func (h *Handlers) DatabasePage(c *gin.Context) {
	dbName := c.Param("db")

	dbInfo, err := h.dbService.GetDatabase(dbName)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":  dbName + " - 数据库管理",
		"dbInfo": dbInfo,
	})
}

// CollectionPage 集合页面
func (h *Handlers) CollectionPage(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	documents, err := h.dbService.GetDocuments(dbName, collectionName, page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 计算翻页数据
	totalPages := (documents.Total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	// 生成页码数组
	var pageNumbers []int64
	start := page - 2
	if start < 1 {
		start = 1
	}
	end := page + 2
	if end > totalPages {
		end = totalPages
	}

	for i := start; i <= end; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	// 计算当前页的记录范围
	startRecord := (page-1)*limit + 1
	endRecord := page * limit
	if endRecord > documents.Total {
		endRecord = documents.Total
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":       collectionName + " - 集合管理",
		"dbName":      dbName,
		"collection":  collectionName,
		"documents":   documents.Documents,
		"total":       documents.Total,
		"page":        documents.Page,
		"limit":       documents.Limit,
		"totalPages":  totalPages,
		"pageNumbers": pageNumbers,
		"start":       startRecord,
		"end":         endRecord,
	})
}

// GetDatabases 获取所有数据库
func (h *Handlers) GetDatabases(c *gin.Context) {
	databases, err := h.dbService.GetDatabases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, databases)
}

// GetDatabase 获取数据库信息
func (h *Handlers) GetDatabase(c *gin.Context) {
	dbName := c.Param("name")

	dbInfo, err := h.dbService.GetDatabase(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dbInfo)
}

// DeleteDatabase 删除数据库
func (h *Handlers) DeleteDatabase(c *gin.Context) {
	dbName := c.Param("name")

	err := h.dbService.DeleteDatabase(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database deleted successfully"})
}

// CreateDatabase 创建数据库
func (h *Handlers) CreateDatabase(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.dbService.CreateDatabase(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database created successfully"})
}

// GetCollections 获取数据库的所有集合
func (h *Handlers) GetCollections(c *gin.Context) {
	dbName := c.Param("db")

	collections, err := h.dbService.GetCollections(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, collections)
}

// CreateCollection 创建集合
func (h *Handlers) CreateCollection(c *gin.Context) {
	dbName := c.Param("db")

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.dbService.CreateCollection(dbName, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection created successfully"})
}

// DeleteCollection 删除集合
func (h *Handlers) DeleteCollection(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")

	err := h.dbService.DeleteCollection(dbName, collectionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}

// GetDocuments 获取集合中的文档
func (h *Handlers) GetDocuments(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	documents, err := h.dbService.GetDocuments(dbName, collectionName, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, documents)
}

// GetDocument 获取单个文档
func (h *Handlers) GetDocument(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")
	id := c.Param("id")

	// URL解码
	decodedID, err := url.QueryUnescape(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 移除ObjectID()包装，如果存在
	if len(decodedID) > 9 && decodedID[:9] == "ObjectID(" && decodedID[len(decodedID)-1] == ')' {
		decodedID = decodedID[9 : len(decodedID)-1]
		// 移除可能的引号
		decodedID = strings.Trim(decodedID, `"`)
	}

	document, err := h.dbService.GetDocument(dbName, collectionName, decodedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, document)
}

// CreateDocument 创建文档
func (h *Handlers) CreateDocument(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")

	var document map[string]interface{}
	if err := c.ShouldBindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.dbService.CreateDocument(dbName, collectionName, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// UpdateDocument 更新文档
func (h *Handlers) UpdateDocument(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")
	id := c.Param("id")

	var document map[string]interface{}
	if err := c.ShouldBindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.dbService.UpdateDocument(dbName, collectionName, id, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// DeleteDocument 删除文档
func (h *Handlers) DeleteDocument(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")
	id := c.Param("id")

	err := h.dbService.DeleteDocument(dbName, collectionName, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// QueryDocuments 查询文档
func (h *Handlers) QueryDocuments(c *gin.Context) {
	dbName := c.Param("db")
	collectionName := c.Param("collection")

	var req struct {
		Query map[string]interface{} `json:"query"`
		Page  int64                  `json:"page"`
		Limit int64                  `json:"limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	documents, err := h.dbService.QueryDocuments(dbName, collectionName, req.Query, req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, documents)
}

// GetStats 获取统计信息
func (h *Handlers) GetStats(c *gin.Context) {
	stats, err := h.dbService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// 连接管理相关处理器

// GetConnections 获取所有连接配置
func (h *Handlers) GetConnections(c *gin.Context) {
	connections := h.connectionManager.GetConnections()
	c.JSON(http.StatusOK, gin.H{
		"connections": connections,
		"currentId":    h.connectionManager.GetCurrentID(),
	})
}

// GetConnection 获取指定连接配置
func (h *Handlers) GetConnection(c *gin.Context) {
	id := c.Param("id")
	connection, err := h.connectionManager.GetConnection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, connection)
}

// AddConnection 添加连接配置
func (h *Handlers) AddConnection(c *gin.Context) {
	var config config.ConnectionConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.connectionManager.AddConnection(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection added successfully", "id": config.ID})
}

// UpdateConnection 更新连接配置
func (h *Handlers) UpdateConnection(c *gin.Context) {
	id := c.Param("id")
	var config config.ConnectionConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.connectionManager.UpdateConnection(id, &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection updated successfully"})
}

// DeleteConnection 删除连接配置
func (h *Handlers) DeleteConnection(c *gin.Context) {
	id := c.Param("id")
	if err := h.connectionManager.DeleteConnection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Connection deleted successfully"})
}

// SetCurrentConnection 设置当前连接
func (h *Handlers) SetCurrentConnection(c *gin.Context) {
	id := c.Param("id")
	if err := h.connectionManager.SetCurrentConnection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Current connection set successfully"})
}

// GetCurrentConnection 获取当前连接
func (h *Handlers) GetCurrentConnection(c *gin.Context) {
	connection := h.connectionManager.GetCurrentConnection()
	if connection == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No current connection"})
		return
	}
	c.JSON(http.StatusOK, connection)
}

// TestConnection 测试连接
func (h *Handlers) TestConnection(c *gin.Context) {
	var config config.ConnectionConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 测试连接
	client, err := database.Connect(config.GetURI())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect: " + err.Error()})
		return
	}
	defer client.Disconnect(context.Background())

	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}