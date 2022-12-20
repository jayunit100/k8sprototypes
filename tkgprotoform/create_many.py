### TODO convert this to golang 


import subprocess

# General TKG Flow

import subprocess
import os

codes = []


def replace_write(fName, oName, r):
    #input file
    fin = open(fName, "rt")
    #output file to write the result to
    fout = open(oName, "wt")
    #for each line in the input file
    for line in fin:
        l = line
        for k in r:
            l = l.replace(k, r[k])
        # Now write the line...
        fout.write(l)
    fin.close()
    fout.close()

def do(task, cname, ip):
    debug=True
    result = None
    err = None
    if task=="create":
        outName="wl-{0}".format(cname)
        replace_write("wlcc.yaml", outName, { "172.16.202.134":ip, "cluster-name":cname } )
        proc = subprocess.Popen(["tanzu","cluster","create","-f",outName, cname],stdout=subprocess.PIPE)
        result,err = proc.communicate()
        exit_code = proc.wait()

    elif task=="destroy":
        proc = subprocess.Popen(["tanzu","cluster","delete","-y", name],stdout=subprocess.PIPE)
        debug=False
        result,err = proc.communicate()
        exit_code = proc.wait()

    else:
        print("not valid",task)
        exit(3)
    if debug:
        print(f"result={result} error={err}")
        codes.append(f"{task} {exit_code}")
        print(exit_code)

# use nmap to find this range, initially, TODO automatethat !
def set_ips():
    base="172.16.202"
    ips = []
    for i in range(48,130):
        ip="{0}.{1}".format(base,i)
        print(ip)
        ips.append(ip)
    return ips

def main():
    ips = set_ips()

    for i in range(2,21,1):
        print(f"start {i}")
        do("create", f"aasdf{i}", ips.pop())
        #do("destroy", f"aasdf{i}")
        print(f"done...{i}")
        print(f"exit codes so far: {codes}")
        i = i + 1

if __name__ == "__main__":
    main()
