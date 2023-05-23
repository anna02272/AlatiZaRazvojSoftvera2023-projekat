package helper

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	posts       = "posts/%s/%s"
	postsLabels = "posts/%s/%s/%s"
	all         = "posts"
)

func generateKey(version string, labels string) (string, string) {
	id := uuid.New().String()
	if labels != "" {
		return fmt.Sprintf(postsLabels, id, version, labels), id
	} else {
		return fmt.Sprintf(posts, id, version), id
	}

}

func constructKey(id string, version string, labels string) string {
	if labels != "" {
		return fmt.Sprintf(postsLabels, id, version, labels)
	} else {
		return fmt.Sprintf(posts, id, version)
	}

}
