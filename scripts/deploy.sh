#!/bin/bash

# Deploy script for application

# Set environment variables
ENV=${1:-development}  # Default to development if no argument provided
APP_NAME="my-app"
DOCKER_REPO="my-docker-repo"

# Function to build and push Docker image
build_and_push() {
    echo "Building Docker image for $ENV environment..."
    docker build -t $DOCKER_REPO/$APP_NAME:$ENV .
    echo "Pushing Docker image to repository..."
    docker push $DOCKER_REPO/$APP_NAME:$ENV
}

# Function to deploy to cloud service (example using AWS ECS)
deploy_to_cloud() {
    echo "Deploying to $ENV environment..."
    aws ecs update-service --cluster $ENV-cluster --service $APP_NAME-service --force-new-deployment
}

# Main deployment process
main() {
    case $ENV in
        development|staging|production)
            echo "Deploying to $ENV environment"
            build_and_push
            deploy_to_cloud
            echo "Deployment to $ENV completed successfully"
            ;;
        *)
            echo "Invalid environment. Use development, staging, or production."
            exit 1
            ;;
    esac
}

# Run the main function
main
