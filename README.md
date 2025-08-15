
# AI CourseBot

Go-based AI chatbot that answers questions about the course catalogs.

It uses:
- ChromaDB (hosted on docker) for semantic course search.
- OpenAI GPT-4o-mini for natural language understanding and reasoning.
- Custom MCP tools for precise queries.
- CSV ingestion to embed course data for retrieval.
- A web server interface for interaction.

## Features

- **Natural Language Q&A** Ask questions like "What CS classes are offered on Mondays?" or "Who teaches Intro to Philosophy?".
- **Vector Search** Finds relevant courses even if your query doesn’t exactly match the CSV wording.
- **Structured Tool Calls**:
  1. courses by instructor
  2. courses by subject
  3. courses by title keyword
  4. courses by meeting day
  5. courses by instruction mode
  6. courses starting before or after a time
  7. calendar invite generation (.ics or details)
- **Web and Terminal Modes** Choose between a browser-based UI or CLI.
- **CSV Loader** Reads and embeds a course catalog into Chroma.
- **Context-Aware Retrieval** Filters search results based on LLM-inferred conditions.
- **Chat History Support** Maintains prior conversation turns so the chatbot can answer follow-up questions with context.

## How It Works

1.  **Database Setup**
    The application connects to a ChromaDB instance hosted inside Docker.
    
2.  **CSV Loading**
    Reads the CSV file with the course catalog and embeds each course description and metadata into the Chroma vector database.
    
3.  **User Query**
   You type a question.
   
5. **LLM Processing** 
    GPT-4o-mini decides whether to answer directly or call a tool.
   Relevant courses are retrieved and tool calls are executed.
    
6.  **Response**
    Returns an AI-generated, formatted answer (or .ics invite if requested).

## Installation

### Prerequisites
- Go 1.21 or higher
- Docker (for running Chroma)
- OpenAI API key stored in .env as OPENAI_API_KEY
- CSV file of course data (default: Fall 2024 Class Schedule 08082024.csv)

### Setup

```bash
# Clone the repo
git clone https://github.com/yourusername/usf-ai-coursebot.git
cd usf-ai-coursebot

# Install dependencies
go mod tidy
```

Create .env:
```
OPENAI_API_KEY=your_api_key_here
```


### Usage

Run:
```go bash
go run $(find . -name "*.go" ! -name "*_test.go")
```
Open browser at:
```
http://localhost:8080
```

### Project Structure

```
├── main.go                # Entry point
├── chatbot.go             # Terminal and web chatbot logic
├── server.go              # Web server routes
├── dbSetup.go             # Loads CSV into Chroma
├── dbInteraction.go       # Retrieves courses from DB
├── llmQuery.go            # Handles LLM calls and tool usage
├── llm_test.go            # Runs tests with sample questions
├── tools.go               # OpenAI tool definitions
├── toolHandlers.go        # Tool execution functions
├── html_files/index.html  # Web UI template
├── Fall 2024 Class Schedule 08082024.csv # Course data
├── go.sum
├── go.mod
└── .env                   # API key storage
```

### Example Queries

- "Which CS classes are on MWF?"
- "Can I learn guitar?"
- "Create a calendar invite for CS110-01."
