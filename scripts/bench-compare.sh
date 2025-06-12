#!/bin/bash
# Benchmark Version Comparison Tool
# 用于对比当前代码与指定 Git 版本的 benchmark 性能差异

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BENCH_DIR=".benchmarks"
BENCH_COUNT=${BENCH_COUNT:-5}
BENCH_TIME=${BENCH_TIME:-1s}

# Functions
print_usage() {
    echo "Usage: $0 [OPTIONS] [VERSION]"
    echo ""
    echo "Compare benchmark results between current code and a specific version."
    echo ""
    echo "Arguments:"
    echo "  VERSION    Git tag, branch, or commit to compare against (default: latest tag)"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -l, --list     List available tags for comparison"
    echo "  -c, --count N  Number of benchmark iterations (default: 5)"
    echo "  -s, --save     Save current results as baseline after comparison"
    echo "  -o, --output   Output file for comparison results"
    echo ""
    echo "Examples:"
    echo "  $0                    # Compare with latest tag"
    echo "  $0 v1.0.0             # Compare with specific tag"
    echo "  $0 -l                 # List available tags"
    echo "  $0 -c 10 v1.0.0       # Run 10 iterations, compare with v1.0.0"
    echo ""
}

list_tags() {
    echo -e "${BLUE}Available tags for comparison:${NC}"
    echo ""
    git tag -l --sort=-v:refname | head -20
    echo ""
    echo -e "${YELLOW}Tip: Use 'git tag -l' to see all tags${NC}"
}

get_latest_tag() {
    git describe --tags --abbrev=0 2>/dev/null || echo ""
}

check_benchstat() {
    if ! command -v benchstat &> /dev/null; then
        echo -e "${RED}Error: benchstat not found${NC}"
        echo ""
        echo "Install benchstat with:"
        echo "  go install golang.org/x/perf/cmd/benchstat@latest"
        echo ""
        exit 1
    fi
}

run_benchmark() {
    local output_file=$1
    echo -e "${BLUE}Running benchmarks (count=$BENCH_COUNT)...${NC}"
    go test -bench=. -benchmem -count=$BENCH_COUNT -benchtime=$BENCH_TIME ./test/... 2>&1 | tee "$output_file"
    echo ""
}

# Parse arguments
SAVE_BASELINE=false
OUTPUT_FILE=""
VERSION=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -l|--list)
            list_tags
            exit 0
            ;;
        -c|--count)
            BENCH_COUNT="$2"
            shift 2
            ;;
        -s|--save)
            SAVE_BASELINE=true
            shift
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        *)
            VERSION="$1"
            shift
            ;;
    esac
done

# Check prerequisites
check_benchstat

# Create benchmark directory
mkdir -p "$BENCH_DIR"

# Determine version to compare
if [ -z "$VERSION" ]; then
    VERSION=$(get_latest_tag)
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Error: No tags found. Please specify a version to compare.${NC}"
        exit 1
    fi
    echo -e "${YELLOW}No version specified, using latest tag: $VERSION${NC}"
fi

# Verify version exists
if ! git rev-parse "$VERSION" &>/dev/null; then
    echo -e "${RED}Error: Version '$VERSION' not found${NC}"
    echo ""
    echo "Available tags:"
    git tag -l --sort=-v:refname | head -10
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Benchmark Comparison Tool${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "Comparing: ${BLUE}current${NC} vs ${BLUE}$VERSION${NC}"
echo -e "Iterations: ${BLUE}$BENCH_COUNT${NC}"
echo ""

# Save current state
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "HEAD")
CURRENT_COMMIT=$(git rev-parse HEAD)
STASH_NEEDED=false

# Check for uncommitted changes
if ! git diff --quiet || ! git diff --cached --quiet; then
    echo -e "${YELLOW}Stashing uncommitted changes...${NC}"
    git stash push -m "bench-compare-temp-stash"
    STASH_NEEDED=true
fi

# Cleanup function
cleanup() {
    echo ""
    echo -e "${BLUE}Cleaning up...${NC}"
    git checkout "$CURRENT_BRANCH" 2>/dev/null || git checkout "$CURRENT_COMMIT"
    if [ "$STASH_NEEDED" = true ]; then
        git stash pop 2>/dev/null || true
    fi
}

trap cleanup EXIT

# Run benchmark on old version
echo -e "${GREEN}Step 1/3: Running benchmark on $VERSION${NC}"
echo "----------------------------------------"
git checkout "$VERSION" --quiet
OLD_BENCH="$BENCH_DIR/old-$VERSION.txt"
run_benchmark "$OLD_BENCH"

# Run benchmark on current version
echo -e "${GREEN}Step 2/3: Running benchmark on current code${NC}"
echo "----------------------------------------"
git checkout "$CURRENT_BRANCH" 2>/dev/null || git checkout "$CURRENT_COMMIT"
if [ "$STASH_NEEDED" = true ]; then
    git stash pop --quiet
    STASH_NEEDED=false
fi
NEW_BENCH="$BENCH_DIR/current.txt"
run_benchmark "$NEW_BENCH"

# Compare results
echo -e "${GREEN}Step 3/3: Comparing results${NC}"
echo "----------------------------------------"
echo ""

COMPARISON_OUTPUT=$(benchstat "$OLD_BENCH" "$NEW_BENCH")
echo "$COMPARISON_OUTPUT"

# Save to file if requested
if [ -n "$OUTPUT_FILE" ]; then
    echo "$COMPARISON_OUTPUT" > "$OUTPUT_FILE"
    echo ""
    echo -e "${GREEN}Results saved to: $OUTPUT_FILE${NC}"
fi

# Save as baseline if requested
if [ "$SAVE_BASELINE" = true ]; then
    cp "$NEW_BENCH" "$BENCH_DIR/baseline.txt"
    echo ""
    echo -e "${GREEN}Current results saved as baseline${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Comparison complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Files generated:"
echo "  - $OLD_BENCH (benchmark for $VERSION)"
echo "  - $NEW_BENCH (benchmark for current)"
echo ""
echo -e "${YELLOW}Legend:${NC}"
echo "  - Negative delta (e.g., -10%) = Performance improved"
echo "  - Positive delta (e.g., +10%) = Performance degraded"
echo "  - ~ = No significant change"
