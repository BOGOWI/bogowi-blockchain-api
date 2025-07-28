#!/bin/bash

# Generate coverage data
echo "🧪 Running tests with coverage..."
go test -coverprofile=coverage.out ./... 2>/dev/null || echo "Some tests failed, but continuing with coverage report..."

# Get coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Generate JSON from coverage data
echo "📊 Generating coverage report..."
go tool cover -func=coverage.out > coverage.txt

# Create a modern HTML report
cat > coverage-report.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>BOGOWI Coverage Report</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
        body { font-family: 'Inter', sans-serif; }
        .coverage-good { background-color: #10b981; }
        .coverage-ok { background-color: #f59e0b; }
        .coverage-bad { background-color: #ef4444; }
    </style>
</head>
<body class="bg-gray-50">
    <div class="min-h-screen">
        <!-- Header -->
        <header class="bg-white shadow-sm border-b">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                <div class="flex items-center justify-between">
                    <h1 class="text-2xl font-bold text-gray-900">🚀 BOGOWI Coverage Report</h1>
                    <div class="text-sm text-gray-500">Generated on <span id="date"></span></div>
                </div>
            </div>
        </header>

        <!-- Summary -->
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="text-sm font-medium text-gray-500 mb-1">Total Coverage</div>
                    <div class="text-3xl font-bold" id="total-coverage">--</div>
                    <div class="mt-2">
                        <div class="w-full bg-gray-200 rounded-full h-2">
                            <div id="coverage-bar" class="h-2 rounded-full transition-all duration-500"></div>
                        </div>
                    </div>
                </div>
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="text-sm font-medium text-gray-500 mb-1">Packages</div>
                    <div class="text-3xl font-bold" id="package-count">--</div>
                </div>
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="text-sm font-medium text-gray-500 mb-1">Files</div>
                    <div class="text-3xl font-bold" id="file-count">--</div>
                </div>
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="text-sm font-medium text-gray-500 mb-1">Status</div>
                    <div class="text-3xl font-bold" id="status">--</div>
                </div>
            </div>

            <!-- Coverage by Package -->
            <div class="bg-white rounded-lg shadow mb-8">
                <div class="px-6 py-4 border-b">
                    <h2 class="text-lg font-semibold">Coverage by Package</h2>
                </div>
                <div class="p-6">
                    <canvas id="packageChart" height="100"></canvas>
                </div>
            </div>

            <!-- Detailed Coverage -->
            <div class="bg-white rounded-lg shadow">
                <div class="px-6 py-4 border-b">
                    <h2 class="text-lg font-semibold">Detailed Coverage</h2>
                </div>
                <div class="overflow-x-auto">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-gray-50">
                            <tr>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">File</th>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Coverage</th>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Lines</th>
                            </tr>
                        </thead>
                        <tbody id="coverage-tbody" class="bg-white divide-y divide-gray-200">
                            <!-- Rows will be inserted here -->
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Set date
        document.getElementById('date').textContent = new Date().toLocaleString();

        // Parse coverage data (this would come from your coverage.txt)
        const coverageData = `COVERAGE_DATA_PLACEHOLDER`;

        // Process the data
        let files = [];
        let packages = {};
        let totalCoverage = 0;

        // Parse coverage.txt format
        const lines = coverageData.trim().split('\n');
        lines.forEach(line => {
            if (line.includes('.go:') && !line.includes('total:')) {
                const parts = line.split('\t');
                if (parts.length >= 3) {
                    const file = parts[0].trim();
                    const coverage = parseFloat(parts[2].replace('%', ''));
                    files.push({ file, coverage });
                    
                    // Extract package
                    const pkg = file.substring(0, file.lastIndexOf('/'));
                    if (!packages[pkg]) packages[pkg] = [];
                    packages[pkg].push(coverage);
                }
            } else if (line.includes('total:')) {
                const parts = line.split('\t');
                if (parts.length >= 3) {
                    totalCoverage = parseFloat(parts[2].replace('%', ''));
                }
            }
        });

        // Update summary
        document.getElementById('total-coverage').textContent = totalCoverage.toFixed(1) + '%';
        document.getElementById('package-count').textContent = Object.keys(packages).length;
        document.getElementById('file-count').textContent = files.length;
        
        // Update coverage bar
        const coverageBar = document.getElementById('coverage-bar');
        coverageBar.style.width = totalCoverage + '%';
        if (totalCoverage >= 80) {
            coverageBar.classList.add('coverage-good');
            document.getElementById('status').textContent = '✅ Good';
            document.getElementById('status').classList.add('text-green-600');
        } else if (totalCoverage >= 60) {
            coverageBar.classList.add('coverage-ok');
            document.getElementById('status').textContent = '⚠️ OK';
            document.getElementById('status').classList.add('text-yellow-600');
        } else {
            coverageBar.classList.add('coverage-bad');
            document.getElementById('status').textContent = '❌ Low';
            document.getElementById('status').classList.add('text-red-600');
        }

        // Create package chart
        const packageLabels = Object.keys(packages);
        const packageData = packageLabels.map(pkg => {
            const coverages = packages[pkg];
            return coverages.reduce((a, b) => a + b, 0) / coverages.length;
        });

        new Chart(document.getElementById('packageChart'), {
            type: 'bar',
            data: {
                labels: packageLabels.map(label => label.split('/').pop()),
                datasets: [{
                    label: 'Coverage %',
                    data: packageData,
                    backgroundColor: packageData.map(coverage => 
                        coverage >= 80 ? '#10b981' : coverage >= 60 ? '#f59e0b' : '#ef4444'
                    ),
                    borderRadius: 4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            callback: function(value) {
                                return value + '%';
                            }
                        }
                    }
                },
                plugins: {
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return context.parsed.y.toFixed(1) + '%';
                            }
                        }
                    }
                }
            }
        });

        // Fill detailed table
        const tbody = document.getElementById('coverage-tbody');
        files.sort((a, b) => b.coverage - a.coverage).forEach(file => {
            const row = document.createElement('tr');
            const coverageClass = file.coverage >= 80 ? 'text-green-600' : 
                                 file.coverage >= 60 ? 'text-yellow-600' : 'text-red-600';
            
            row.innerHTML = `
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    ${file.file}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm ${coverageClass} font-semibold">
                    ${file.coverage.toFixed(1)}%
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <div class="w-32 bg-gray-200 rounded-full h-2">
                        <div class="h-2 rounded-full ${file.coverage >= 80 ? 'coverage-good' : file.coverage >= 60 ? 'coverage-ok' : 'coverage-bad'}" 
                             style="width: ${file.coverage}%"></div>
                    </div>
                </td>
            `;
            tbody.appendChild(row);
        });
    </script>
</body>
</html>
EOF

# Insert actual coverage data into HTML
COVERAGE_DATA=$(cat coverage.txt)
sed -i.bak "s|COVERAGE_DATA_PLACEHOLDER|${COVERAGE_DATA}|g" coverage-report.html
rm coverage-report.html.bak

echo "✨ Modern coverage report generated: coverage-report.html"
echo "📈 Total coverage: ${COVERAGE}%"

# Open in browser
if [[ "$OSTYPE" == "darwin"* ]]; then
    open coverage-report.html
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    xdg-open coverage-report.html
fi