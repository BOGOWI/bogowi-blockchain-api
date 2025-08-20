#!/bin/bash

# Simple coverage table formatter

COVERAGE_FILE=${1:-coverage.out}

# Header
echo "----------------------------|----------|----------|----------|----------|----------------|"
echo "File                        |  % Stmts | % Branch |  % Funcs |  % Lines |Uncovered Lines |"
echo "----------------------------|----------|----------|----------|----------|----------------|"

# Process coverage by package
go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep -E "^bogowi-blockchain-go|^total:" | \
sed 's/bogowi-blockchain-go\///' | \
awk '
{
    if ($1 == "total:") {
        gsub(/%/, "", $3)
        total = $3
    } else {
        # Store for later processing
        path = $1
        gsub(/%/, "", $3) 
        coverage = $3
        
        # Extract directory structure
        n = split(path, parts, "/")
        
        # Group by directory
        if (n >= 2) {
            dir = parts[1] "/" parts[2] "/"
            if (!(dir in dirs)) {
                dirs[dir] = coverage
                counts[dir] = 1
            } else {
                dirs[dir] += coverage
                counts[dir]++
            }
        }
    }
}
END {
    # Print grouped results
    for (dir in dirs) {
        avg = dirs[dir] / counts[dir]
        # Clean display name
        display = dir
        if (length(display) > 27) display = substr(display, 1, 24) "..."
        printf " %-27s|  %6.1f  |    -     |    -     |  %6.1f  |                |\n", display, avg, avg
    }
    
    # Print total
    if (total) {
        print "----------------------------|----------|----------|----------|----------|----------------|"
        printf "All files                   |  %6.1f  |    -     |    -     |  %6.1f  |                |\n", total, total
        print "----------------------------|----------|----------|----------|----------|----------------|"
    }
}'