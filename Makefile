watch_build:
	if command -v entr > /dev/null; then find . -print | grep -i '.*[.]go' | entr -c go run main.go; else go run main.go; echo "\nInstall entr(1) to automatically rebuild documentation when files change. \nSee http://entrproject.org/"; fi
watch_test:
	if command -v entr > /dev/null; then find . -print | grep -i '.*[.]go' | entr -c go test ./...; else go test ./...; echo "\nInstall entr(1) to automatically rebuild documentation when files change. \nSee http://entrproject.org/"; fi
