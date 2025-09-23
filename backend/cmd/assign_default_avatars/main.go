package main

import (
    "fmt"
    "log"
    "strings"

    "godad-backend/config"
    "godad-backend/models"
    "godad-backend/services"

    "github.com/joho/godotenv"
)

// backfill empty avatars using the same rule as registration
func main() {
    // Load .env if present (try repo root and backend dir)
    if err := godotenv.Load(".env"); err != nil {
        _ = godotenv.Load("backend/.env")
    }

    // Init config + DB
    cfg := config.LoadConfig()
    if err := config.InitDatabase(cfg); err != nil {
        log.Fatalf("init db failed: %v", err)
    }
    defer config.CloseDatabase()
    db := config.GetDB()
    svc := services.NewUserService()

    var users []models.User
    if err := db.Where("(avatar = '' OR avatar IS NULL) AND status = 1").Find(&users).Error; err != nil {
        log.Fatalf("query users: %v", err)
    }
    if len(users) == 0 {
        fmt.Println("No users need backfill.")
        return
    }
    updated := 0
    for _, u := range users {
        url := svc.ComputeDefaultAvatarURL(u.Username, strings.TrimSpace(u.Nickname))
        if url == "" {
            continue
        }
        if err := db.Model(&u).Update("avatar", url).Error; err != nil {
            log.Printf("update user %d failed: %v", u.ID, err)
            continue
        }
        updated++
    }
    fmt.Printf("Backfill done. Updated %d users.\n", updated)
}
