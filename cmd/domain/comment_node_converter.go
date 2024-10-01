package domain

import "github.com/kkwon1/apod-forum-backend/cmd/models"

func ConvertToCommentNodes(comments []models.Comment) []*models.CommentNode {
	commentMap := make(map[string]*models.CommentNode)
	for _, comment := range comments {
		commentNode := &models.CommentNode{
			CommentID: comment.CommentID,
			Author:    comment.Author,
			Comment:   comment.Comment,
			Children:  []*models.CommentNode{},
		}
		commentMap[comment.CommentID] = commentNode
	}

	pidToCidMap := make(map[string][]string)
	for _, comment := range comments {
		parentID := comment.ParentID
		pidToCidMap[parentID] = append(pidToCidMap[parentID], comment.CommentID)
	}

	result := []*models.CommentNode{}

	for parentID, childIDs := range pidToCidMap {
		// No parent ID, so we attach as root children
		for _, childID := range childIDs {
			if parentID == "" {
				result = append(result, commentMap[childID])
			} else {
				parentNode := commentMap[parentID]
				parentNode.Children = append(parentNode.Children, commentMap[childID])
			}
		}
	}

	return result
}