"""
#run do
def generate_dataclass(name, fields)
    puts "from dataclasses import dataclass\n\n@dataclass\nclass \#{name}:"
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
"""


Person = class (name: str, age: int, gender: str)


def visualize(person: Person):
    print(f"{person.name} is {person.age} years old and {person.gender}")


fn looptwice(person: Person) {
    for 0..2 {
        visualize(person)
    }
}


me = Person("Rontero", 25, "Male")

looptwice(me)
