package runner

// PathsToResults represents the mapping between paths and FilesToResults.
type PathsToResults map[string]FilesToResults

// FilesToResults represents the mapping between files and OptimizersToResults.
type FilesToResults map[string]OptimizersToResults

// OptimizersToResults represents the mapping between optimizer name
// and optimization value got from this optimizer.
type OptimizersToResults map[string]int
