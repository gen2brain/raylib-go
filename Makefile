PACKAGES= raylib raygui easings physics rres

GO?= go

all: packages

packages:
	@for pkg in ${PACKAGES}; do \
		echo "Building package github.com/gen2brain/raylib-go/$$pkg..."; \
		${GO} build github.com/gen2brain/raylib-go/$$pkg || exit 1; \
	done
