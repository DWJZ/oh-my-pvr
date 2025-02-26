# Stage 1: Build
FROM node:18-alpine AS builder
WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Copy package files and install dependencies using pnpm.
# Make sure your project has a pnpm-lock.yaml for consistent installs.
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# Copy the rest of the code and build the Next.js app.
COPY . .
RUN pnpm build

# Stage 2: Production container
FROM node:18-alpine AS runner

ENV NODE_ENV=production
WORKDIR /app

# Install pnpm in the production image.
RUN npm install -g pnpm

# Copy built files and necessary code from the builder stage.
COPY --from=builder /app/public ./public

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Expose Next.js default port.
EXPOSE 3000

ENV PORT=3000

# server.js is created by next build from the standalone output
# https://nextjs.org/docs/pages/api-reference/config/next-config-js/output
ENV HOSTNAME="0.0.0.0"
CMD ["node", "server.js"]