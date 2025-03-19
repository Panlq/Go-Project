package upload

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type ChunkInfo struct {
	Index    int    `json:"index"`
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
}

type MergeInfo struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Size     int64  `json:"size"`
}

var (
	uploadDir     = "data"
	chunkDir      = filepath.Join(uploadDir, "chunks")
	uploadingLock sync.Map
)

// 获取文件的存储路径
func getFilePath(hash, filename string) string {
	return filepath.Join(uploadDir, hash, filename)
}

// 获取分片状态文件路径
func getStatusFilePath(hash string) string {
	return filepath.Join(chunkDir, hash, "status")
}

// 读取分片状态
func readChunkStatus(hash string) ([]int, error) {
	statusPath := getStatusFilePath(hash)
	if _, err := os.Stat(statusPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(statusPath)
	if err != nil {
		return nil, err
	}

	var chunks []int
	for i := 0; i < len(data); i++ {
		if data[i] == 1 {
			chunks = append(chunks, i)
		}
	}
	return chunks, nil
}

// 更新分片状态
func updateChunkStatus(hash string, index int) error {
	statusPath := getStatusFilePath(hash)

	// 打开状态文件
	f, err := os.OpenFile(statusPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// 将对应位置设为1
	if _, err := f.Seek(int64(index), 0); err != nil {
		return err
	}
	if _, err := f.Write([]byte{1}); err != nil {
		return err
	}

	return nil
}

func init() {
	// 创建必要的目录
	for _, dir := range []string{uploadDir, chunkDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}
}

// 检查文件是否已存在，并返回已上传的分片信息
func HandleCheck(c *fiber.Ctx) error {
	hash := c.Query("hash")
	filename := c.Query("filename")
	if hash == "" || filename == "" {
		return c.Status(400).JSON(fiber.Map{"error": "hash and filename are required"})
	}

	filePath := getFilePath(hash, filename)
	_, err := os.Stat(filePath)
	exists := !os.IsNotExist(err)
	if exists {
		// 如果文件已存在，返回成功
		return c.JSON(fiber.Map{
			"exists":         true,
			"uploadedChunks": nil,
		})
	}

	// 获取前端传入的总分片数
	totalChunks, err := strconv.Atoi(c.Query("totalChunks"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid totalChunks"})
	}

	// 创建分片目录和状态文件
	chunkPath := filepath.Join(chunkDir, hash)
	if err := os.MkdirAll(chunkPath, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create chunk directory"})
	}

	// 初始化状态文件
	statusPath := getStatusFilePath(hash)
	if _, err := os.Stat(statusPath); os.IsNotExist(err) {
		if err := os.WriteFile(statusPath, make([]byte, totalChunks), 0644); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to create status file"})
		}
	}

	// 获取已上传的分片信息
	uploadedChunks, err := readChunkStatus(hash)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read chunk status"})
	}

	return c.JSON(fiber.Map{
		"exists":         exists,
		"uploadedChunks": uploadedChunks,
	})
}

// 处理分片上传
func HandleUploadChunk(c *fiber.Ctx) error {
	// 设置请求超时

	file, err := c.FormFile("chunk")
	if err != nil {
		fmt.Printf("[ERROR] Failed to get chunk file: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "chunk is required"})
	}
	fmt.Printf("[INFO] Received chunk file: %s, size: %d\n", file.Filename, file.Size)

	hash := c.FormValue("hash")
	if hash == "" {
		return c.Status(400).JSON(fiber.Map{"error": "hash is required"})
	}

	// 创建分片目录
	chunkPath := filepath.Join(chunkDir, hash)
	if err := os.MkdirAll(chunkPath, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create chunk directory"})
	}

	// 保存分片文件
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to open uploaded file"})
	}
	defer src.Close()

	index := c.FormValue("index")
	chunkIndex, err := strconv.Atoi(index)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid chunk index"})
	}

	dst, err := os.Create(filepath.Join(chunkPath, fmt.Sprintf("%d", chunkIndex)))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create chunk file"})
	}
	defer dst.Close()

	// 使用带缓冲的写入
	bufWriter := bufio.NewWriter(dst)
	if _, err = io.Copy(bufWriter, src); err != nil {
		if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "EPIPE") {
			return c.Status(500).JSON(fiber.Map{"error": "connection interrupted, please retry"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to save chunk file"})
	}

	// 确保所有数据都写入磁盘
	if err = bufWriter.Flush(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to flush chunk file"})
	}

	// 更新分片状态
	if err := updateChunkStatus(hash, chunkIndex); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update chunk status"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// 合并文件分片
func HandleMerge(c *fiber.Ctx) error {
	var info MergeInfo
	if err := c.BodyParser(&info); err != nil {
		fmt.Printf("[ERROR] Failed to parse merge request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	fmt.Printf("[INFO] Starting merge process for file: %s, hash: %s, size: %d\n", info.Filename, info.Hash, info.Size)

	if info.Hash == "" || info.Filename == "" {
		return c.Status(400).JSON(fiber.Map{"error": "hash and filename are required"})
	}

	// 检查是否有其他合并操作正在进行
	if _, exists := uploadingLock.LoadOrStore(info.Hash, true); exists {
		return c.Status(409).JSON(fiber.Map{"error": "file is being processed"})
	}
	defer uploadingLock.Delete(info.Hash)

	// 获取所有分片
	chunkPath := filepath.Join(chunkDir, info.Hash)
	if _, err := os.Stat(chunkPath); os.IsNotExist(err) {
		return c.Status(400).JSON(fiber.Map{"error": "no chunks found"})
	}

	chunks, err := os.ReadDir(chunkPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read chunks"})
	}

	if len(chunks) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "no chunks found"})
	}

	// 按索引排序分片，过滤掉status文件
	chunkFiles := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		if !chunk.IsDir() && chunk.Name() != "status" {
			chunkFiles = append(chunkFiles, chunk.Name())
		}
	}

	sort.Slice(chunkFiles, func(i, j int) bool {
		var ni, nj int
		fmt.Sscanf(chunkFiles[i], "%d", &ni)
		fmt.Sscanf(chunkFiles[j], "%d", &nj)
		return ni < nj
	})

	// 创建目标文件，使用原始文件名
	// 创建hash目录
	fileDir := filepath.Join(uploadDir, info.Hash)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create file directory"})
	}

	// 创建目标文件
	dstPath := getFilePath(info.Hash, info.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create merged file"})
	}
	defer dst.Close()

	// 合并分片
	hasher := md5.New()

	// 使用MultiWriter将数据同时写入文件和hasher
	writer := io.MultiWriter(dst, hasher)

	// 合并所有分片
	for i, chunkFile := range chunkFiles {
		chunkFilePath := filepath.Join(chunkPath, chunkFile)
		src, err := os.Open(chunkFilePath)
		if err != nil {
			fmt.Printf("[ERROR] Failed to open chunk file %s: %v\n", chunkFilePath, err)
			os.Remove(dstPath)
			return c.Status(500).JSON(fiber.Map{"error": "failed to open chunk file"})
		}

		// 使用io.Copy进行文件复制
		written, err := io.Copy(writer, src)
		src.Close()
		if err != nil {
			fmt.Printf("[ERROR] Failed to copy chunk data from %s: %v\n", chunkFilePath, err)
			os.Remove(dstPath)
			return c.Status(500).JSON(fiber.Map{"error": "failed to copy chunk data"})
		}

		fmt.Printf("[INFO] Processed chunk %d/%d: size=%d\n", i+1, len(chunkFiles), written)
	}

	// 验证文件完整性
	fileHash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("[INFO] Verifying file hash: expected %s, actual %s\n", info.Hash, fileHash)
	if fileHash != info.Hash {
		fmt.Printf("[ERROR] File hash verification failed for %s\n", info.Filename)
		os.Remove(dstPath)
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("file hash mismatch: expected %s, got %s", info.Hash, fileHash)})
	}
	fmt.Printf("[INFO] File %s merged successfully\n", info.Filename)

	// 清理分片文件和状态
	os.RemoveAll(chunkPath)

	return c.JSON(fiber.Map{"success": true})
}
