# Demo Order Service

The service integrates with PostgreSQL for data storage and utilizes NATS Streaming for message queuing. 
The service includes caching mechanisms to enhance performance and fault tolerance.

## Features
 1) PostgreSQL Integration: Set up a local PostgreSQL database, create tables, and configure a user for the service.
 2) NATS Streaming Integration: Establish connection and subscription to a NATS Streaming channel to receive order data.
 3) Data Storage: Store received data in the PostgreSQL database.
 4) Caching Mechanism: Implement in-memory caching of received data for improved performance. The service restores cache from the database in case of service failure.
 5) HTTP Server: Launch an HTTP server to serve cached data based on order IDs.
 6) User Interface: Develop a basic interface to display order data based on order IDs.