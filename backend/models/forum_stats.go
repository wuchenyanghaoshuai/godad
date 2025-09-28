package models

// ForumStats 论坛统计数据
type ForumStats struct {
    ForumPostCount    int64 `json:"forumPostCount"`
    ForumReplyCount   int64 `json:"forumReplyCount"`
    ActiveUserCount   int64 `json:"activeUserCount"`
    FeaturedPostCount int64 `json:"featuredPostCount"`
    PinnedPostCount   int64 `json:"pinnedPostCount"`
    TodayNewPosts     int64 `json:"todayNewPosts"`
    TodayNewReplies   int64 `json:"todayNewReplies"`
}

