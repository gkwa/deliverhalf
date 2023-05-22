function UpdateConfig(){
    $tempDirectory = [System.IO.Path]::GetTempPath()
    Set-Location $tempDirectory
    & "C:\Program Files\Deliverhalf\deliverhalf.exe" config fetch

    if (Test-Path -Path $tempDirectory\.deliverhalf.yaml) {
        $oldConfig = "$programdataDirectory\.deliverhalf.$(Get-Date -f MM-dd-yyyy_HH_mm_ss).yaml"
        Move-Item $programdataDirectory\.deliverhalf.yaml $oldConfig
        Move-Item $tempDirectory\.deliverhalf.yaml $programdataDirectory\.deliverhalf.yaml
    }

    if (Test-Path -Path $programdataDirectory\.deliverhalf.yaml) {
        $oldConfig = "$programdataDirectory\.deliverhalf.$(Get-Date -f MM-dd-yyyy_HH_mm_ss).yaml"
        Move-Item $programdataDirectory\.deliverhalf.yaml $oldConfig
        New-Item -Path $directoryPath -ItemType Directory | Out-Null
    }
}

function StopService() {
    $s = Get-Service | Where-Object { $_.Name -eq 'deliverhalf' }
    if ($s -and $s.Status -eq 'Running') {
        Stop-Service deliverhalf -PassThru
        Write-Host "Stopped deliverhalf service."
    }
}

function DownloadFile($url, $filePath) {
    Invoke-WebRequest -Uri $url -OutFile $filePath
}

function ExtractZipFile($zipFilePath, $destinationPath) {
    Expand-Archive -Path $zipFilePath -DestinationPath $destinationPath -Force
}

function CopyFile($sourcePath, $destinationPath) {
    Copy-Item -Path $sourcePath -Destination $destinationPath -Force
}

function CreateDirectory($directoryPath) {
    if (-not (Test-Path -Path $directoryPath -PathType Container)) {
        New-Item -Path $directoryPath -ItemType Directory | Out-Null
    }
}

function InstallNssm($nssmZipFilePath, $nssmSourcePath, $targetDirectory) {
    Expand-Archive -Path $nssmZipFilePath -DestinationPath $tempDirectory -Force
    CopyFile -sourcePath $nssmSourcePath -destinationPath $targetDirectory
}

function CheckAndDownloadFile($url, $filePath, $expiryMinutes) {
    if (-not (Test-Path -Path $filePath) -or (Get-Item -Path $filePath).LastWriteTime -lt (Get-Date).AddMinutes(-$expiryMinutes)) {
        Write-Host "Downloading file: $url"
        DownloadFile -url $url -filePath $filePath
    } else {
        Write-Host "File already up to date: $filePath"
    }
}

$downloadUrl = "https://github.com/taylormonacelli/deliverhalf/releases/latest/download/deliverhalf_Windows_x86_64.zip"
$nssmDownloadUrl = "https://nssm.cc/release/nssm-2.24.zip"

$targetDirectory = "C:\Program Files\Deliverhalf"
$env:PATH += ";$targetDirectory"

$tempDirectory = [System.IO.Path]::GetTempPath()
$zipFilePath = Join-Path -Path $tempDirectory -ChildPath "deliverhalf_Windows_x86_64.zip"
$exeFilePath = Join-Path -Path $targetDirectory -ChildPath "deliverhalf.exe"
$nssmFilePath = Join-Path -Path $targetDirectory -ChildPath "nssm.exe"
$nssmZipFilePath = Join-Path -Path $tempDirectory -ChildPath "nssm-2.24.zip"
$nssmSourcePath = Join-Path -Path $tempDirectory -ChildPath "nssm-2.24\win64\nssm.exe"
$programdataDirectory = "C:\Programdata\deliverhalf"

$global:ProgressPreference = "SilentlyContinue"

function InstallService() {
    $s = Get-Service | Where-Object { $_.Name -eq 'deliverhalf' }
    if (-not $s) {
        & "$targetDirectory\nssm.exe" install deliverhalf "$exeFilePath"
    }
    & "$targetDirectory\nssm.exe" set deliverhalf Start SERVICE_AUTO_START
    & "$targetDirectory\nssm.exe" set deliverhalf DisplayName deliverhalf
    & "$targetDirectory\nssm.exe" set deliverhalf AppDirectory $programdataDirectory
    & "$targetDirectory\nssm.exe" set deliverhalf AppParameters "client send forever --config '""$programdataDirectory\.deliverhalf.yaml""'"

    Get-Service deliverhalf
    Start-Service deliverhalf -PassThru
}

function Main() {
    StopService
    CreateDirectory $targetDirectory
    CreateDirectory $programdataDirectory

    CheckAndDownloadFile -url $downloadUrl -filePath $zipFilePath -expiryMinutes 60
    ExtractZipFile -zipFilePath $zipFilePath -destinationPath $targetDirectory

    Write-Host "Deliverhalf executable extracted to: $exeFilePath"

    CheckAndDownloadFile -url $nssmDownloadUrl -filePath $nssmZipFilePath -expiryMinutes 60
    InstallNssm -nssmZipFilePath $nssmZipFilePath -nssmSourcePath $nssmSourcePath -targetDirectory $targetDirectory

    Get-ChildItem -Recurse -Path $programdataDirectory, $targetDirectory | Select-Object -ExpandProperty Fullname
    UpdateConfig
    InstallService
}

Main
