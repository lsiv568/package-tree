require 'rspec'
require 'command'

describe Command do
  it 'returns a object if valid command' do
    expectations = {
      Command.new("INDEX|b|c\n") => {:command => :INDEX, :package => 'b', :dependencies => ['c']},
      Command.new("INDEX|b|c") => {:command => :INDEX, :package => 'b', :dependencies => ['c']},
      Command.new("INDEX|c|") => {:command => :INDEX, :package => 'c', :dependencies => []},
      Command.new("REMOVE|a|") => {:command => :REMOVE, :package => 'a', :dependencies => []},
      Command.new("QUERY|a|") => {:command => :QUERY, :package => 'a', :dependencies => []},
    }

    expectations.each { |actual,expected|
      expect(actual.command).to eq(expected[:command])
      expect(actual.package).to eq(expected[:package])
      expect(actual.dependencies).to eq(expected[:dependencies])
    }
  end

  it 'throws exception if broken msg' do
    expect {Command.new("")}.to raise_error(Exception)
  end

  it 'throws exception if unknown command'
end
