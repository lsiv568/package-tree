class Command
  attr_reader :command, :package, :dependencies

  def initialize(msg)
    parts = msg.split("|")

    command = parts[0].upcase.to_sym
    package = parts[1]
    deps = parts[2].to_s.split(",")
    @command = command
    @package = package
    @dependencies = deps
  end
end
