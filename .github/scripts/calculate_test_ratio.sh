#!/bin/bash

# Check if cloc is installed
if ! command -v cloc &> /dev/null; then
    echo "Error: cloc is not installed. Please install it first."
    echo "You can install it with: brew install cloc"
    exit 1
fi

# Function to get source code stats (excluding tests)
get_source_stats() {
    local dir=$1
    cloc --exclude-dir=_test --exclude-ext=test.go "$dir" --json | jq -r '.Go.code'
}

# Function to get test code stats
get_test_stats() {
    local dir=$1
    cloc --include-ext=test.go "$dir" --json | jq -r '.Go.code'
}

# Calculate for keytool
keytool_source=$(get_source_stats "keytool")
keytool_tests=$(get_test_stats "keytool")

# Calculate for server
server_source=$(get_source_stats "server")
server_tests=$(get_test_stats "server")

# Calculate totals
total_source=$((keytool_source + server_source))
total_tests=$((keytool_tests + server_tests))

# Calculate ratio
if [ "$total_source" -eq 0 ]; then
    ratio="N/A (no source code)"
else
    ratio=$(echo "scale=2; $total_tests / $total_source" | bc)
fi

# Define thresholds
MIN_RATIO=0.5  # Minimum acceptable ratio (50% test-to-code)
WARN_RATIO=0.7  # Warning threshold (70% test-to-code)
TARGET_RATIO=1.0  # Target ratio (100% test-to-code)

echo "Test-to-Code Ratio Report"
echo "========================"
echo "This is a standard metric that measures the relative size of test code compared to production code."
echo "A ratio of 1.0 means equal amounts of test and production code."
echo "Thresholds:"
echo "  Minimum: $MIN_RATIO (50% test-to-code)"
echo "  Warning: $WARN_RATIO (70% test-to-code)"
echo "  Target:  $TARGET_RATIO (100% test-to-code)"
echo
echo "Keytool:"
echo "  Source lines: $keytool_source"
echo "  Test lines:   $keytool_tests"
echo "  Ratio:        $(echo "scale=2; $keytool_tests / $keytool_source" | bc)"
echo
echo "Server:"
echo "  Source lines: $server_source"
echo "  Test lines:   $server_tests"
echo "  Ratio:        $(echo "scale=2; $server_tests / $server_source" | bc)"
echo
echo "Total:"
echo "  Source lines: $total_source"
echo "  Test lines:   $total_tests"
echo "  Ratio:        $ratio"

# Set GitHub step output for workflow summary
if [ "$total_source" -ne 0 ]; then
    echo "test_ratio=$ratio" >> $GITHUB_OUTPUT
    echo "test_ratio_keytool=$(echo "scale=2; $keytool_tests / $keytool_source" | bc)" >> $GITHUB_OUTPUT
    echo "test_ratio_server=$(echo "scale=2; $server_tests / $server_source" | bc)" >> $GITHUB_OUTPUT

    # Check against thresholds
    if (( $(echo "$ratio < $MIN_RATIO" | bc -l) )); then
        echo "ratio_status=❌ Below minimum threshold" >> $GITHUB_OUTPUT
        if [ "$1" = "enforce" ]; then
            exit 1
        fi
    elif (( $(echo "$ratio < $WARN_RATIO" | bc -l) )); then
        echo "ratio_status=⚠️ Below warning threshold" >> $GITHUB_OUTPUT
    elif (( $(echo "$ratio < $TARGET_RATIO" | bc -l) )); then
        echo "ratio_status=✅ Below target threshold" >> $GITHUB_OUTPUT
    else
        echo "ratio_status=🎯 Above target threshold" >> $GITHUB_OUTPUT
    fi
fi
