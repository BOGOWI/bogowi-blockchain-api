#!/bin/bash

# Istanbul-style coverage reporter for Go
# Mimics the output format of Istanbul/nyc used by JavaScript/Solidity projects

COVERAGE_FILE=${1:-coverage.out}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

# Function to get color based on percentage
get_color() {
    local val=$1
    if (( $(echo "$val >= 80" | bc -l) )); then
        echo -n "$GREEN"
    elif (( $(echo "$val >= 50" | bc -l) )); then
        echo -n "$YELLOW"
    else
        echo -n "$RED"
    fi
}

# Process coverage data and get unique files with their coverage
declare -A file_coverage
declare -A package_files

# Parse coverage data - aggregate by file
go tool cover -func="$COVERAGE_FILE" 2>/dev/null | while read -r line; do
    # Skip empty lines and total
    [[ -z "$line" ]] && continue
    [[ "$line" == total:* ]] && continue
    
    # Parse the line
    filepath=$(echo "$line" | awk '{print $1}')
    coverage=$(echo "$line" | awk '{print $NF}' | sed 's/%//')
    
    # Skip if no valid coverage number
    [[ ! "$coverage" =~ ^[0-9]+(\.[0-9]+)?$ ]] && continue
    
    # Extract clean path and file
    clean_path=$(echo "$filepath" | sed 's/bogowi-blockchain-go\///' | cut -d':' -f1)
    
    # Output in format: filepath|coverage
    echo "$clean_path|$coverage"
done | awk -F'|' '
{
    # Aggregate coverage by file (take the first occurrence which is usually the file summary)
    file = $1
    cov = $2
    if (!(file in seen)) {
        seen[file] = 1
        coverage[file] = cov
        
        # Extract package
        n = split(file, parts, "/")
        pkg = ""
        for (i = 1; i < n; i++) {
            if (i == 1) pkg = parts[i]
            else pkg = pkg "/" parts[i]
        }
        packages[pkg] = packages[pkg] " " parts[n]
        pkg_sum[pkg] += cov
        pkg_count[pkg]++
    }
}
END {
    # Print header
    print "-------------------------|---------|----------|---------|---------|-------------------"
    print "File                     | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s"
    print "-------------------------|---------|----------|---------|---------|-------------------"
    
    # Sort and print by package
    PROCINFO["sorted_in"] = "@ind_str_asc"
    for (pkg in packages) {
        # Calculate package average
        avg = pkg_sum[pkg] / pkg_count[pkg]
        
        # Determine color
        if (avg >= 80) color = "'$GREEN'"
        else if (avg >= 50) color = "'$YELLOW'"
        else color = "'$RED'"
        
        # Print package header
        printf " '$CYAN'%-23s'$NC' | %s%6.1f'$NC' | '$GRAY'   —   '$NC' | '$GRAY'   —   '$NC' | %s%6.1f'$NC' |                   \n", \
            pkg "/", color, avg, color, avg
        
        # Print files in package
        split(packages[pkg], files, " ")
        for (i in files) {
            if (files[i] == "") continue
            filepath = pkg "/" files[i]
            cov = coverage[filepath]
            
            # Determine color for file
            if (cov >= 80) fcolor = "'$GREEN'"
            else if (cov >= 50) fcolor = "'$YELLOW'"
            else fcolor = "'$RED'"
            
            # Truncate long filenames
            filename = files[i]
            if (length(filename) > 22) filename = substr(filename, 1, 19) "..."
            
            printf "  %-22s | %s%6.1f'$NC' | %s%6.1f'$NC' | %s%6.1f'$NC' | %s%6.1f'$NC' |                   \n", \
                filename, fcolor, cov, fcolor, cov, fcolor, cov, fcolor, cov
        }
    }
}'

# Get total coverage
total_coverage=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "^total:" | awk '{print $3}' | sed 's/%//')

# Print footer with total
echo "-------------------------|---------|----------|---------|---------|-------------------"
total_color=$(get_color "$total_coverage")
printf " ${CYAN}All files${NC}               | ${total_color}%6.1f${NC} | ${total_color}%6.1f${NC} | ${total_color}%6.1f${NC} | ${total_color}%6.1f${NC} |                   \n" \
    "$total_coverage" "$total_coverage" "$total_coverage" "$total_coverage"
echo "-------------------------|---------|----------|---------|---------|-------------------"

# Summary statistics
echo ""
if (( $(echo "$total_coverage >= 80" | bc -l) )); then
    echo "✅ Coverage: ${GREEN}${total_coverage}%${NC} PASSED"
elif (( $(echo "$total_coverage >= 60" | bc -l) )); then
    echo "⚠️  Coverage: ${YELLOW}${total_coverage}%${NC} WARNING"
else
    echo "❌ Coverage: ${RED}${total_coverage}%${NC} FAILED"
fi

echo ""
echo "Legend:"
echo "  • Statements, Branches, Functions show same value (Go limitation)"
echo "  • ${GRAY}—${NC} = metric not available separately in Go"
echo "  • Add ${CYAN}-covermode=count${NC} for more detailed branch coverage"