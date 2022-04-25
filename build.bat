@echo off
set QT_DEBUG_CONSOLE=true
set PATH=C:\Users\Public\env_windows_amd64\5.13.0\mingw73_64\bin;C:/Users/Public/env_windows_amd64_Tools/mingw730_64\bin;%PATH%
windres icon.rc -o icon_windows.syso
qtdeploy build desktop
cp icon.ico deploy/windows/idle.ico
cp busy.ico deploy/windows/busy.ico
cp -r PaddleOCRServer/* deploy/windows

