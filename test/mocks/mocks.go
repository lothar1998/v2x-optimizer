package mocks

//v2x-optimizer
//go:generate mockgen --build_flags=-mod=mod -destination=data/decoder.go -package=mocks github.com/lothar1998/v2x-optimizer/pkg/data EncoderDecoder

//v2x-optimizer-performance
//go:generate mockgen --build_flags=-mod=mod -destination=performance/executor/process.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/executor Process
//go:generate mockgen --build_flags=-mod=mod -destination=performance/executor/executor.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/executor Executor
//go:generate mockgen --build_flags=-mod=mod -destination=performance/optimizer/optimizer.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/optimizer IdentifiableOptimizer
//go:generate mockgen --build_flags=-mod=mod -destination=performance/runner/file/runner.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/runner FileRunner
//go:generate mockgen --build_flags=-mod=mod -destination=performance/runner/path/runner.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/runner PathRunner
//go:generate mockgen --build_flags=-mod=mod -destination=performance/runner/view/view.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/runner/view DirectoryView
//go:generate mockgen --build_flags=-mod=mod -destination=performance/cache/cache.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/cache Cache
