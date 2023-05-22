function Uninstall(){
    $s = Get-Service | Where-Object { $_.Name -eq 'deliverhalf' }
    if ($s -and $s.Status -eq 'Running') {
        Stop-Service deliverhalf -PassThru
        Write-Host "Stopped deliverhalf service."
    }
    if ($s) {
        sc.exe delete deliverhalf
        Write-Host "Removed deliverhalf service."
    }

    $paths = @(
        "C:\Programdata\deliverhalf"
        "C:\Program Files\Deliverhalf"
        "$env:USERPROFILE\.deliverhalf.yaml"
    )

    Set-Location C:/Windows/Temp
    Foreach ($path in $paths){
        if (Test-Path -Path $path) {
            Remove-Item -Force -Recurse $path
            if (-not (Test-Path -Path $path)) {
                Write-Host "$path has been removed"
            } else {
                Write-Host "could not remove $path"
            }
        }
    }
}

function UpdateConfig() {
    $path = "$programdataDirectory\deliverhalf\.deliverhalf.yaml"

    if (Test-Path -Path $path) {
        $oldConfig = "$programdataDirectory\deliverhalf\.deliverhalf.$(Get-Date -f MM-dd-yyyy_HH_mm_ss).yaml"
        Move-Item $path $oldConfig
    }

    if (Test-Path -Path $env:USERPROFILE\.deliverhalf.yaml) {
        $oldConfig = "$programdataDirectory\deliverhalf\.deliverhalf.$(Get-Date -f MM-dd-yyyy_HH_mm_ss).yaml"

        if (Test-Path -Path $path) {
            Move-Item $path $oldConfig
        }
    }

    & "C:\Program Files\Deliverhalf\deliverhalf.exe" config fetch

    if (Test-Path -Path $env:USERPROFILE\.deliverhalf.yaml) {
        New-Item -Type "directory" -Force -Path $programdataDirectory\deliverhalf | Out-Null
        Move-Item $env:USERPROFILE\.deliverhalf.yaml $programdataDirectory\deliverhalf\.deliverhalf.yaml
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
    & "$targetDirectory\nssm.exe" set deliverhalf AppParameters "client send forever --log-level trace --config '""$programdataDirectory\deliverhalf\.deliverhalf.yaml""'"

    Get-Service deliverhalf
    Start-Service deliverhalf -PassThru
}

function FindConfigs() {
    $configs = @(
        "$programdataDirectory\deliverhalf\.deliverhalf.yaml"
        "$env:USERPROFILE\.deliverhalf.yaml"
    )
    foreach ($config in $configs){
        if (Test-Path -Path $config) {
            Write-Host "exists: $config"
        } else {
            Write-Host "gone: $config"
        }
    }
}

function fixup(){
    $targetDirectory = "C:\Program Files\Deliverhalf"
    $env:PATH += ";$targetDirectory"
    StopService
    FindConfigs
    Set-Location C:\Programdata\deliverhalf
    deliverhalf config fetch --log-level trace
    FindConfigs
    if (Test-Path -Path $env:USERPROFILE\.deliverhalf.yaml) {
        Write-Host "removing $env:USERPROFILE\.deliverhalf.yaml"
        Remove-Item $env:USERPROFILE\.deliverhalf.yaml
    }
    FindConfigs
    deliverhalf client send --log-level trace
    &"$targetDirectory\nssm.exe" set deliverhalf AppParameters "client send forever"
    nssm get deliverhalf AppParameters
    &"$targetDirectory\nssm.exe" set deliverhalf AppParameters client send forever --delay 1m --log-level trace --config '""C:\Users\Administrator\.deliverhalf.yaml""'
    nssm get deliverhalf AppParameters
    Start-Service deliverhalf
    Start-Sleep -Seconds 2
    Get-Content C:\Programdata\deliverhalf\deliverhalf.log -Tail 5
    FindConfigs
    StopService
    deliverhalf config fetch --log-level trace
    FindConfigs
    Start-Service deliverhalf
    Start-Sleep -Seconds 2
    Get-Content C:\Programdata\deliverhalf\deliverhalf.log -Tail 5
    Get-Service deliverhalf
    FindConfigs

    StopService
    UpdateConfig
    deliverhalf config fetch --log-level trace
    FindConfigs
    Start-Service deliverhalf
    Start-Sleep -Seconds 2
    Get-Content C:\Programdata\deliverhalf\deliverhalf.log -Tail 5
    Get-Service deliverhalf
    FindConfigs

    StopService
    $path = "$programdataDirectory\deliverhalf\.deliverhalf.yaml"
    if (Test-Path -Path $path) {
        $oldConfig = "$programdataDirectory\deliverhalf\.deliverhalf.$(Get-Date -f MM-dd-yyyy_HH_mm_ss).yaml"
        Move-Item $path $oldConfig
    }
    deliverhalf config fetch --log-level trace
    FindConfigs
    Start-Service deliverhalf
    Start-Sleep -Seconds 2
    Get-Content C:\Programdata\deliverhalf\deliverhalf.log -Tail 5
}

function Install() {
    StopService
    CreateDirectory $targetDirectory
    CreateDirectory $programdataDirectory

    CheckAndDownloadFile -url $downloadUrl -filePath $zipFilePath -expiryMinutes 60
    ExtractZipFile -zipFilePath $zipFilePath -destinationPath $targetDirectory

    StopService
    Write-Host "Deliverhalf executable extracted to: $exeFilePath"

    CheckAndDownloadFile -url $nssmDownloadUrl -filePath $nssmZipFilePath -expiryMinutes 60
    InstallNssm -nssmZipFilePath $nssmZipFilePath -nssmSourcePath $nssmSourcePath -targetDirectory $targetDirectory

    Get-ChildItem -Recurse -Path $programdataDirectory, $targetDirectory | Select-Object -ExpandProperty Fullname
    UpdateConfig
    InstallService
    fixup
}
