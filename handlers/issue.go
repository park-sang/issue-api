package handlers

import (
	"issue-api/models"
	"issue-api/store"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateIssue(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		UserID      *uint  `json:"userId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청", "code": 400})
		return
	}

	var user *models.User
	if input.UserID != nil {
		for _, u := range models.Users {
			if u.ID == *input.UserID {
				user = &u
				break
			}
		}
		if user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "존재하지 않는 사용자", "code": 400})
			return
		}
	}

	now := time.Now()

	store.Mutex.Lock()
	defer store.Mutex.Unlock()

	issue := &models.Issue{
		ID:          store.NextID,
		Title:       input.Title,
		Description: input.Description,
		Status:      "PENDING",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if user != nil {
		issue.User = user
		issue.Status = "IN_PROGRESS"
	}

	store.Issues[store.NextID] = issue
	store.NextID++

	c.JSON(http.StatusCreated, issue)
}

func ListIssues(c *gin.Context) {
	statusFilter := c.Query("status")

	var result []models.Issue

	store.Mutex.Lock()
	defer store.Mutex.Unlock()

	for _, issue := range store.Issues {
		if statusFilter == "" || issue.Status == statusFilter {
			result = append(result, *issue)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"issues": result,
	})
}

func GetIssueByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 ID입니다", "code": 400})
		return
	}

	store.Mutex.Lock()
	defer store.Mutex.Unlock()

	issue, exists := store.Issues[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "해당 이슈를 찾을 수 없습니다", "code": 404})
		return
	}

	c.JSON(http.StatusOK, issue)
}

func UpdateIssue(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 ID입니다", "code": 400})
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
		UserID      *uint   `json:"userId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청", "code": 400})
		return
	}

	validStatuses := map[string]bool{
		"PENDING":     true,
		"IN_PROGRESS": true,
		"COMPLETED":   true,
		"CANCELLED":   true,
	}

	store.Mutex.Lock()
	defer store.Mutex.Unlock()

	issue, exists := store.Issues[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "해당 이슈를 찾을 수 없습니다", "code": 404})
		return
	}

	if issue.Status == "COMPLETED" || issue.Status == "CANCELLED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "COMPLETED 또는 CANCELLED 상태의 이슈는 수정할 수 없습니다", "code": 400})
		return
	}

	// 제목 변경
	if input.Title != nil {
		issue.Title = *input.Title
	}

	// 설명 변경
	if input.Description != nil {
		issue.Description = *input.Description
	}

	// 상태 변경 유효성 체크 및 적용
	if input.Status != nil {
		if !validStatuses[*input.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 상태값입니다", "code": 400})
			return
		}

		// 담당자가 없으면 PENDING, CANCELLED 외 상태 변경 불가
		if issue.User == nil && *input.Status != "PENDING" && *input.Status != "CANCELLED" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "담당자가 없는 상태에서 PENDING 또는 CANCELLED 외 상태로 변경할 수 없습니다", "code": 400})
			return
		}

		issue.Status = *input.Status
	}

	// 담당자 변경
	if input.UserID != nil {
		if *input.UserID == 0 {
			// 담당자 제거: 상태 PENDING으로 변경
			issue.User = nil
			issue.Status = "PENDING"
		} else {
			// 담당자 할당
			var newUser *models.User
			for _, u := range models.Users {
				if u.ID == *input.UserID {
					newUser = &u
					break
				}
			}
			if newUser == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "존재하지 않는 사용자", "code": 400})
				return
			}
			issue.User = newUser
			// 상태가 PENDING일 때, 담당자 할당 시 상태 변경
			if issue.Status == "PENDING" {
				issue.Status = "IN_PROGRESS"
			}
		}
	}

	issue.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, issue)
}
