# 🌐 go_webpage_analyzer
A lightweight Golang web server that analyzes the contents of web pages by URL.
<br/>
It features simple in-memory caching and processes URLs concurrently with a limit, using a semaphore to control concurrency.


# 🚀 Features
* Analyze web pages by submitting a URL
* In-memory caching to avoid redundant requests
* Concurrent processing of page links with configurable limits
* Server has a gracefull shutdown 

# 🛠 Build the Docker Image
```
docker build -t webpage_analyzer_server .
```

# ▶️ Run the Docker Container
```
docker run -p 8080:8080 webpage_analyzer_server
```

The server will be accessible at: `http://localhost:8080`

# 📌 Notes
* Port Configuration:
Currently, the server port is hardcoded. It can be made configurable via environment variables or command-line flags for more flexibility.

* Caching:
The app uses simple in-memory caching, which is cleared when the container restarts. This could be improved by integrating an external cache like Redis, making the cache persistent and decoupled from the Go server.

* Concurrency Control:
The concurrency limit for processing URLs is fixed. This could be made configurable at runtime to better suit different workloads or system capacities.
