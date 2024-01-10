package funcs

type Post_json struct {
	User_ID       string   `json:"user_id"`
	User_Name     string   `json:"user_name"`
	Title         string   `json:"title"`
	Text          string   `json:"text"`
	Category      []string `json:"category"`
	Edited        bool     `json:"edited"`
	Creation_Date string   `json:"creation_date"`
	PostLikes     int      `json:"post_likes"`
	PostDisLikes  int      `json:"post_dislikes"`
	Post_ID       string   `json:"post_id"`
	IsLiked       int      `json:"isLiked"`
}
