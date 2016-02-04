class InvalidMessageException < Exception
end

class Command
  VALID_COMMANDS = %w(INDEX QUERY REMOVE).map(&:to_sym)

  attr_reader :command, :package, :dependencies

  VALID_MESSAGE = /^[A-Z]+\|[\w\-\+]+\|[\w\-\+,]*$/

  def initialize(msg)
    raise InvalidMessageException, "Invalid message #{msg.chomp}: broken syntax" unless msg.match(VALID_MESSAGE)
    parts = msg.chomp.split("|")

    command = parts[0].upcase.to_sym

    raise InvalidMessageException, "Invalid message #{msg.chomp}: command #{command} must be one of #{VALID_COMMANDS}" unless VALID_COMMANDS.include?(command)

    package = parts[1]
    deps = parts[2].to_s.split(",")
    @command = command
    @package = package
    @dependencies = deps
  end
end
