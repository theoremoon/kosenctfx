require 'sinatra'
require 'uri'
require 'socket'

def connect()
  sock = TCPSocket.open("redis", 6379)

  if not ping(sock) then
    exit
  end

  return sock
end

def query(sock, cmd)
  sock.write(cmd + "\r\n")
end

def recv(sock)
  data = sock.gets
  if data == nil then
    return nil
  elsif data[0] == "+" then
    return data[1..-1].strip
  elsif data[0] == "$" then
    if data == "$-1\r\n" then
      return nil
    end
    return sock.gets.strip
  end

  return nil
end

def ping(sock)
  query(sock, "ping")
  return recv(sock) == "PONG"
end

def set(sock, key, value)
  query(sock, "SET #{key} #{value}")
  return recv(sock) == "OK"
end

def get(sock, key)
  query(sock, "GET #{key}")
  return recv(sock)
end

before do
  sock = connect()
  set(sock, "flag", File.read("flag.txt").strip)
end

get '/' do
  if params.has_key?(:q) then
    q = params[:q]
    if not (q =~ /^[0-9a-f]{16}$/)
      return
    end

    sock = connect()
    url = get(sock, q)
    redirect url
  end

  send_file 'index.html'
end

post '/' do
  if not params.has_key?(:url) then
    return
  end

  url = params[:url]
  if not (url =~ URI.regexp) then
    return
  end

  key = Random.urandom(8).unpack("H*")[0]
  sock = connect()
  set(sock, key, url)

  "#{request.host}:#{request.port}/?q=#{key}"
end

