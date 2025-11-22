package site

import (
	"github.com/osteele/gojekyll/collection"
)

func (s *Site) findPostCollection() *collection.Collection {
	for _, c := range s.Collections {
		if c.Name == "posts" {
			return c
		}
	}
	return nil
}

func (s *Site) setPostVariables() {
	c := s.findPostCollection()
	if c == nil {
		return
	}
	var (
		ps      = c.Pages()
		related = ps
	)

	// Use LSI to find related posts if enabled
	if s.cfg.LSI && len(ps) > 0 {
		// Calculate average similarity for each post to find most "central" posts
		related = s.findCentralPosts(ps, 10)
	} else {
		// Default behavior: return most recent posts
		if len(related) > 10 {
			related = related[:10]
		}
	}

	s.drop["categories"] = s.groupPagesBy(func(p Page) []string { return p.Categories() })
	s.drop["tags"] = s.groupPagesBy(func(p Page) []string { return p.Tags() })
	s.drop["related_posts"] = related
}

// findCentralPosts finds posts with highest average similarity to all other posts
func (s *Site) findCentralPosts(posts []Page, limit int) []Page {
	if len(posts) == 0 {
		return nil
	}
	if len(posts) <= limit {
		return posts
	}

	// Build TF-IDF vectors for all posts
	docs := buildTFIDFVectors(posts)

	// Calculate average similarity for each post
	type postScore struct {
		page  Page
		score float64
	}

	scores := make([]postScore, len(docs))
	for i, doc := range docs {
		var totalSim float64
		for j, otherDoc := range docs {
			if i != j {
				totalSim += cosineSimilarity(doc, otherDoc)
			}
		}
		avgSim := totalSim / float64(len(docs)-1)
		scores[i] = postScore{page: doc.page, score: avgSim}
	}

	// Sort by average similarity (descending)
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score > scores[i].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Return top results
	result := make([]Page, 0, limit)
	for i := 0; i < len(scores) && i < limit; i++ {
		result = append(result, scores[i].page)
	}

	return result
}

func (s *Site) groupPagesBy(getter func(Page) []string) map[string][]Page {
	categories := map[string][]Page{}
	for _, p := range s.Pages() {
		for _, k := range p.Categories() {
			ps, found := categories[k]
			if !found {
				ps = []Page{}
			}
			categories[k] = append(ps, p)
		}
	}
	return categories
}
