package runner

type PathsToResults map[string]FilesToResults

type FilesToResults map[string]OptimizersToResults

type OptimizersToResults map[string]int
