#!/bin/bash

# Simple coverage formatter that matches contracts style

COVERAGE_FILE=${1:-coverage.out}

# Get coverage data and process it
go tool cover -func="$COVERAGE_FILE" 2>/dev/null | \
awk '
BEGIN {
    print "----------------------------|----------|----------|----------|----------|----------------|"
    print "File                        |  % Stmts | % Branch |  % Funcs |  % Lines |Uncovered Lines |"
    print "----------------------------|----------|----------|----------|----------|----------------|"
    
    # Initialize package tracking
    current_pkg = ""
    pkg_total = 0
    pkg_count = 0
    pkg_printed = 0
}

# Process each line
{
    if ($1 == "total:") {
        # Print last package summary if needed
        if (current_pkg != "" && pkg_count > 0 && !pkg_printed) {
            avg = pkg_total / pkg_count
            printf " %-27s|  %6.1f  |    -     |    -     |  %6.1f  |                |\n", current_pkg, avg, avg
        }
        
        # Print total
        gsub(/%/, "", $3)
        total = $3 + 0
        print "----------------------------|----------|----------|----------|----------|----------------|"
        printf "All files                   |  %6.1f  |    -     |    -     |  %6.1f  |                |\n", total, total
        print "----------------------------|----------|----------|----------|----------|----------------|"
        next
    }
    
    # Skip non-project files
    if ($1 !~ /^bogowi-blockchain-go/) next
    
    # Parse the path
    path = $1
    gsub(/%/, "", $3)
    coverage = $3 + 0
    
    # Extract package and file
    split(path, parts, "/")
    
    # Determine package (everything except last part)
    pkg = ""
    for (i = 2; i < length(parts); i++) {
        pkg = pkg parts[i] "/"
    }
    
    # Get the file name (last part, remove :line)
    file = parts[length(parts)]
    gsub(/:.*/, "", file)
    gsub(/\.go$/, "", file)
    
    # Check if package changed
    if (pkg != current_pkg) {
        # Print previous package summary
        if (current_pkg != "" && pkg_count > 0 && !pkg_printed) {
            avg = pkg_total / pkg_count
            printf " %-27s|  %6.1f  |    -     |    -     |  %6.1f  |                |\n", current_pkg, avg, avg
        }
        
        # Reset for new package
        current_pkg = pkg
        pkg_total = 0
        pkg_count = 0
        pkg_printed = 0
        
        # Collect all files for this package first
        pkg_total = coverage
        pkg_count = 1
        
        # Print package header immediately
        printf " %-27s|", pkg
        
        # We will calculate average after seeing all files
        pkg_printed = 1
    } else {
        # Same package, accumulate
        pkg_total += coverage
        pkg_count++
    }
}

END {
    # Make sure last package gets its summary
    if (current_pkg != "" && pkg_count > 0 && pkg_printed) {
        avg = pkg_total / pkg_count
        printf "  %6.1f  |    -     |    -     |  %6.1f  |                |\n", avg, avg
    }
}'