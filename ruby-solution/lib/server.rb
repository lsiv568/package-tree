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
      port, host = connection.peeraddr[1,2]
      client = "#{host}:#{port}"

      puts "#{client} is connected"
      Thread.new(connection) { |c| accept_connection(c) }
    end

  rescue Exception => e
    puts "Error: #{e}"
    exit(1)
  end

  private
  def accept_connection(conn)
    begin
      loop do
        line = conn.readline
        response = process_line(line)
        conn.puts(response)
      end

    rescue Errno::EPIPE, Errno::ECONNRESET, EOFError => e
      puts "#{client} has disconnected (#{e})"
      conn.close
    rescue IOError => e
      puts "Error handling connection: #{e}"
      conn.close
    rescue Exception => e
      puts "Fatal error handling connection: #{e.inspect}"
      conn.close
      exit(1)
    end
  end

  def process_line(line)
    begin
      command = Command.new(line)
      result = @package_repository.execute(command)
      result ? 'OK' : 'FAIL'
    rescue InvalidMessageException => e
      'FAIL'
    end
  end
end

Server.new.run

