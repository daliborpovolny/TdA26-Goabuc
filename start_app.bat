@echo off

set loc="%CD%"

echo Running backend...
if exist .\apps\server\cmd\tourbackend (
    cd .\apps\server\cmd\tourbackend
    start cmd /k go run .
)

cd "%loc%"

echo Running frontend...
if exist .\apps\web (
    cd .\apps\web
    start cmd /k npm run dev
)

exit /b