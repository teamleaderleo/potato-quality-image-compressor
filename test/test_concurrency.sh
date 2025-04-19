#!/bin/bash
# Concurrency test using awk for calculations

# Configuration
SERVICE_URL="http://localhost:8080"
TEST_IMAGE="test.png"  # Using PNG
REQUESTS=16            # Number of concurrent requests
WORKER_COUNT=16         # Should match our service configuration

# Terminal colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'           # No Color

echo -e "${BLUE}Image Compression Service - Concurrency Test${NC}"
echo "Testing with $REQUESTS concurrent requests using $WORKER_COUNT workers"
echo "Using image: $TEST_IMAGE"
echo

# Check if test image exists
if [ ! -f "$TEST_IMAGE" ]; then
    echo -e "${RED}Error: Test image '$TEST_IMAGE' not found!${NC}"
    echo "Please place a PNG image named 'test.png' in the current directory."
    exit 1
fi

# Get original size
ORIG_SIZE=$(stat -c%s "$TEST_IMAGE" 2>/dev/null || stat -f%z "$TEST_IMAGE" 2>/dev/null || ls -l "$TEST_IMAGE" | awk '{print $5}')
echo "Original image size: $ORIG_SIZE bytes ($(awk "BEGIN {printf \"%.2f\", $ORIG_SIZE/1024/1024}") MB)"
echo

echo "Starting $REQUESTS concurrent requests..."
echo "This may take a few seconds..."

# Use date with nanoseconds for more precision
START_TIME=$(date +%s.%N)

# Run concurrent requests
for i in $(seq 1 $REQUESTS); do
    # Use curl in background with unique output file
    curl -s -X POST -F "image=@$TEST_IMAGE" -F "quality=80" -F "format=webp" \
         "$SERVICE_URL/compress" -o "concurrent_${i}.webp" &
done

# Wait for all background processes to complete
wait

# Measure end time with nanoseconds
END_TIME=$(date +%s.%N)

# Calculate elapsed time using awk
ELAPSED=$(awk "BEGIN {printf \"%.2f\", $END_TIME - $START_TIME}")
echo "Total processing time: $ELAPSED seconds"

# Count successful responses
SUCCESS_COUNT=0
TOTAL_COMPRESSED_SIZE=0

for i in $(seq 1 $REQUESTS); do
    OUTPUT_FILE="concurrent_${i}.webp"
    if [ -f "$OUTPUT_FILE" ] && [ -s "$OUTPUT_FILE" ]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        
        # Get size of compressed file
        COMP_SIZE=$(stat -c%s "$OUTPUT_FILE" 2>/dev/null || stat -f%z "$OUTPUT_FILE" 2>/dev/null || ls -l "$OUTPUT_FILE" | awk '{print $5}')
        TOTAL_COMPRESSED_SIZE=$((TOTAL_COMPRESSED_SIZE + COMP_SIZE))
        
        # Calculate ratio using awk
        RATIO=$(awk "BEGIN {printf \"%.2f\", $ORIG_SIZE/$COMP_SIZE}")
        
        echo -e "Request $i: ${GREEN}Success${NC} - Compressed size: $COMP_SIZE bytes ($(awk "BEGIN {printf \"%.2f\", $COMP_SIZE/1024}") KB, ${RATIO}x smaller)"
    else
        echo -e "Request $i: ${RED}Failed${NC}"
    fi
done

# Calculate average compression ratio if any succeeded
if [ $SUCCESS_COUNT -gt 0 ]; then
    AVG_COMP_SIZE=$(awk "BEGIN {printf \"%.2f\", $TOTAL_COMPRESSED_SIZE/$SUCCESS_COUNT}")
    AVG_RATIO=$(awk "BEGIN {printf \"%.2f\", $ORIG_SIZE/($TOTAL_COMPRESSED_SIZE/$SUCCESS_COUNT)}")
    echo
    echo -e "${GREEN}Average compression ratio: ${AVG_RATIO}x${NC}"
    echo -e "Average compressed size: $(awk "BEGIN {printf \"%.2f\", $TOTAL_COMPRESSED_SIZE/$SUCCESS_COUNT/1024}") KB (from original $(awk "BEGIN {printf \"%.2f\", $ORIG_SIZE/1024/1024}") MB)"
fi

# Calculate success rate
SUCCESS_RATE=$(awk "BEGIN {printf \"%.0f\", ($SUCCESS_COUNT/$REQUESTS)*100}")

echo
echo -e "${BLUE}Results:${NC}"
echo "Success rate: ${SUCCESS_RATE}% ($SUCCESS_COUNT/$REQUESTS)"

# Now run a sequential test for comparison
echo
echo -e "${BLUE}Running sequential test for comparison...${NC}"
rm -f concurrent_*.webp  # Clean up previous files

# Use date with nanoseconds for more precision
SEQ_START_TIME=$(date +%s.%N)

# Run sequential requests (one after another)
for i in $(seq 1 $REQUESTS); do
    echo -n "Processing request $i... "
    curl -s -X POST -F "image=@$TEST_IMAGE" -F "quality=80" -F "format=webp" \
         "$SERVICE_URL/compress" -o "sequential_${i}.webp"
    echo "done"
done

# Measure end time with nanoseconds
SEQ_END_TIME=$(date +%s.%N)

# Calculate sequential elapsed time
SEQ_ELAPSED=$(awk "BEGIN {printf \"%.2f\", $SEQ_END_TIME - $SEQ_START_TIME}")
echo "Sequential processing time: $SEQ_ELAPSED seconds"

# Calculate actual speedup
SPEEDUP=$(awk "BEGIN {printf \"%.2f\", $SEQ_ELAPSED/$ELAPSED}")
echo -e "${GREEN}Actual speedup from worker pool: ${SPEEDUP}x${NC}"

echo
echo -e "${BLUE}Worker Pool Analysis:${NC}"
echo "With $WORKER_COUNT workers, we would theoretically expect up to ${WORKER_COUNT}x speedup"
echo "The measured speedup shows how well your concurrency implementation is working"

# Calculate efficiency
EFFICIENCY=$(awk "BEGIN {printf \"%.1f\", ($SPEEDUP/$WORKER_COUNT)*100}")
echo -e "Worker pool efficiency: ${EFFICIENCY}% of theoretical maximum"

# Clean up
echo
echo "Cleaning up test files..."
rm -f concurrent_*.webp sequential_*.webp

echo
echo "Test completed!"