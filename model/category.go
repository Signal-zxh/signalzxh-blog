package model

// Category 分类结构
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Tag 标签结构
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PostWithCategoryTag 带分类和标签的文章结构
type PostWithCategoryTag struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	UserID   int      `json:"user_id"`
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}
