#!/bin/bash

API_KEY="YOUR_API_KEY" # Replace with your actual API key
OUTPUT_FILE="pwned.txt"
TEMP_FILE=$(mktemp)

GREEN='\033[0;32m'
NC='\033[0m'

read -p "Enter the path to the mail list file: " INPUT_FILE

total_lines=$(wc -l < "$INPUT_FILE")
current_line=1

while IFS= read -r email; do
    response=$(curl -s -o /dev/null -w "%{http_code}" -H "Hibp-Api-Key: $API_KEY" "https://haveibeenpwned.com/api/v3/breachedaccount/$email")
    if [ "$response" -eq 200 ]; then
        echo -e "${GREEN}Pwned!! : $email${NC}"
        echo "$email" >> "$TEMP_FILE"
    elif [ "$response" -ne 404 ]; then
        echo "Failed to check pwned status for email $email. Response code: $response"
    fi

    progress=$((current_line * 100 / total_lines))
    echo -ne "Scanning progress: $progress% \r"

    current_line=$((current_line + 1))
    sleep 6
done < "$INPUT_FILE"

if [ -s "$TEMP_FILE" ]; then
    cat "$TEMP_FILE" >> "$OUTPUT_FILE"
    echo "Pwned email addresses saved to $OUTPUT_FILE"
fi

rm "$TEMP_FILE"
