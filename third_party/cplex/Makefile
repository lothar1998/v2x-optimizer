.PHONY: optimize

optimize:
ifdef data
ifneq (,$(findstring /,$(data)))
# Provided data is located outside the repo data directory
		@oplrun v2x-optimization.mod $(data)
else
# Provided data is located in the repo data directory
		@oplrun v2x-optimization.mod ./data/$(data)
endif
else
	@echo "Data file not defined! Required 'data = [ filepath | filename ]'"
endif
