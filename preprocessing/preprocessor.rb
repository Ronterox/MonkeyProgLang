require 'stringio'

DEBUG = true

def once
  yield
end

def d(message)
  return unless DEBUG

  puts message
end

class Preprocessor
  def initialize
    @definitions = {}
    @identation = nil
    @command = nil
    @body = ''
  end

  def preprocess(line)
    @definitions.each do |pattern, definition|
      d "replace: '#{pattern}' with '#{definition}'" if line.include?(pattern)
      line.gsub!(pattern, definition)
    end

    if /^\s*#define\s+(?<pattern>\w+)\s+(?<definition>.+)$/ =~ line
      d "definition: #{pattern} => #{definition}"
      @definitions[pattern.strip] = definition.strip
      return ''
    elsif !@command && /^\s*(?<identation>#+)(?<command>.*)/ =~ line
      d "command: #{command}"
      @command = command.strip
      @identation = identation
      return ''
    elsif @command && /^\s*#{@identation}(?!#)(?<close>.*)/ =~ line
      bprocessor = Preprocessor.new
      d "preprocess: #{@command}"
      d "body: #{@body}"

      body = ''
      close = $~[:close]
      @body.each_line do |line|
        line = bprocessor.preprocess(line)
        body << line unless line.strip.empty?
      end

      captured = StringIO.new
      original = $stdout
      $stdout = captured

      begin
        eval <<~COMMAND
          #{@command}
          print <<~HEREDOC
          #{body}
          HEREDOC
          #{close}
        COMMAND
      ensure
        $stdout = original
        @command = nil
        @body = ''
      end

      return captured.string
    elsif @command
      d "append to body: #{line}"
      @body += line unless line.strip.empty?
      return ''
    end

    line
  end
end

File.open('example.mky', 'r') do |file|
  preprocessor = Preprocessor.new

  file.each_line do |line|
    d "LINE: #{line}"
    line = preprocessor.preprocess(line)
    print line unless DEBUG
  end
end
