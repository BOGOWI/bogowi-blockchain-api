#!/bin/bash

# Format Go coverage output to match contracts coverage table style

# Read coverage data from stdin or file
COVERAGE_FILE=${1:-coverage.out}

# Generate function coverage data
go tool cover -func="$COVERAGE_FILE" > /tmp/coverage_raw.txt 2>/dev/null

# Print header
echo "----------------------------|----------|----------|----------|----------|----------------|"
echo "File                        |  % Stmts | % Branch |  % Funcs |  % Lines |Uncovered Lines |"
echo "----------------------------|----------|----------|----------|----------|----------------|"

# Process and group coverage data by package
awk '
BEGIN {
    # Store all data first
    idx = 0
}

# Skip empty lines and non-relevant lines
/^$/ { next }
!/^bogowi-blockchain-go/ && !/^total:/ { next }

# Store lines for processing
{
    lines[idx] = $0
    idx++
}

END {
    current_pkg = ""
    pkg_files[0] = ""
    pkg_coverage[0] = ""
    file_count = 0
    
    # First pass: organize by package
    for (i = 0; i < idx; i++) {
        line = lines[i]
        
        if (line ~ /^total:/) {
            # Handle total line at the end
            split(line, fields, /[ \t]+/)
            gsub(/%/, "", fields[3])
            total_coverage = fields[3]
            continue
        }
        
        # Parse the line
        split(line, fields, /[ \t]+/)
        path = fields[1]
        gsub(/%/, "", fields[3])
        coverage = fields[3] + 0
        
        # Extract package and filename
        n = split(path, parts, "/")
        
        # Build package name (skip bogowi-blockchain-go prefix)
        pkg = ""
        filename = ""
        
        for (j = 2; j <= n; j++) {
            if (parts[j] ~ /:/) {
                # This is the file:line part
                split(parts[j], fileparts, ":")
                filename = fileparts[1]
                gsub(/\.go$/, "", filename)
            } else {
                # Build package path
                if (j == n) {
                    # Last part might be the file if no colon
                    if (parts[j] ~ /\.go/) {
                        filename = parts[j]
                        gsub(/\.go.*/, "", filename)
                    } else {
                        pkg = (pkg == "") ? parts[j] "/" : pkg parts[j] "/"
                    }
                } else {
                    pkg = (pkg == "") ? parts[j] "/" : pkg parts[j] "/"
                }
            }
        }
        
        # Store by package
        if (pkg != current_pkg) {
            if (current_pkg != "" && file_count > 0) {
                # Print previous package
                print_package(current_pkg, pkg_files, pkg_coverage, file_count)
            }
            current_pkg = pkg
            file_count = 0
        }
        
        pkg_files[file_count] = filename
        pkg_coverage[file_count] = coverage
        file_count++
    }
    
    # Print last package
    if (file_count > 0) {
        print_package(current_pkg, pkg_files, pkg_coverage, file_count)
    }
    
    # Print total
    if (total_coverage != "") {
        print "----------------------------|----------|----------|----------|----------|----------------|"
        printf "All files                   |  %6.1f  |    -     |    -     |  %6.1f  |                |\n", 
               total_coverage, total_coverage
        print "----------------------------|----------|----------|----------|----------|----------------|"
    }
}

function print_package(pkg_name, files, coverages, count) {
    # Calculate package average
    sum = 0
    for (i = 0; i < count; i++) {
        sum += coverages[i]
    }
    avg = (count > 0) ? sum / count : 0
    
    # Print package header
    printf " %-27s|  %6.1f  |    -     |    -     |  %6.1f  |                |\n", 
           pkg_name, avg, avg
    
    # Print files in package
    for (i = 0; i < count; i++) {
        if (files[i] != "") {
            printf " %-27s|  %6.1f  |    -     |  %6.1f  |  %6.1f  |                |\n", 
                   "  " files[i], coverages[i], coverages[i], coverages[i]
        }
    }
}
' /tmp/coverage_raw.txt

# Clean up
rm -f /tmp/coverage_raw.txt