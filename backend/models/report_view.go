package models

// ReportWithDetails 举报视图（带用户信息）
type ReportWithDetails struct {
    Report
    ReporterUsername string `json:"reporter_username"`
    ReporterNickname string `json:"reporter_nickname"`
    ReporterAvatar   string `json:"reporter_avatar,omitempty"`
    HandlerUsername  string `json:"handler_username,omitempty"`
    HandlerNickname  string `json:"handler_nickname,omitempty"`
    HandlerAvatar    string `json:"handler_avatar,omitempty"`
    TargetTitle      string `json:"target_title,omitempty"`
}
