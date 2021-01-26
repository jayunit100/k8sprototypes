# Input file is any text/log file.  This program will remove the specificities
# from the input logs and print them out as a frequency plot.
# The example file here is a kubelet error file.
input_file="/data/workspace_files/Kubelet.Error"


def join(ar):
    out = ""
    for x in ar:
        out = out + x
    return out[0:25].replace("\"","")

def splitAny(inp):
    outputs = []
    buff = ""
    for i in inp:
        if i == " " or i == "-":
            # anything w/ numbers we just ignore
            if any(i.isdigit() for i in join(buff)):
                outputs += ""
            else:
                outputs += join(buff)
            buff = " "
        else:
            buff=buff + i
    return outputs




# Example of how to use those utils functions...
x = splitAny("A---v  g g 2o4i23mof02i3f -- aasdf -232342342 sfdsdfa")
print(join(x))

# 1 count lines in the file...
count = 0
f = open(input_file)
for line in f.readlines(  ): count += 1
print("TOTAL LINES")
print(count)
    

# 2 Build an index of all log types
line_set={"":0}
line_map={"":0}

for line in f.readlines(  ):
    split=splitAny(line)
    ss = join(split)
    if line_set.get(ss):
        line_set[ss] = line_set.get(ss)+1
    else:
        line_map[ss] = len(line_map)
        line_set[ss] = 1

print(line_set)
print(line_map)
#print(len(line_set))
#print(len(line_map))

# 3 now we have integer ids for all unique log "types", 
print("((((((((((((((((((((((((((((((((((((")
counttt=0
f = open(input_file)
for line in f.readlines(  ):
    split=splitAny(line)
    ss = join(split)

    print(line_map[ss])
print("))))))))))))))))))))))))))))))))))))")
