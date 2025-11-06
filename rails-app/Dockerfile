FROM ruby:3.2-alpine

# Install dependencies
RUN apk add --no-cache \
    build-base \
    postgresql-dev \
    nodejs \
    npm \
    yarn \
    tzdata \
    git \
    curl

# Set working directory
WORKDIR /app

# Copy Gemfile and install gems
COPY Gemfile Gemfile.lock ./
RUN bundle config set --local without 'development test' && \
    bundle install --jobs 4 --retry 3

# Copy package.json and install node modules
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile --production

# Copy application code
COPY . .

# Precompile assets (will be done at runtime but prepare environment)
RUN bundle exec rails assets:precompile RAILS_ENV=production SECRET_KEY_BASE=dummy || true

# Create necessary directories
RUN mkdir -p tmp/pids tmp/sockets log public/uploads

# Expose port
EXPOSE 3000

# Set environment
ENV RAILS_ENV=production \
    RAILS_LOG_TO_STDOUT=true \
    RAILS_SERVE_STATIC_FILES=true

# Start server
CMD ["bundle", "exec", "puma", "-C", "config/puma.rb"]
