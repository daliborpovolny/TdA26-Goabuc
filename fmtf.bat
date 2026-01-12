@echo off

set loc="%CD%"

if exist .\apps\web (
    cd .\apps\web
    npm run format
)

cd %loc%

exit /b