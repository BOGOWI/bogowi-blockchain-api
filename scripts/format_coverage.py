#!/usr/bin/env python3

import subprocess
import sys
from collections import defaultdict

def get_coverage_data(coverage_file):
    """Run go tool cover and get coverage data"""
    try:
        result = subprocess.run(
            ['go', 'tool', 'cover', '-func=' + coverage_file],
            capture_output=True,
            text=True
        )
        return result.stdout
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

def parse_coverage(coverage_output):
    """Parse coverage output into structured data"""
    packages = defaultdict(list)
    total_coverage = 0
    
    for line in coverage_output.strip().split('\n'):
        if not line or line.startswith('\t'):
            continue
            
        parts = line.split('\t')
        if len(parts) < 3:
            continue
            
        path = parts[0]
        coverage_str = parts[-1].strip()
        
        # Handle total line
        if path == 'total:':
            total_coverage = float(coverage_str.replace('%', ''))
            continue
            
        # Skip non-project files
        if not path.startswith('bogowi-blockchain-go'):
            continue
            
        # Parse coverage percentage
        coverage = float(coverage_str.replace('%', ''))
        
        # Extract package and file
        path_parts = path.split('/')
        
        # Skip the module name
        path_parts = path_parts[1:]
        
        # Extract file and line info
        if ':' in path_parts[-1]:
            file_part = path_parts[-1].split(':')[0]
            package_path = '/'.join(path_parts[:-1]) + '/'
        else:
            file_part = path_parts[-1] if path_parts else ''
            package_path = '/'.join(path_parts[:-1]) + '/' if len(path_parts) > 1 else ''
        
        # Clean up file name
        file_name = file_part.replace('.go', '')
        
        if package_path and file_name:
            packages[package_path].append((file_name, coverage))
    
    return packages, total_coverage

def print_coverage_table(packages, total_coverage):
    """Print coverage in table format like contracts"""
    
    # Header
    print("----------------------------|----------|----------|----------|----------|----------------|")
    print("File                        |  % Stmts | % Branch |  % Funcs |  % Lines |Uncovered Lines |")
    print("----------------------------|----------|----------|----------|----------|----------------|")
    
    # Sort packages for consistent output
    sorted_packages = sorted(packages.keys())
    
    for package in sorted_packages:
        files = packages[package]
        if not files:
            continue
            
        # Calculate package average
        pkg_avg = sum(f[1] for f in files) / len(files) if files else 0
        
        # Print package header
        pkg_display = package if len(package) <= 27 else package[:24] + "..."
        print(f" {pkg_display:<27}|  {pkg_avg:6.1f}  |    -     |    -     |  {pkg_avg:6.1f}  |                |")
        
        # Print files in package (deduplicate by file name)
        seen_files = set()
        for file_name, coverage in files:
            if file_name not in seen_files:
                seen_files.add(file_name)
                file_display = file_name if len(file_name) <= 25 else file_name[:22] + "..."
                print(f"   {file_display:<25}|  {coverage:6.1f}  |    -     |  {coverage:6.1f}  |  {coverage:6.1f}  |                |")
    
    # Footer
    print("----------------------------|----------|----------|----------|----------|----------------|")
    print(f"All files                   |  {total_coverage:6.1f}  |    -     |    -     |  {total_coverage:6.1f}  |                |")
    print("----------------------------|----------|----------|----------|----------|----------------|")

def main():
    coverage_file = sys.argv[1] if len(sys.argv) > 1 else 'coverage.out'
    
    # Get coverage data
    coverage_output = get_coverage_data(coverage_file)
    
    # Parse coverage
    packages, total_coverage = parse_coverage(coverage_output)
    
    # Print table
    print_coverage_table(packages, total_coverage)

if __name__ == '__main__':
    main()