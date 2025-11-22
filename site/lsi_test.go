package site

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple text",
			input:    "This is a test",
			expected: []string{"this", "test"},
		},
		{
			name:     "text with punctuation",
			input:    "Hello, world! How are you?",
			expected: []string{"hello", "world", "how", "you"},
		},
		{
			name:     "text with numbers",
			input:    "Test 123 with numbers",
			expected: []string{"test", "123", "numbers"}, // "with" is a stop word
		},
		{
			name:     "filters stop words",
			input:    "the cat is on the mat",
			expected: []string{"cat", "mat"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil, // returns nil for empty input
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenize(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTermFrequency(t *testing.T) {
	tokens := []string{"cat", "dog", "cat", "bird"}
	tf := termFrequency(tokens)

	require.InDelta(t, 0.5, tf["cat"], 0.01)   // 2/4 = 0.5
	require.InDelta(t, 0.25, tf["dog"], 0.01)  // 1/4 = 0.25
	require.InDelta(t, 0.25, tf["bird"], 0.01) // 1/4 = 0.25
}

func TestTermFrequencyEmpty(t *testing.T) {
	tokens := []string{}
	tf := termFrequency(tokens)

	require.Empty(t, tf)
}

func TestInverseDocumentFrequency(t *testing.T) {
	documents := [][]string{
		{"cat", "dog"},
		{"cat", "bird"},
		{"dog", "bird"},
	}

	idf := inverseDocumentFrequency(documents)

	// "cat" appears in 2 out of 3 documents: log(3/2) ≈ 0.405
	require.InDelta(t, 0.405, idf["cat"], 0.01)

	// "dog" appears in 2 out of 3 documents: log(3/2) ≈ 0.405
	require.InDelta(t, 0.405, idf["dog"], 0.01)

	// "bird" appears in 2 out of 3 documents: log(3/2) ≈ 0.405
	require.InDelta(t, 0.405, idf["bird"], 0.01)
}

func TestCosineSimilarity(t *testing.T) {
	// Two identical vectors should have similarity of 1.0
	doc1 := LSIDocument{
		vector: map[string]float64{"cat": 1.0, "dog": 1.0},
		norm:   1.414, // sqrt(1^2 + 1^2)
	}
	doc2 := LSIDocument{
		vector: map[string]float64{"cat": 1.0, "dog": 1.0},
		norm:   1.414,
	}

	similarity := cosineSimilarity(doc1, doc2)
	require.InDelta(t, 1.0, similarity, 0.01)
}

func TestCosineSimilarityOrthogonal(t *testing.T) {
	// Two orthogonal vectors should have similarity of 0.0
	doc1 := LSIDocument{
		vector: map[string]float64{"cat": 1.0},
		norm:   1.0,
	}
	doc2 := LSIDocument{
		vector: map[string]float64{"dog": 1.0},
		norm:   1.0,
	}

	similarity := cosineSimilarity(doc1, doc2)
	require.InDelta(t, 0.0, similarity, 0.01)
}

func TestCosineSimilarityZeroNorm(t *testing.T) {
	// Vector with zero norm should return 0
	doc1 := LSIDocument{
		vector: map[string]float64{"cat": 1.0},
		norm:   0.0,
	}
	doc2 := LSIDocument{
		vector: map[string]float64{"cat": 1.0},
		norm:   1.0,
	}

	similarity := cosineSimilarity(doc1, doc2)
	require.Equal(t, 0.0, similarity)
}
