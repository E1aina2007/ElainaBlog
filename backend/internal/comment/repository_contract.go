// repository_contract.go 定义 comment 仓储层接口，用于依赖注入和测试 mock
package comment

// Repository 接口封装评论数据访问操作。
type Repository interface {
	GetCommentByID(id int64) (*Comment, error)
	GetCommentListByArticleID(articleID int64) ([]*CommentVO, error)
	GetAllCommentList() ([]*CommentVO, error)
	CreateComment(comment *Comment) (int64, error)
	DeleteComment(id int64) error
	DeleteCommentsByArticleID(articleID int64) error
}
