require 'stringio'

DEBUG = false

def once
  yield
end

def d(message)
  return unless DEBUG

  puts message
end

# Line by line code preprocessor
class Preprocessor
  def initialize
    @definitions = []
    @identation = nil
    @command = nil
    @body = ''
  end

  # @param line [String]
  def preprocess(line)
    @definitions.each do |pattern, definition|
      d "replace: '#{pattern}' with '#{definition}'" if line.include?(pattern)
      line.gsub!(Regexp.new(pattern), definition)
    end

    if /^\s*#define\s+(?<pattern>.+)\s+(?<definition>.+)$/ =~ line
      d "definition: #{pattern} => #{definition}"
      @definitions << [pattern, definition]
      return ''
    elsif !@command && /^\s*(?<identation>#+)(?<command>.*)/ =~ line
      d "command: #{command}"
      @command = command.strip
      @identation = identation
      return ''
    elsif @command && /^\s*#{@identation}(?!#)(?<close>.*)/ =~ line
      bprocessor = Preprocessor.new
      close = $~[:close]

      d "execute: #{@command}"
      d "body: #{@body}"

      captured = StringIO.new
      original = $stdout
      $stdout = captured

      begin
        eval <<~COMMAND
          #{@command}
          print <<HEREDOC
          #{@body}
          HEREDOC
          #{close}
        COMMAND
      ensure
        $stdout = original
        @identation = nil
        @command = nil
        @body = ''
      end

      output = captured.string
      body = ''

      d "preprocess: #{output}"

      output.each_line do |line|
        line = bprocessor.preprocess(line)
        body << line unless line.strip.empty?
      end

      return body
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
