MODEL_FILE ?= v2x-optimization.mod
DEFAULT_DATA_DIR ?= ./data
OUTPUT_REGEX = "RRH_COUNT|RRH"
THREADS_PARAM = -D threads=0 # use default CPLEX value

.PHONY: optimize

optimize:
ifdef data
ifneq (,$(findstring /,$(data)))
# Provided data is located outside the repo data directory
	@oplrun $(THREADS_PARAM) $(MODEL_FILE) $(data) | grep -E $(OUTPUT_REGEX)
else
# Provided data is located in the repo data directory
	@oplrun $(THREADS_PARAM) $(MODEL_FILE) $(DEFAULT_DATA_DIR)/$(data) | grep -E $(OUTPUT_REGEX)
endif
else
	@echo "Data file not defined! Required 'data = [ filepath | filename ]'"
endif
