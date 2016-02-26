FROM google/golang
RUN go get -u github.com/golang/lint/golint
RUN go get -u github.com/jteeuwen/go-bindata/...

RUN apt-get update -y && apt-get install -y ruby
ENV GEM_HOME /usr/local/bundle
ENV PATH $GEM_HOME/bin:$PATH

ENV BUNDLER_VERSION 1.11.2

ADD ruby-solution/Gemfile Gemfile 

RUN gem install bundler --version "$BUNDLER_VERSION" \
    && bundle config --global path "$GEM_HOME" \
    && bundle config --global bin "$GEM_HOME/bin" \
    && bundle config --global silence_root_warning 1 \
    && bundle install


# don't create ".bundle" in all our apps
ENV BUNDLE_APP_CONFIG $GEM_HOME
