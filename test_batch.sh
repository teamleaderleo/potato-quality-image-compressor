#!/bin/bash
# Test batch processing functionality

# Configuration
SERVICE_URL="http://localhost:8080"
TEST_IMAGE="test.png"  # Original test image
NUM_IMAGES=5          # Number of test images to create for batch testing

# Terminal colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'           # No Color

echo -e "${BLUE}Image Compression Service - Batch Processing Test${NC}"
echo

# Check if test image exists
if [ ! -f "$TEST_IMAGE" ]; then
    echo -e "${RED}Error: Test image '$TEST_IMAGE' not found!${NC}"
    echo "Please place a PNG image named 'test.png' in the current directory."
    exit 1
fi

# Create duplicate test images with different names
echo "Creating $NUM_IMAGES test images for batch processing..."
TOTAL_SIZE=0
for i in $(seq 1 $NUM_IMAGES); do
    cp "$TEST_IMAGE" "batch_test_${i}.png"
    
    FILE_SIZE=$(stat -c%s "batch_test_${i}.png" 2>/dev/null || stat -f%z "batch_test_${i}.png")
    TOTAL_SIZE=$((TOTAL_SIZE + FILE_SIZE))
done

echo "Total size of all images: $TOTAL_SIZE bytes"
echo

# Test individual processing
echo -e "${BLUE}First testing sequential processing of $NUM_IMAGES images...${NC}"
SEQUENTIAL_START=$(date +%s.%N)

for i in $(seq 1 $NUM_IMAGES); do
    echo "Processing image $i individually..."
    curl -s -X POST -F "image=@batch_test_${i}.png" -F "quality=80" -F "format=webp" \
         "$SERVICE_URL/compress" -o "sequential_${i}.webp"
done

SEQUENTIAL_END=$(date +%s.%N)
SEQUENTIAL_TIME=$(echo "$SEQUENTIAL_END - $SEQUENTIAL_START" | bc)

echo "Sequential processing complete in $SEQUENTIAL_TIME seconds"
echo

# Test batch processing
echo -e "${BLUE}Now testing batch processing of the same $NUM_IMAGES images...${NC}"
BATCH_START=$(date +%s.%N)

# Build curl args as an array (safe and quote-proof)
CURL_ARGS=(-s -X POST)
for i in $(seq 1 $NUM_IMAGES); do
    CURL_ARGS+=(-F "images=@batch_test_${i}.png")
done
CURL_ARGS+=(-F "quality=80" -F "format=webp" "$SERVICE_URL/batch-compress" -o batch_result.zip)

# Execute curl
echo "Sending batch request..."
curl "${CURL_ARGS[@]}"

echo -e "\n--- batch_result.zip (first 200 bytes) ---"
head -c 200 batch_result.zip
echo -e "\n--- end ---"

BATCH_END=$(date +%s.%N)
BATCH_TIME=$(echo "$BATCH_END - $BATCH_START" | bc)

# Check if batch processing was successful
if [ -f "batch_result.zip" ] && [ -s "batch_result.zip" ]; then
    ZIP_SIZE=$(stat -c%s "batch_result.zip" 2>/dev/null || stat -f%z "batch_result.zip")
    echo -e "${GREEN}Batch processing complete in $BATCH_TIME seconds${NC}"
    echo "Result ZIP size: $ZIP_SIZE bytes"
    
    # Calculate speedup
    SPEEDUP=$(echo "scale=2; $SEQUENTIAL_TIME / $BATCH_TIME" | bc)
    echo -e "${GREEN}Batch processing was ${SPEEDUP}x faster than sequential processing${NC}"
    
    # Try to unzip the result to verify contents
    echo "Extracting ZIP file to verify contents..."
    mkdir -p batch_output
    unzip -q batch_result.zip -d batch_output
    
    # Count files in the extracted directory
    FILE_COUNT=$(ls batch_output | wc -l)
    echo "Found $FILE_COUNT files in the ZIP archive"
    
    if [ "$FILE_COUNT" -eq "$NUM_IMAGES" ]; then
        echo -e "${GREEN}All images were successfully processed in batch!${NC}"
    else
        echo -e "${RED}Expected $NUM_IMAGES files but found $FILE_COUNT${NC}"
    fi
else
    echo -e "${RED}Batch processing failed or returned an invalid ZIP${NC}"
fi

echo
echo -e "${BLUE}Results Summary:${NC}"
echo "Sequential processing time: $SEQUENTIAL_TIME seconds"
echo "Batch processing time:      $BATCH_TIME seconds"
if [ -f "batch_result.zip" ] && [ -s "batch_result.zip" ]; then
    echo "Performance improvement:    ${SPEEDUP}x faster with batch processing"
fi

# Clean up
echo
echo "Cleaning up test files..."
rm -f batch_test_*.png sequential_*.webp batch_result.zip
rm -rf batch_output

echo
echo "Test completed!"