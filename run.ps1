# Copy test file and run binary
Copy-Item "testfiles/cayinput566000.csv" "input/AP566000.csv"
./go-mail.exe
