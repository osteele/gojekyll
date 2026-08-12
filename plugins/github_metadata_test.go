package plugins

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/require"
)

func TestGitHubRepoCacheReusesAndExpiresSuccessfulRequests(t *testing.T) {
	cache := newGitHubRepoCache(10 * time.Minute)
	now := time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)
	var calls int
	fetch := func(nwo string) (*github.Repository, error) {
		calls++
		return &github.Repository{Name: github.String(nwo)}, nil
	}

	first, err := cache.get("osteele/gojekyll", now, fetch)
	require.NoError(t, err)
	second, err := cache.get("osteele/gojekyll", now.Add(9*time.Minute), fetch)
	require.NoError(t, err)
	require.Same(t, first, second)
	require.Equal(t, 1, calls)

	third, err := cache.get("osteele/gojekyll", now.Add(10*time.Minute), fetch)
	require.NoError(t, err)
	require.NotSame(t, first, third)
	require.Equal(t, 2, calls)
}

func TestGitHubRepoCacheDoesNotCacheFailures(t *testing.T) {
	cache := newGitHubRepoCache(10 * time.Minute)
	now := time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)
	var calls int
	fetch := func(string) (*github.Repository, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("temporary failure")
		}
		return &github.Repository{Name: github.String("gojekyll")}, nil
	}

	_, err := cache.get("osteele/gojekyll", now, fetch)
	require.EqualError(t, err, "temporary failure")
	_, err = cache.get("osteele/gojekyll", now, fetch)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestGitHubRepoCacheCoalescesConcurrentRequests(t *testing.T) {
	cache := newGitHubRepoCache(10 * time.Minute)
	now := time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)
	started := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int32
	fetch := func(string) (*github.Repository, error) {
		calls.Add(1)
		close(started)
		<-release
		return &github.Repository{Name: github.String("gojekyll")}, nil
	}

	const requests = 8
	var wg sync.WaitGroup
	results := make(chan *github.Repository, requests)
	errs := make(chan error, requests)
	for range requests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo, err := cache.get("osteele/gojekyll", now, fetch)
			results <- repo
			errs <- err
		}()
	}
	<-started
	close(release)
	wg.Wait()
	close(results)
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
	var first *github.Repository
	for repo := range results {
		if first == nil {
			first = repo
		}
		require.Same(t, first, repo)
	}
	require.Equal(t, int32(1), calls.Load())
}
