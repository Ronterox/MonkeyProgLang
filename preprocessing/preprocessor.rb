#!/usr/bin/env ruby
require 'stringio'
require 'irb/ruby-lex'

def can_evaluate_command?(command)
  ['end', '}', ')'].each do |end_word|
    full_command = "#{command}\n#{end_word}"
    next unless Ripper.sexp(full_command)

    begin
      eval(full_command, TOPLEVEL_BINDING)
    rescue StandardError
      return nil
    end

    return true
  end

  nil
end

OPTIONS = ['-d', '--debug'].freeze
DEBUG = ARGV.any? { |arg| OPTIONS.include?(arg.downcase) }

def once
  yield
end

def d(message)
  return unless DEBUG

  puts message
end

# Captures stdout and returns it as a string
# @return [String]
def capture_stdout
  old = $stdout
  $stdout = StringIO.new
  begin
    yield
  ensure
    output = $stdout.string
    $stdout = old
  end
  output
end

def run(&block) = eval(capture_stdout(&block), TOPLEVEL_BINDING)

# Line by line code preprocessor
class Preprocessor
  def initialize(definitions = [])
    @definitions = definitions
    @identation = nil
    @command = nil
    @body = ''
  end

  # @param line [String]
  def preprocess(line)
    # d "definitions: #{@definitions.length}"
    substitutions = @definitions.count do |pattern, definition|
      # d "pattern: #{pattern} #{line.match?(pattern) ? 'matches' : 'does not match'}"
      line.gsub!(pattern, definition).tap { |v| d "replace: '#{pattern}' with '#{definition}'" if v }
    end

    return line.each_line.map { |l| preprocess(l) }.join if substitutions.positive?

    if /^\s*#\s*define\s+(?<pattern>.+?)\s+(?<definition>.+)$/ =~ line
      d "definition: #{pattern} => #{definition}"
      @definitions << [Regexp.new(pattern), "\"#{definition}\"".undump]
      return ''
    elsif !@command && /^\s*(?<identation>#+)(?<command>.*)/ =~ line
      d "command: #{command}"
      if can_evaluate_command?(command)
        @command = command.strip
        @identation = identation
      end
      return ''
    elsif @command && /^\s*#{@identation}(?!#)(?<close>.*)/ =~ line
      bprocessor = Preprocessor.new(@definitions)
      close = $~[:close]

      d "execute: #{@command}"
      d "body: #{@body}"
      d "close: #{close}"

      output = capture_stdout do
        eval <<~COMMAND
          #{@command}
          print <<HEREDOC
          #{@body}
          HEREDOC
          #{close}
        COMMAND
      end

      @identation = nil
      @command = nil
      @body = ''

      d "preprocess: #{output}"

      return output.each_line.map(&bprocessor.method(:preprocess)).join
    elsif @command
      d "append to body: #{line}"
      @body += line unless line.strip.empty?
      return ''
    end

    line
  end
end

ARGV.each do |arg|
  next if OPTIONS.include?(arg)

  File.open(arg, 'r') do |file|
    preprocessor = Preprocessor.new

    file.each_line do |line|
      d "LINE: #{line}"
      line = preprocessor.preprocess(line)
      print line unless DEBUG
    end
  end
end
