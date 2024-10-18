# This script is used to push the Docker image to Docker Hub
# The first argument is the username
# The second argument is the repository name
# The third argument is the tag

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Usage: $0 <dockerhub-username> <repository-name> <tag>"
    echo "Example: $0 intothevoid likho latest"
    exit 1
fi

docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t $1/$2:$3 --push .
