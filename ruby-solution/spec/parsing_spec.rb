require 'rspec'
require 'command'

describe Command do
  it 'returns a object if valid command' do
    expectations = {
      Command.new("INSTALL|b|c\n") => {:command => :INSTALL, :package => 'b', :dependencies => ['c']},
      Command.new("INSTALL|b|c") => {:command => :INSTALL, :package => 'b', :dependencies => ['c']},
      Command.new("INSTALL|c|") => {:command => :INSTALL, :package => 'c', :dependencies => []},
      Command.new("UNINSTALL|a|") => {:command => :UNINSTALL, :package => 'a', :dependencies => []},
      Command.new("QUERY|a|") => {:command => :QUERY, :package => 'a', :dependencies => []},
    }

    expectations.each { |actual,expected|
      expect(actual.command).to eq(expected[:command])
      expect(actual.package).to eq(expected[:package])
      expect(actual.dependencies).to eq(expected[:dependencies])
    }
  end

  it 'throws exception if broken msg'
  it 'throws exception if unknown command'
end
