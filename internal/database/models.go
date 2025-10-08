package database

// DatabaseInfo 数据库信息
type DatabaseInfo struct {
	Name       string        `json:"name"`
	Collections []string     `json:"collections"`
	Stats      *DatabaseStats `json:"stats"`
}

// DatabaseStats 数据库统计信息
type DatabaseStats struct {
	Collections int   `json:"collections"`
	Objects     int64 `json:"objects"`
	DataSize    int64 `json:"dataSize,omitempty"`
	Indexes     int   `json:"indexes,omitempty"`
	IndexSize   int64 `json:"indexSize,omitempty"`
	StorageSize int64 `json:"storageSize,omitempty"`
}

// DocumentsResponse 文档响应
type DocumentsResponse struct {
	Documents []map[string]interface{} `json:"documents"`
	Total     int64                    `json:"total"`
	Page      int64                    `json:"page"`
	Limit     int64                    `json:"limit"`
}

// ServerStats 服务器统计信息
type ServerStats struct {
	Version      string    `json:"version"`
	Uptime       int64     `json:"uptime"`
	Connections  Connections `json:"connections"`
	Memory       Memory    `json:"mem"`
	DatabaseCount int       `json:"databaseCount"`
}

// Connections 连接信息
type Connections struct {
	Current      int `json:"current"`
	Available    int `json:"available"`
	TotalCreated int `json:"totalCreated"`
}

// Memory 内存信息
type Memory struct {
	Resident int64 `json:"resident"`
	Virtual  int64 `json:"virtual"`
	Mapped   int64 `json:"mapped,omitempty"`
}