#!/bin/bash

# Enhanced coverage reporter with detailed metrics

COVERAGE_FILE=${1:-coverage.out}
COVERAGE_HTML=${2:-coverage.html}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Generate HTML report for detailed analysis
go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML" 2>/dev/null

# Header
echo "═══════════════════════════════════════════════════════════════════════════════════════════════"
echo "                                    COVERAGE ANALYSIS REPORT                                    "
echo "═══════════════════════════════════════════════════════════════════════════════════════════════"
echo ""

# Function to determine color based on coverage percentage
get_color() {
    local percentage=$1
    if (( $(echo "$percentage >= 80" | bc -l) )); then
        echo "$GREEN"
    elif (( $(echo "$percentage >= 60" | bc -l) )); then
        echo "$YELLOW"
    else
        echo "$RED"
    fi
}

# Detailed file-by-file coverage
echo "┌─────────────────────────────────────────────────────────────────────────────────────────────┐"
echo "│ FILE-BY-FILE COVERAGE                                                                          │"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"
printf "│ %-60s │ %10s │ %10s │\n" "File Path" "Coverage" "Status"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"

# Process each file
go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep -v "^total:" | while IFS=$'\t' read -r file func coverage; do
    # Skip empty lines
    [[ -z "$file" ]] && continue
    
    # Clean up the file path
    clean_file=$(echo "$file" | sed 's/bogowi-blockchain-go\///')
    
    # Extract percentage
    percent=$(echo "$coverage" | sed 's/%//')
    
    # Skip function-level lines, only show file summaries
    if [[ "$clean_file" == *".go"* ]] && [[ "$func" == *")"* ]]; then
        continue
    fi
    
    # Determine color and status
    color=$(get_color "$percent")
    if (( $(echo "$percent >= 80" | bc -l) )); then
        status="✓ Good"
    elif (( $(echo "$percent >= 60" | bc -l) )); then
        status="⚠ Fair"
    else
        status="✗ Low"
    fi
    
    # Truncate long paths
    if [ ${#clean_file} -gt 60 ]; then
        clean_file="${clean_file:0:57}..."
    fi
    
    printf "│ %-60s │ ${color}%9s%%${NC} │ ${color}%10s${NC} │\n" "$clean_file" "$percent" "$status"
done

echo "└─────────────────────────────────────────────────────────────────────────────────────────────┘"
echo ""

# Package-level summary
echo "┌─────────────────────────────────────────────────────────────────────────────────────────────┐"
echo "│ PACKAGE SUMMARY                                                                                │"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"
printf "│ %-60s │ %10s │ %10s │\n" "Package" "Avg Coverage" "Files"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"

# Calculate package averages
declare -A pkg_totals
declare -A pkg_counts

go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep -v "^total:" | while IFS=$'\t' read -r file func coverage; do
    [[ -z "$file" ]] && continue
    
    # Extract package name
    pkg=$(dirname "$file" | sed 's/bogowi-blockchain-go\///')
    
    # Skip function-level entries
    if [[ "$func" == *")"* ]]; then
        continue
    fi
    
    # Extract percentage
    percent=$(echo "$coverage" | sed 's/%//')
    
    # Store in associative arrays (this won't persist outside the loop due to subshell)
    echo "$pkg|$percent"
done | awk -F'|' '
{
    pkg = $1
    percent = $2
    if (pkg != "") {
        totals[pkg] += percent
        counts[pkg]++
    }
}
END {
    for (pkg in totals) {
        avg = totals[pkg] / counts[pkg]
        printf "│ %-60s │ %9.1f%% │ %10d │\n", pkg, avg, counts[pkg]
    }
}' | sort

echo "└─────────────────────────────────────────────────────────────────────────────────────────────┘"
echo ""

# Uncovered lines analysis
echo "┌─────────────────────────────────────────────────────────────────────────────────────────────┐"
echo "│ UNCOVERED CODE ANALYSIS                                                                        │"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"

# Find files with low coverage
echo "│ Files with coverage < 50%:                                                                     │"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"

go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep -v "^total:" | while IFS=$'\t' read -r file func coverage; do
    [[ -z "$file" ]] && continue
    [[ "$func" == *")"* ]] && continue
    
    percent=$(echo "$coverage" | sed 's/%//')
    
    if (( $(echo "$percent < 50" | bc -l) )); then
        clean_file=$(echo "$file" | sed 's/bogowi-blockchain-go\///')
        printf "│   ${RED}%-77s %10s%%${NC} │\n" "$clean_file" "$percent"
    fi
done

echo "└─────────────────────────────────────────────────────────────────────────────────────────────┘"
echo ""

# Overall statistics
echo "┌─────────────────────────────────────────────────────────────────────────────────────────────┐"
echo "│ OVERALL STATISTICS                                                                             │"
echo "├─────────────────────────────────────────────────────────────────────────────────────────────┤"

total=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "^total:" | awk '{print $3}' | sed 's/%//')
color=$(get_color "$total")

printf "│ Total Coverage: ${color}%76s%%${NC} │\n" "$total"

# Count files
total_files=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep -v "^total:" | grep -v ")" | wc -l)
covered_files=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep -v "^total:" | grep -v ")" | awk '{gsub(/%/, "", $3); if ($3 > 0) print}' | wc -l)

printf "│ Files Analyzed: %76d │\n" "$total_files"
printf "│ Files with Coverage: %70d │\n" "$covered_files"
printf "│ Files without Coverage: %67d │\n" "$((total_files - covered_files))"

echo "└─────────────────────────────────────────────────────────────────────────────────────────────┘"
echo ""

# Recommendations
echo "📊 Coverage Report: $COVERAGE_HTML"
echo ""

if (( $(echo "$total < 80" | bc -l) )); then
    echo "💡 Recommendations:"
    echo "   • Current coverage is below 80%. Consider adding more tests."
    echo "   • Focus on files with < 50% coverage first."
    echo "   • Use 'go test -coverprofile=coverage.out -v ./...' to generate detailed coverage."
    echo ""
fi