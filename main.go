package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"html/template"
	"math"

	"github.com/gin-gonic/gin"
)

// 文件分享项结构
type ShareItem struct {
	ID           string    `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	UploadTime   time.Time `json:"upload_time"`
	AccessCount  int       `json:"access_count"`
}

// 全局变量
var (
	shareItems    = make(map[string]*ShareItem)
	adminPassword = "admin123"  // 硬编码的管理员口令
	uploadDir     = "./uploads" // 上传文件存储目录
)

func main() {
	// 创建上传目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("创建上传目录失败: %v\n", err)
		return
	}

	r := gin.Default()

	// 添加模板函数
	r.SetFuncMap(template.FuncMap{
		"formatFileSize": formatFileSize,
		"isImage":        isImage,
		"isVideo":        isVideo,
	})

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 静态文件服务
	r.Static("/static", "./static")

	// 路由设置
	r.GET("/", homePage)
	r.GET("/upload", uploadPage)
	r.POST("/upload", handleUpload)
	r.GET("/share/:id", viewShare)
	r.GET("/download/:id", downloadFile)
	r.GET("/admin", adminPage)
	r.POST("/admin/login", adminLogin)
	r.GET("/admin/dashboard", adminDashboard)
	r.POST("/admin/delete/:id", deleteShare)

	fmt.Printf("服务器启动在 http://localhost:8080\n")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

// 首页
func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "快速文件分享",
	})
}

// 上传页面
func uploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "上传文件",
	})
}

// 处理文件上传
func handleUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}
	defer file.Close()

	// 生成唯一ID
	id := generateID()

	// 保存文件
	filename := id + filepath.Ext(header.Filename)
	filepath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
			return
		}
	}(dst)

	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 创建分享项
	shareItem := &ShareItem{
		ID:           id,
		Filename:     filename,
		OriginalName: header.Filename,
		Size:         header.Size,
		UploadTime:   time.Now(),
		AccessCount:  0,
	}

	shareItems[id] = shareItem

	// 返回分享链接
	shareURL := fmt.Sprintf("http://%s/share/%s", c.Request.Host, id)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"share_url": shareURL,
		"id":        id,
	})
}

// 查看分享
func viewShare(c *gin.Context) {
	id := c.Param("id")
	item, exists := shareItems[id]
	if !exists {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "分享链接不存在或已过期",
		})
		return
	}

	// 增加访问计数
	item.AccessCount++

	// 检查文件是否存在
	filepath := filepath.Join(uploadDir, item.Filename)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "文件不存在",
		})
		return
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"item":         item,
		"download_url": fmt.Sprintf("/download/%s", id),
	})
}

// 下载文件
func downloadFile(c *gin.Context) {
	id := c.Param("id")
	item, exists := shareItems[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	filepath := filepath.Join(uploadDir, item.Filename)
	c.FileAttachment(filepath, item.OriginalName)
}

// 管理员页面
func adminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", gin.H{
		"title": "管理员登录",
	})
}

// 管理员登录
func adminLogin(c *gin.Context) {
	password := c.PostForm("password")
	if password == adminPassword {
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
			"error": "口令错误",
		})
	}
}

// 管理员仪表板
func adminDashboard(c *gin.Context) {
	// 转换map为slice以便在模板中使用
	var items []*ShareItem
	for _, item := range shareItems {
		items = append(items, item)
	}

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"items": items,
	})
}

// 删除分享
func deleteShare(c *gin.Context) {
	id := c.Param("id")
	item, exists := shareItems[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享项不存在"})
		return
	}

	// 删除文件
	filepath := filepath.Join(uploadDir, item.Filename)
	os.Remove(filepath)

	// 删除分享项
	delete(shareItems, id)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 生成唯一ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 格式化文件大小
func formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 Bytes"
	}

	const unit = 1024
	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	sizes := []string{"Bytes", "KB", "MB", "GB", "TB"}

	if exp >= len(sizes) {
		exp = len(sizes) - 1
	}

	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(unit, float64(exp)), sizes[exp])
}

// 检查是否为图片文件
func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}

	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// 检查是否为视频文件
func isVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mkv", ".m4v"}

	for _, videoExt := range videoExts {
		if ext == videoExt {
			return true
		}
	}
	return false
}
