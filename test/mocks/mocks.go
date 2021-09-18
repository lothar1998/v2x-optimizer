package mocks

//v2x-optimizer
//go:generate mockgen --build_flags=-mod=mod -destination=pkg_encoder_decoder.go -package=mocks github.com/lothar1998/v2x-optimizer/pkg/data EncoderDecoder

//v2x-optimizer-performance
//go:generate mockgen --build_flags=-mod=mod -destination=internal_process.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/executor Process
//go:generate mockgen --build_flags=-mod=mod -destination=internal_executor.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/executor Executor
//go:generate mockgen --build_flags=-mod=mod -destination=internal_identifiable_optimizer.go -package=mocks github.com/lothar1998/v2x-optimizer/internal/performance/optimizer IdentifiableOptimizer
