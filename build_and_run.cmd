cd core 
go build .
cd ..\hello 
go build .
move /Y hello.exe .. 
cd ..
.\hello.exe