PACKAGES= raylib raygui raymath easings physics rres

GO?= go

all: packages

packages:
	@for pkg in ${PACKAGES}; do \
		echo "Building package github.com/xzebra/raylib-go/$$pkg..."; \
		${GO} build github.com/xzebra/raylib-go/$$pkg || exit 1; \
		${GO} install github.com/xzebra/raylib-go/$$pkg || exit 1; \
	done
