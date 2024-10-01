package domain

import (
	"reflect"
	"testing"

	"github.com/kkwon1/apod-forum-backend/cmd/models"
)

func TestConvertToCommentNodes(t *testing.T) {
	tests := []struct {
		name     string
		comments []models.Comment
		expected []*models.CommentNode
	}{
		{
			name: "single comment",
			comments: []models.Comment{
				{CommentID: "1", Author: "Author1", Comment: "Comment1", ParentID: ""},
			},
			expected: []*models.CommentNode{
				{
					CommentID: "1",
					Author:    "Author1",
					Comment:   "Comment1",
					Children:  []*models.CommentNode{},
				},
			},
		},
		{
			name: "multiple root comments",
			comments: []models.Comment{
				{CommentID: "1", Author: "Author1", Comment: "Comment1", ParentID: ""},
				{CommentID: "2", Author: "Author2", Comment: "Comment2", ParentID: ""},
			},
			expected: []*models.CommentNode{
				{
					CommentID: "1",
					Author:    "Author1",
					Comment:   "Comment1",
					Children:  []*models.CommentNode{},
				},
				{
					CommentID: "2",
					Author:    "Author2",
					Comment:   "Comment2",
					Children:  []*models.CommentNode{},
				},
			},
		},
		{
			name: "nested comments",
			comments: []models.Comment{
				{CommentID: "1", Author: "Author1", Comment: "Comment1", ParentID: ""},
				{CommentID: "2", Author: "Author2", Comment: "Comment2", ParentID: "1"},
				{CommentID: "3", Author: "Author3", Comment: "Comment3", ParentID: "2"},
			},
			expected: []*models.CommentNode{
				{
					CommentID: "1",
					Author:    "Author1",
					Comment:   "Comment1",
					Children: []*models.CommentNode{
						{
							CommentID: "2",
							Author:    "Author2",
							Comment:   "Comment2",
							Children: []*models.CommentNode{
								{
									CommentID: "3",
									Author:    "Author3",
									Comment:   "Comment3",
									Children:  []*models.CommentNode{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple nested comments",
			comments: []models.Comment{
				{CommentID: "1", Author: "Author1", Comment: "Comment1", ParentID: ""},
				{CommentID: "2", Author: "Author2", Comment: "Comment2", ParentID: "1"},
				{CommentID: "3", Author: "Author3", Comment: "Comment3", ParentID: "1"},
				{CommentID: "4", Author: "Author4", Comment: "Comment4", ParentID: "2"},
			},
			expected: []*models.CommentNode{
				{
					CommentID: "1",
					Author:    "Author1",
					Comment:   "Comment1",
					Children: []*models.CommentNode{
						{
							CommentID: "2",
							Author:    "Author2",
							Comment:   "Comment2",
							Children: []*models.CommentNode{
								{
									CommentID: "4",
									Author:    "Author4",
									Comment:   "Comment4",
									Children:  []*models.CommentNode{},
								},
							},
						},
						{
							CommentID: "3",
							Author:    "Author3",
							Comment:   "Comment3",
							Children:  []*models.CommentNode{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToCommentNodes(tt.comments)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertToCommentNodes() = %v, expected %v", result, tt.expected)
			}
		})
	}
}