# Blog-Aggregator

Users can login and add RSS feeds which the aggregator would use to collect the posts from the feeds and display them

## Instructions
1. Open your terminal and go to the src folder of the backend folder of the blog-aggregator project
2. Create a .env file and fill it out with the following lines:
PORT="8080"

JWT_SECRET = ""

DB_URL=""

3. Run the command 'openssl rand -base64 64' in your terminal and copy and paste the output into the quotes for the JWT_SECRET line of the .env file

4. Run 'go build && ./src'