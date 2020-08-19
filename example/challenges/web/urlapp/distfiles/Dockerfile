FROM ruby:2.7

RUN mkdir /app
WORKDIR app

ADD challenge/Gemfile Gemfile
ADD challenge/Gemfile.lock Gemfile.lock

RUN bundle install

ADD challenge/app.rb app.rb
ADD challenge/index.html index.html
ADD challenge/flag.txt flag.txt

CMD bundle exec ruby app.rb -s WEBrick -e production -p 4567
