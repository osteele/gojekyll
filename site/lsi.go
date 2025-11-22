package site

import (
	"math"
	"regexp"
	"strings"
)

// LSIDocument represents a post with its TF-IDF vector
type LSIDocument struct {
	page   Page
	vector map[string]float64
	norm   float64 // cached magnitude for performance
}

// tokenize extracts words from text
func tokenize(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Extract words (alphanumeric sequences)
	re := regexp.MustCompile(`\b[a-z0-9]+\b`)
	words := re.FindAllString(text, -1)

	// Filter out very short words and common stop words
	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	}

	var filtered []string
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// termFrequency calculates TF for a document
func termFrequency(tokens []string) map[string]float64 {
	tf := make(map[string]float64)
	total := float64(len(tokens))

	if total == 0 {
		return tf
	}

	// Count occurrences
	for _, token := range tokens {
		tf[token]++
	}

	// Normalize by document length
	for token := range tf {
		tf[token] = tf[token] / total
	}

	return tf
}

// inverseDocumentFrequency calculates IDF for the corpus
func inverseDocumentFrequency(documents [][]string) map[string]float64 {
	idf := make(map[string]float64)
	docCount := float64(len(documents))

	if docCount == 0 {
		return idf
	}

	// Count how many documents contain each term
	for _, tokens := range documents {
		seen := make(map[string]bool)
		for _, token := range tokens {
			if !seen[token] {
				idf[token]++
				seen[token] = true
			}
		}
	}

	// Calculate IDF
	for token, count := range idf {
		idf[token] = math.Log(docCount / count)
	}

	return idf
}

// buildTFIDFVectors creates TF-IDF vectors for all documents
func buildTFIDFVectors(pages []Page) []LSIDocument {
	if len(pages) == 0 {
		return nil
	}

	// Extract all tokens from all documents
	allTokens := make([][]string, len(pages))
	for i, page := range pages {
		// Get page front matter
		fm := page.FrontMatter()

		// Build text from title and content
		text := ""
		if title, ok := fm["title"].(string); ok {
			text = title + " "
		}
		// Note: Content may not be available yet during initialization,
		// so we primarily use the title for LSI analysis
		allTokens[i] = tokenize(text)
	}

	// Calculate IDF for the corpus
	idf := inverseDocumentFrequency(allTokens)

	// Build TF-IDF vectors
	docs := make([]LSIDocument, len(pages))
	for i, page := range pages {
		tf := termFrequency(allTokens[i])
		vector := make(map[string]float64)

		// Calculate TF-IDF
		for term, tfVal := range tf {
			vector[term] = tfVal * idf[term]
		}

		// Calculate and cache the magnitude
		var norm float64
		for _, val := range vector {
			norm += val * val
		}
		norm = math.Sqrt(norm)

		docs[i] = LSIDocument{
			page:   page,
			vector: vector,
			norm:   norm,
		}
	}

	return docs
}

// cosineSimilarity calculates the cosine similarity between two documents
func cosineSimilarity(doc1, doc2 LSIDocument) float64 {
	// Handle edge cases
	if doc1.norm == 0 || doc2.norm == 0 {
		return 0
	}

	// Calculate dot product
	var dotProduct float64
	for term, val1 := range doc1.vector {
		if val2, exists := doc2.vector[term]; exists {
			dotProduct += val1 * val2
		}
	}

	// Return cosine similarity
	return dotProduct / (doc1.norm * doc2.norm)
}

// findRelatedPosts finds the most similar posts to a given post using LSI
func findRelatedPosts(targetPage Page, allPages []Page, limit int) []Page {
	if len(allPages) <= 1 {
		return nil
	}

	// Build TF-IDF vectors for all documents
	docs := buildTFIDFVectors(allPages)

	// Find the target document
	var targetDoc LSIDocument
	targetIndex := -1
	for i, doc := range docs {
		if doc.page == targetPage {
			targetDoc = doc
			targetIndex = i
			break
		}
	}

	// If target not found, return empty
	if targetIndex == -1 {
		return nil
	}

	// Calculate similarities
	type similarity struct {
		page  Page
		score float64
	}

	similarities := make([]similarity, 0, len(docs)-1)
	for i, doc := range docs {
		if i == targetIndex {
			continue // Skip the target itself
		}
		score := cosineSimilarity(targetDoc, doc)
		similarities = append(similarities, similarity{page: doc.page, score: score})
	}

	// Sort by similarity (descending)
	// Using a simple bubble sort for small collections
	for i := 0; i < len(similarities); i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[j].score > similarities[i].score {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}

	// Return top results
	result := make([]Page, 0, limit)
	for i := 0; i < len(similarities) && i < limit; i++ {
		result = append(result, similarities[i].page)
	}

	return result
}
