require 'socket'

puts "Starting..."
server = TCPServer.new(8080)


puts "Waiting..."
while (connection = server.accept)
  Thread.new(connection) do |conn|
    port, host = conn.peeraddr[1,2]
    client = "#{host}:#{port}"
    puts "#{client} is connected"
    begin
      loop do
        line = conn.readline
        puts "#{client} says: #{line}"
        conn.puts(line)
      end
    rescue e => EOFError
      conn.close
      puts "#{client} has disconnected #{e}"
    end
  end
end
