node ./src/index.js
$urls = Get-Content -Path 'urls.txt'

foreach ($url in $urls) {
    if ($url.Trim() -ne '') {
        Write-Host "正在下载: $url"
        $folderPath = $url -replace '^https?://[^/]+/', '' -replace '/[^/]+$', '' -replace '%20', ' '
        New-Item -ItemType Directory -Force -Path $folderPath
        $fileName = $url -replace '^.+/', '' -replace '%20', ' '
        $filePath = Join-Path $folderPath $fileName
        curl -o $filePath $url
    }
}

Remove-Item -Path 'urls.txt'

Read-Host -Prompt "下载完成。按任意键退出"
