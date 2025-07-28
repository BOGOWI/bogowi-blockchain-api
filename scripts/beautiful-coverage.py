#!/usr/bin/env python3

import subprocess
import sys
import re
from pathlib import Path

def get_coverage_data():
    """Run tests and get coverage data"""
    print("🧪 Running tests with coverage...")
    subprocess.run(["go", "test", "-coverprofile=coverage.out", "./..."], 
                   stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    
    # Get coverage summary
    result = subprocess.run(["go", "tool", "cover", "-func=coverage.out"], 
                          capture_output=True, text=True)
    return result.stdout

def parse_coverage(coverage_output):
    """Parse coverage output into structured data"""
    lines = coverage_output.strip().split('\n')
    files = []
    total_coverage = 0
    
    for line in lines:
        if line and '\t' in line:
            parts = line.split('\t')
            if len(parts) >= 3:
                file_path = parts[0].strip()
                coverage_str = parts[-1].strip()
                if coverage_str.endswith('%'):
                    coverage = float(coverage_str[:-1])
                    
                    if 'total:' in line:
                        total_coverage = coverage
                    else:
                        # Extract package and file
                        if '/' in file_path:
                            package = '/'.join(file_path.split('/')[:-1])
                            filename = file_path.split('/')[-1].split(':')[0]
                            files.append({
                                'package': package,
                                'file': filename,
                                'path': file_path,
                                'coverage': coverage
                            })
    
    return files, total_coverage

def generate_html(files, total_coverage):
    """Generate beautiful HTML coverage report"""
    
    # Group by package
    packages = {}
    for file in files:
        pkg = file['package']
        if pkg not in packages:
            packages[pkg] = []
        packages[pkg].append(file)
    
    # Calculate package averages
    package_stats = {}
    for pkg, pkg_files in packages.items():
        avg = sum(f['coverage'] for f in pkg_files) / len(pkg_files)
        package_stats[pkg] = {
            'average': avg,
            'files': pkg_files,
            'count': len(pkg_files)
        }
    
    html = f"""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>BOGOWI Coverage Report</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.2.19/tailwind.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
        body {{ font-family: 'Inter', sans-serif; }}
        .gradient-text {{
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }}
        .coverage-bar {{
            background: linear-gradient(90deg, #f59e0b 0%, #10b981 50%, #10b981 100%);
            background-size: 200% 100%;
            background-position: {100 - total_coverage}% 0;
            transition: all 0.5s ease;
        }}
        .hover-lift {{
            transition: transform 0.2s ease;
        }}
        .hover-lift:hover {{
            transform: translateY(-2px);
        }}
    </style>
</head>
<body class="bg-gray-900 text-gray-100">
    <div class="min-h-screen">
        <!-- Header -->
        <header class="bg-gray-800 border-b border-gray-700">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                <div class="flex items-center justify-between">
                    <div>
                        <h1 class="text-3xl font-bold gradient-text">BOGOWI Coverage Report</h1>
                        <p class="text-gray-400 mt-1">Test coverage analysis for your blockchain project</p>
                    </div>
                    <div class="text-right">
                        <p class="text-sm text-gray-400">Generated on</p>
                        <p class="text-gray-300">{'{}'}</p>
                    </div>
                </div>
            </div>
        </header>

        <!-- Summary Cards -->
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
                <!-- Total Coverage Card -->
                <div class="bg-gray-800 rounded-xl p-6 hover-lift">
                    <div class="flex items-center justify-between mb-4">
                        <div class="text-gray-400">
                            <i class="fas fa-chart-line text-2xl"></i>
                        </div>
                        <div class="text-right">
                            <p class="text-3xl font-bold text-{'green' if total_coverage >= 80 else 'yellow' if total_coverage >= 60 else 'red'}-400">
                                {total_coverage:.1f}%
                            </p>
                            <p class="text-sm text-gray-400">Total Coverage</p>
                        </div>
                    </div>
                    <div class="w-full bg-gray-700 rounded-full h-3 overflow-hidden">
                        <div class="coverage-bar h-full rounded-full" style="width: {total_coverage}%"></div>
                    </div>
                </div>

                <!-- Packages Card -->
                <div class="bg-gray-800 rounded-xl p-6 hover-lift">
                    <div class="flex items-center justify-between">
                        <div class="text-gray-400">
                            <i class="fas fa-cube text-2xl"></i>
                        </div>
                        <div class="text-right">
                            <p class="text-3xl font-bold text-blue-400">{len(packages)}</p>
                            <p class="text-sm text-gray-400">Packages</p>
                        </div>
                    </div>
                </div>

                <!-- Files Card -->
                <div class="bg-gray-800 rounded-xl p-6 hover-lift">
                    <div class="flex items-center justify-between">
                        <div class="text-gray-400">
                            <i class="fas fa-file-code text-2xl"></i>
                        </div>
                        <div class="text-right">
                            <p class="text-3xl font-bold text-purple-400">{len(files)}</p>
                            <p class="text-sm text-gray-400">Files</p>
                        </div>
                    </div>
                </div>

                <!-- Status Card -->
                <div class="bg-gray-800 rounded-xl p-6 hover-lift">
                    <div class="flex items-center justify-between">
                        <div class="text-gray-400">
                            <i class="fas fa-{'check-circle' if total_coverage >= 80 else 'exclamation-triangle' if total_coverage >= 60 else 'times-circle'} text-2xl"></i>
                        </div>
                        <div class="text-right">
                            <p class="text-2xl font-bold text-{'green' if total_coverage >= 80 else 'yellow' if total_coverage >= 60 else 'red'}-400">
                                {'Excellent' if total_coverage >= 80 else 'Good' if total_coverage >= 60 else 'Needs Work'}
                            </p>
                            <p class="text-sm text-gray-400">Status</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Package Breakdown -->
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
            <h2 class="text-xl font-semibold mb-6 text-gray-200">Package Coverage</h2>
            <div class="space-y-4">
    """
    
    # Add package cards
    for pkg, stats in sorted(package_stats.items(), key=lambda x: x[1]['average'], reverse=True):
        coverage = stats['average']
        color = 'green' if coverage >= 80 else 'yellow' if coverage >= 60 else 'red'
        
        html += f"""
                <div class="bg-gray-800 rounded-lg overflow-hidden hover-lift">
                    <div class="p-6">
                        <div class="flex items-center justify-between mb-4">
                            <h3 class="text-lg font-medium text-gray-200">
                                <i class="fas fa-folder text-gray-400 mr-2"></i>
                                {pkg.split('/')[-1]}
                            </h3>
                            <div class="flex items-center space-x-4">
                                <span class="text-sm text-gray-400">{stats['count']} files</span>
                                <span class="text-2xl font-bold text-{color}-400">{coverage:.1f}%</span>
                            </div>
                        </div>
                        <div class="w-full bg-gray-700 rounded-full h-2 mb-4">
                            <div class="bg-{color}-400 h-2 rounded-full transition-all duration-500" style="width: {coverage}%"></div>
                        </div>
                        <div class="space-y-2">
        """
        
        # Add files in package
        for file in sorted(stats['files'], key=lambda x: x['coverage'], reverse=True):
            file_color = 'green' if file['coverage'] >= 80 else 'yellow' if file['coverage'] >= 60 else 'red'
            html += f"""
                            <div class="flex items-center justify-between py-2 px-3 bg-gray-900 rounded">
                                <span class="text-sm text-gray-300">
                                    <i class="fas fa-file-code text-gray-500 mr-2"></i>
                                    {file['file']}
                                </span>
                                <div class="flex items-center space-x-3">
                                    <div class="w-24 bg-gray-700 rounded-full h-1.5">
                                        <div class="bg-{file_color}-400 h-1.5 rounded-full" style="width: {file['coverage']}%"></div>
                                    </div>
                                    <span class="text-sm font-medium text-{file_color}-400 w-12 text-right">{file['coverage']:.1f}%</span>
                                </div>
                            </div>
            """
        
        html += """
                        </div>
                    </div>
                </div>
        """
    
    html += """
            </div>
        </div>
    </div>

    <script>
        // Add date
        document.querySelector('.text-gray-300').textContent = new Date().toLocaleString();
        
        // Animate coverage bars on load
        setTimeout(() => {
            document.querySelectorAll('[style*="width: 0"]').forEach(el => {
                el.style.width = el.getAttribute('data-width') + '%';
            });
        }, 100);
    </script>
</body>
</html>
    """
    
    return html.format(
        import_datetime="from datetime import datetime",
        datetime_now="datetime.now().strftime('%Y-%m-%d %H:%M:%S')"
    )

def main():
    # Get coverage data
    coverage_output = get_coverage_data()
    
    # Parse it
    files, total_coverage = parse_coverage(coverage_output)
    
    # Generate HTML
    html = generate_html(files, total_coverage)
    
    # Write to file
    with open('beautiful-coverage.html', 'w') as f:
        f.write(html)
    
    print(f"✨ Beautiful coverage report generated: beautiful-coverage.html")
    print(f"📈 Total coverage: {total_coverage:.1f}%")
    
    # Open in browser
    import platform
    import os
    
    if platform.system() == 'Darwin':
        os.system('open beautiful-coverage.html')
    elif platform.system() == 'Linux':
        os.system('xdg-open beautiful-coverage.html')

if __name__ == '__main__':
    main()