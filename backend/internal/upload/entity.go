package upload

// CleanupResult 清理结果
type CleanupResult struct {
	ScannedFiles    int      // 扫描的文件数
	ReferencedFiles int      // 被引用的文件数
	DeletedFiles    int      // 删除的文件数
	DeletedPaths    []string // 删除的文件路径列表
	Errors          []error  // 删除失败的错误列表
}
