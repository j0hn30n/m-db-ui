package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// ConnectionConfig 数据库连接配置
type ConnectionConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Database    string `json:"database"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	AuthDB      string `json:"authDB,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

// GetURI 获取MongoDB连接URI
func (c *ConnectionConfig) GetURI() string {
	if c.Username != "" && c.Password != "" {
		authDB := c.AuthDB
		if authDB == "" {
			authDB = c.Database
			if authDB == "" {
				authDB = "admin"
			}
		}
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			c.Username, c.Password, c.Host, c.Port, c.Database, authDB)
	}
	return fmt.Sprintf("mongodb://%s:%d/%s", c.Host, c.Port, c.Database)
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*ConnectionConfig
	currentID   string
	mutex       sync.RWMutex
	filePath    string
}

var GlobalConnectionManager *ConnectionManager

// NewConnectionManager 创建连接管理器
func NewConnectionManager(filePath string) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*ConnectionConfig),
		filePath:    filePath,
	}
}

// LoadConnections 加载连接配置
func (cm *ConnectionManager) LoadConnections() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, err := os.Stat(cm.filePath); os.IsNotExist(err) {
		// 创建默认配置
		defaultConfig := &ConnectionConfig{
			ID:          "default",
			Name:        "本地MongoDB",
			Host:        "localhost",
			Port:        27017,
			Database:    "",
			Description: "默认本地MongoDB连接",
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}
		cm.connections[defaultConfig.ID] = defaultConfig
		cm.currentID = defaultConfig.ID
		return cm.saveConnections()
	}

	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		return err
	}

	var connections []*ConnectionConfig
	if err := json.Unmarshal(data, &connections); err != nil {
		return err
	}

	cm.connections = make(map[string]*ConnectionConfig)
	for _, conn := range connections {
		cm.connections[conn.ID] = conn
		if cm.currentID == "" {
			cm.currentID = conn.ID
		}
	}

	return nil
}

// SaveConnections 保存连接配置
func (cm *ConnectionManager) saveConnections() error {
	var connections []*ConnectionConfig
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}

	data, err := json.MarshalIndent(connections, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.filePath, data, 0644)
}

// AddConnection 添加连接
func (cm *ConnectionManager) AddConnection(config *ConnectionConfig) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if config.ID == "" {
		config.ID = fmt.Sprintf("conn_%d", time.Now().Unix())
	}
	config.CreatedAt = time.Now().Unix()
	config.UpdatedAt = time.Now().Unix()

	cm.connections[config.ID] = config
	return cm.saveConnections()
}

// UpdateConnection 更新连接
func (cm *ConnectionManager) UpdateConnection(id string, config *ConnectionConfig) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.connections[id]; !exists {
		return fmt.Errorf("connection not found")
	}

	config.ID = id
	config.UpdatedAt = time.Now().Unix()
	cm.connections[id] = config
	return cm.saveConnections()
}

// DeleteConnection 删除连接
func (cm *ConnectionManager) DeleteConnection(id string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.connections[id]; !exists {
		return fmt.Errorf("connection not found")
	}

	delete(cm.connections, id)

	// 如果删除的是当前连接，切换到第一个可用连接
	if cm.currentID == id {
		for connID := range cm.connections {
			cm.currentID = connID
			break
		}
	}

	return cm.saveConnections()
}

// GetConnections 获取所有连接
func (cm *ConnectionManager) GetConnections() []*ConnectionConfig {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var connections []*ConnectionConfig
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}
	return connections
}

// GetConnection 获取指定连接
func (cm *ConnectionManager) GetConnection(id string) (*ConnectionConfig, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	conn, exists := cm.connections[id]
	if !exists {
		return nil, fmt.Errorf("connection not found")
	}
	return conn, nil
}

// GetCurrentConnection 获取当前连接
func (cm *ConnectionManager) GetCurrentConnection() *ConnectionConfig {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if conn, exists := cm.connections[cm.currentID]; exists {
		return conn
	}
	return nil
}

// SetCurrentConnection 设置当前连接
func (cm *ConnectionManager) SetCurrentConnection(id string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.connections[id]; !exists {
		return fmt.Errorf("connection not found")
	}

	cm.currentID = id
	return nil
}

// GetCurrentID 获取当前连接ID
func (cm *ConnectionManager) GetCurrentID() string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.currentID
}