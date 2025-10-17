#!/bin/bash

set -e

echo -e "\nHello, this is a setup agent to create a Go project interactively"
echo -e "\nGithub username:"
read GITHUB_USERNAME

if [ -z "$GITHUB_USERNAME" ]; then
    echo -e "\nUsername cannot be empty!"
    echo "Program is terminated."
    exit 1
fi

echo -e "\nProject name:"
read PROJECT_NAME

if [ -z "$PROJECT_NAME" ]; then
    echo -e "\nProject name cannot be empty!"
    echo "Program is terminated."
    exit 1
fi

if [ -d "$PROJECT_NAME" ]; then
    echo -e "\nProject '$PROJECT_NAME' already exists!"
    echo "Program in terminated to avoid overwriting."
    exit 1
fi

echo -e "\n----------------------------- Creating $PROJECT_NAME -----------------------------\n"

MODULE_NAME="github.com/$GITHUB_USERNAME/$PROJECT_NAME"

mkdir -p "$PROJECT_NAME"/{cmd,internal,pkg}
cd "$PROJECT_NAME"

go mod init "$MODULE_NAME"

cat > cmd/main.go << EOF
package main

import "fmt"

func main() {
    fmt.Println("Hello from $PROJECT_NAME!")
}
EOF

cat > .gitignore << EOF
bin/
vedndor/
*.exe
*.log
EOF

echo -e "\nGo project '$PROJECT_NAME' created successfully!"
echo "Directory structure:"
tree .