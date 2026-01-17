# Build script for Slip - Cross-platform compilation (PowerShell)
# Usage: .\scripts\build.ps1 [version]

param(
    [string]$Version = "dev"
)

# Configuration
$BuildDir = "build"
$BinaryName = "slip"
$Package = "slip/cmd/slip"

# Colors for output
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Clean build directory
function Clean-BuildDir {
    Write-Info "Cleaning build directory..."
    if (Test-Path $BuildDir) {
        Remove-Item -Recurse -Force $BuildDir
    }
    New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null
    Write-Success "Build directory cleaned"
}

# Get git commit hash (short)
function Get-CommitHash {
    try {
        $hash = git rev-parse --short HEAD 2>$null
        if ($LASTEXITCODE -eq 0) {
            return $hash
        }
    } catch {}
    return "unknown"
}

# Get build timestamp
function Get-Timestamp {
    return (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd_HH:mm:ss_UTC")
}

# Build for a specific platform
function Build-Platform {
    param(
        [string]$Os,
        [string]$Arch
    )

    $OutputName = $BinaryName
    if ($Os -eq "windows") {
        $OutputName = "${BinaryName}.exe"
    }

    $OutputPath = Join-Path $BuildDir "${BinaryName}-${Os}-${Arch}"
    if ($Os -eq "windows") {
        $OutputPath = "${OutputPath}.exe"
    }

    Write-Info "Building for ${Os}/${Arch}..."

    # Set environment variables
    $env:GOOS = $Os
    $env:GOARCH = $Arch

    # Build with version info and optimizations
    $commit = Get-CommitHash
    $buildTime = Get-Timestamp
    $ldflags = "-s -w -X main.Version=$Version -X main.Commit=$commit -X main.BuildTime=$buildTime"

    & go build -ldflags=$ldflags -o $OutputPath $Package

    if ($LASTEXITCODE -eq 0) {
        $size = (Get-Item $OutputPath).Length / 1MB
        $sizeStr = "{0:N2} MB" -f $size
        Write-Success "Built ${Os}/${Arch} (${sizeStr})"

        # Create SHA256 checksum
        $hash = (Get-FileHash -Algorithm SHA256 $OutputPath).Hash
        $checksumFile = "${OutputPath}.sha256"
        $fileName = Split-Path $OutputPath -Leaf
        "${hash}  ${fileName}" | Set-Content $checksumFile
    } else {
        Write-Error-Custom "Failed to build ${Os}/${Arch}"
        return $false
    }

    return $true
}

# Build all platforms
function Build-All {
    Write-Info "Starting cross-platform build for version: $Version"
    Write-Host ""

    $platforms = @(
        @{Os="linux"; Arch="amd64"},
        @{Os="linux"; Arch="arm64"},
        @{Os="darwin"; Arch="amd64"},
        @{Os="darwin"; Arch="arm64"},
        @{Os="windows"; Arch="amd64"}
    )

    $success = $true
    foreach ($platform in $platforms) {
        $result = Build-Platform -Os $platform.Os -Arch $platform.Arch
        if (-not $result) {
            $success = $false
        }
    }

    Write-Host ""
    if ($success) {
        Write-Success "All builds completed!"
    } else {
        Write-Warn "Some builds failed"
    }

    return $success
}

# Create release archives
function Create-Archives {
    Write-Info "Creating release archives..."

    Push-Location $BuildDir

    Get-ChildItem -Filter "slip-*" | Where-Object { -not $_.Name.EndsWith(".sha256") } | ForEach-Object {
        $file = $_.Name

        if ($file -match "\.exe$") {
            # Windows - create ZIP
            $archiveName = $file -replace "\.exe$", ".zip"
            $checksumFile = "${file}.sha256"

            if (Test-Path $checksumFile) {
                Compress-Archive -Path $file, $checksumFile -DestinationPath $archiveName -Force
            } else {
                Compress-Archive -Path $file -DestinationPath $archiveName -Force
            }
        } else {
            # Unix - create tar.gz
            $archiveName = "${file}.tar.gz"
            $checksumFile = "${file}.sha256"

            # Use tar if available (Windows 10+)
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                if (Test-Path $checksumFile) {
                    & tar -czf $archiveName $file $checksumFile 2>$null
                } else {
                    & tar -czf $archiveName $file 2>$null
                }
            }
        }

        if (Test-Path $archiveName) {
            Write-Success "Created $archiveName"
        }
    }

    Pop-Location
}

# Show build summary
function Show-Summary {
    Write-Host ""
    Write-Info "Build Summary:"
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Host "Version:    $Version"
    Write-Host "Commit:     $(Get-CommitHash)"
    Write-Host "Build Time: $(Get-Timestamp)"
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Host ""
    Get-ChildItem $BuildDir | Format-Table Name, Length, LastWriteTime
    Write-Host ""
    Write-Success "Build artifacts are in: $BuildDir\"
}

# Main execution
function Main {
    Write-Info "Slip Build Script v1.0"
    Write-Host ""

    # Check if Go is installed
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Error-Custom "Go is not installed or not in PATH"
        exit 1
    }

    $goVersion = & go version
    Write-Info "Using: $goVersion"
    Write-Host ""

    Clean-BuildDir
    $success = Build-All
    Create-Archives
    Show-Summary

    if ($success) {
        exit 0
    } else {
        exit 1
    }
}

# Run main function
Main
