# Skill Swap: Microservice Application

## Project Overview

Skill Swap is a microservice-based application built in Go that enables users to post skills they can offer or want to learn, and provides matching between users based on complementary skills. This project is designed as a learning exercise for various modern technologies including Go, gRPC, GraphQL, MongoDB, RabbitMQ, and Elasticsearch.

## Architecture

The application is structured as a collection of microservices:

1. **User Service**: Handles user registration, authentication, and profile management
2. **Skill Service**: Manages skill offerings and learning requests
3. **Match Service**: Matches users based on complementary skills
4. **Search Service**: Provides search functionality using Elasticsearch
5. **Notification Service**: Manages notifications for matches and messages
6. **API Gateway**: GraphQL gateway that serves as the entry point for client applications

## Technology Stack

- **Programming Language**: Go
- **Database**: MongoDB (document storage)
- **Inter-service Communication**: gRPC
- **Message Queue**: RabbitMQ (for asynchronous processing)
- **Search Engine**: Elasticsearch
- **API Gateway**: GraphQL
- **Containerization**: Docker & Docker Compose

## Getting Started

### Prerequisites

- Go 1.18 or later
- Docker and Docker Compose
- MongoDB
- RabbitMQ
- Elasticsearch

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/BhuvanMM/skill-swap.git
   cd skill-swap
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start the infrastructure services using Docker Compose:
   ```bash
   docker-compose up -d mongodb rabbitmq elasticsearch
   ```

4. Generate gRPC code from protocol buffers:
   ```bash
   # Install protoc compiler if you haven't already
   # Then generate code from .proto files
   protoc --go_out=. --go-grpc_out=. common/proto/*.proto
   ```

5. Run each service:
   ```bash
   # Start each service in a separate terminal
   go run user-service/main.go
   go run skill-service/main.go
   go run match-service/main.go
   go run search-service/main.go
   go run notification-service/main.go
   go run api-gateway/main.go
   ```

### Environment Variables

Each service can be configured using environment variables:

```bash
# MongoDB connection
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=skillswap

# RabbitMQ connection
RABBITMQ_URI=amqp://guest:guest@localhost:5672/

# Elasticsearch connection
ELASTICSEARCH_URI=http://localhost:9200

# JWT authentication
JWT_SECRET=your-secret-key

# Service ports
GRPC_PORT=50051  # Different for each service
GRAPHQL_PORT=8080  # For API Gateway
```

## Service Details

### User Service

The user service manages user registration, authentication, and profile information. It provides:
- User registration and login
- Profile management
- Authentication via JWT

### Skill Service

The skill service handles the creation and management of skills that users can offer or want to learn:
- Creation of skills (offering or learning)
- Listing skills by user
- Categorizing and describing skills

### Match Service

The match service connects users based on complementary skills:
- Finds potential skill matches
- Manages match status (pending, accepted, rejected)
- Notifies users of new matches

### Search Service

The search service provides advanced search capabilities using Elasticsearch:
- Searching for skills by name, description, or category
- Finding users by skills they offer or want to learn
- Fuzzy matching and relevance sorting

### Notification Service

The notification service manages notifications for matches and messages:
- Sends notifications when a match is found
- Alerts users to new messages or system events
- Manages notification preferences

### API Gateway

The GraphQL API gateway serves as the entry point for client applications:
- Combines data from multiple services
- Provides a unified API for clients
- Handles authentication and authorization

## Data Models

### User
- ID
- Name
- Email
- Password (hashed)
- Skills offering (references to skills)
- Skills learning (references to skills)
- Created/Updated timestamps

### Skill
- ID
- User ID (reference to user)
- Name
- Description
- Category
- Type (offering or learning)
- Proficiency level (1-5)
- Created/Updated timestamps

### Match
- ID
- User ID A (reference to user)
- Skill ID A (reference to skill)
- User ID B (reference to user)
- Skill ID B (reference to skill)
- Status (pending, accepted, rejected)
- Created/Updated timestamps

## API Examples

### GraphQL API (API Gateway)

```graphql
# Query to get a user with their skills
query {
  user(id: "user_id") {
    id
    name
    email
    skillsOffering {
      id
      name
      description
      category
      proficiencyLevel
    }
    skillsLearning {
      id
      name
      description
      category
      proficiencyLevel
    }
  }
}

# Mutation to create a new skill
mutation {
  createSkill(
    userId: "user_id"
    name: "JavaScript Programming"
    description: "Web development with JavaScript"
    category: "Programming"
    type: "offering"
    proficiencyLevel: 4
  ) {
    id
    name
    description
  }
}
```

## Development Notes

- This project uses a clean architecture pattern with separated concerns
- Each service is independently deployable and scalable
- Common code is shared via the `common` package
- The project emphasizes simplicity and learning rather than production-readiness

## Future Enhancements

- Add websocket support for real-time notifications
- Implement service discovery
- Add metrics and monitoring
- Implement caching for performance
- Add rate limiting and security features
- Create mobile and web client applications

## License

This project is licensed under the MIT License - see the LICENSE file for details.
