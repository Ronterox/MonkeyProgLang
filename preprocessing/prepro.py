# define """% \1
# define %""" \1

"""
#run do
def generate_dataclass(name, fields)
    puts "\nfrom dataclasses import dataclass\n\n@dataclass\nclass \#{name}:"
    fields.split(',').each do |field|
        name, type = field.split(':')
        puts "    \#{name.strip}: \#{type.strip}"
    end
end
#end

#define for\s+(\d+)\.\.(\d+) for i in range(\1,\2)
#define \{$ :
#define \}$ pass
#define fn(\s+) def\1
#define (\w+)\s*=\s*class\s*\((.*)\) #run do\ngenerate_dataclass('\1','\2')\n#end
#define catch\s+(.*) try:\n\t\1\nexcept:\n\tpass
"""

# This is a code generated class
Person = class (name: str, age: int, gender: str)

"""%
# run do
3.times do |i|
    generate_dataclass("Carlitos\#{i}", 'power:float,level:int')
end
# end
%"""

# 3.times do |i|
print("Test #{i}")
# end


# idk what to do
def visualize(person: Person):
    print(f"{person.name} is {person.age} years old and {person.gender}")


fn looptwice(person: Person) {
    for 0..2 {
        visualize(person)
    }
}

name, age, gender = "Default", "0", "Default"
catch name, age, gender = open("file.txt", "r").read().split()
looptwice(Person(name, int(age), gender))
