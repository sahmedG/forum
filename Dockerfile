# Use a specific Go version as the base image
FROM golang:1.21.3

# Set the working directory inside the image
WORKDIR /app

# Copy the source code and other files into the image
COPY . .

# Download the dependencies
RUN go mod download

# Build the application and create an executable binary
RUN go build -o main ./cmd/main.go

# Expose port 8080 for the application
EXPOSE 8080

# Run the application when the container starts
CMD ["./main"]
LABEL "version"="1.0"
LABEL "project name"="Re4um"
LABEL "description"="This project is all about building a user-friendly web forum where users can communicate through sharing posts, comments and interact by liking and disliking. We use cool technologies like SQLite for storing data, Docker to keep things organized, and Go to make it all work seamlessly. while also using essential concepts such as session management, encryption, and database manipulation. By adhering to best practices and incorporating essential packages like bcrypt and UUID."
LABEL "Author"="emahfood, amali, malsamma, sahmed, akhaled"
