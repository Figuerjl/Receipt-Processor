Receipt Processor
Receipt Processor is a GoLang application that processes receipts and calculates points based on certain rules. 
This application provides an API endpoint for submitting receipts and retrieving the awarded points.

How to Use
To run the Receipt Processor application, follow the steps below:

Clone the repository:

git clone https://github.com/Figuerjl/Receipt-Processor

Change into the project directory:

cd Receipt-Processor

Build the Docker image:

dcker build -t receipt-processor .
Run the Docker container:

docker run -p 8080:8080 receipt-processor
The Receipt Processor application will be running and listening on http://localhost:8080.

API Endpoints
The Receipt Processor application provides the following API endpoints:

POST /receipts/process: Submits a receipt for processing. The receipt data should be sent in the request body as JSON. The response will include the ID of the processed receipt and the awarded points.

GET /receipts/{id}/points: Retrieves the awarded points for a processed receipt based on its ID.

Pl