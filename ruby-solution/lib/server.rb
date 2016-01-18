require 'socket'
require_relative 'command'
require_relative 'package_repository'

class Server
  def initialize
    @socket_server = TCPServer.new(8080)
    @package_repository = PackageRepository.new
  end

  def run
    puts "Waiting..."
    while (connection = @socket_server.accept)
      Thread.new(connection) do |conn|
        port, host = connection.peeraddr[1,2]
        client = "#{host}:#{port}"

        puts "#{client} is connected"
        begin
          loop do
            line = connection.readline
            command = Command.new(line)

            result = @package_repository.execute(command)
            response = result ? 1 : 0

            connection.puts(response)
          end
        rescue EOFError => e
          connection.close
          puts "#{client} has disconnected #{e}"
        end
      end
    end
  end
end

Server.new.run

