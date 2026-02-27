require 'stringio'

def once
  yield
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
      line.gsub!(pattern, definition)
    end

    if /^\s*#define\s+(?<pattern>\w+)\s+(?<definition>.+)$/ =~ line
      @definitions[pattern.strip] = definition.strip
      return ''
    elsif !@command && /^\s*(?<identation>#+)(?<command>.*)/ =~ line
      @command = command.strip
      @identation = identation
      return ''
    elsif @command && /^\s*#{@identation}(?!#)(?<close>.*)/ =~ line
      bprocessor = Preprocessor.new

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
      @body += line unless line.strip.empty?
      return ''
    end

    line
  end
end

File.open('example.mky', 'r') do |file|
  preprocessor = Preprocessor.new

  file.each_line do |line|
    print preprocessor.preprocess(line)
  end
end
