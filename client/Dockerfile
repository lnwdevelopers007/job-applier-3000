# Client Dockerfile
FROM oven/bun:1.3
WORKDIR /app

# Install dependencies
COPY package.json ./
RUN bun install

# Copy everything
COPY . .

# Load .env file
COPY .env.production .env

EXPOSE 5173

CMD ["bun", "run", "dev", "--host"]
