PACKAGES= raylib raygui easings physics rres

GO?= go

all: packages

packages:
	@for pkg in ${PACKAGES}; do \
		echo "Building package github.com/icodealot/raylib-go-headless/$$pkg..."; \
		${GO} build github.com/icodealot/raylib-go-headless/$$pkg || exit 1; \
	done
