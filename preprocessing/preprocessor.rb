require 'stringio'

DEBUG = false

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

$*.each do |path|
  File.open(path, 'r') do |file|
    preprocessor = Preprocessor.new

    file.each_line do |line|
      d "LINE: #{line}"
      line = preprocessor.preprocess(line)
      print line unless DEBUG
    end
  end
end
