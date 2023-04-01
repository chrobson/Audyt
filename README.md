# PokerStars Audit Analyzer

This is a command-line tool for analyzing audit data from PokerStars, a popular online poker platform. The tool retrieves an encrypted CSV file from a URL, unzips it, and parses the data into a list of audit records. It then calculates the total amount of rewards earned from chests in the audit data and prints the result to the console.

## Features:
- Retrieves and decrypts a password-encrypted CSV file from a URL
- Parses the CSV data into a list of audit records
- Calculates the total amount of rewards earned from chests in the audit data
- Prints the result to the console

## Usage:
1. Set the environment variables FILE_URL and FILE_PASSWORD to the URL of the CSV file and its password, respectively.
2. Run the command "go run main.go" in the terminal to start the program.
3. The program will output the total sum of rewards earned from chests in the audit data.
