require 'socket'
require_relative 'command'
require_relative 'package_repository'

class Server
  def initialize
    puts "Creating the server"
  end

  def run
    @socket_server = TCPServer.new(8080)
    @package_repository = PackageRepository.new

    puts "Waiting..."
    while (connection = @socket_server.accept)
      Thread.new(connection) do |conn|
        port, host = conn.peeraddr[1,2]
        client = "#{host}:#{port}"

        puts "#{client} is connected"
        begin
          loop do
            line = conn.readline
            command = Command.new(line)

            result = @package_repository.execute(command)

            response = result ? 1 : 0

            conn.puts(response)
          end
        rescue EOFError => e
          conn.close
          puts "#{client} has disconnected #{e}"
        end
      end
    end
  end
end

Server.new.run

