# parse sonobuoy junit.xml into a usable format

import xml.dom.minidom
import statistics
# valid tags, from https://github.com/kubernetes/test-infra/blob/master/jobs/validOwners.json 
TAGS = {
   "sig-api-machin",
   "sig-apps",
   "sig-auth",
#   "sig-autoscaling",
#    "sig-azure",
#   "sig-big-data",
   "sig-cli",
#   "sig-cloud-provider-gcp",
#   "sig-cluster-lifecycle",
#   "sig-contributor-experience",
#   "sig-docs",
#   "sig-multicluster",
   "sig-instrum",
   "sig-netw",
   "sig-node",
#    "sig-onprem",
#    "sig-openstack",
#    "kubernetes-release",
#   "sig-scalability",
   "sig-scheduling",
#   "sig-service-catalog",
   "sig-storage",
#   "sig-testing",
#   "sig-ui",
#   "sig-windows"
#  "sig-release"
}

def string(confTest):
    return str(confTest.name)+","+str(confTest.time)+"[[[["+str(confTest.tags)+"]]]]"

class ConformanceRun:
    '''
        sig-net -> x
                   y
                   a
                   b
        sig-win -> x
                   y 
    '''
    def __init__(self):
        self.allTests = []

    def summary(self) -> str: 
        return str(" ---> " + str(len(self.allTests)))

    def add(self, t):
        self.allTests.append(t)
    
    def total(self) -> float:
        return len(self.allTests)

    def filter(self, tags) -> list:
        #print(len(self.allTests))
        def matches(cTest):
            bb = cTest.match(tags)
            return bb

        return list(filter(matches, self.allTests))

    def total(self, tags ):
        f = self.filter(tags)
        sum = 0
        for i in f:
            sum += i.time
        return sum

    # count total # of tests w/ these tags 
    def count(self, tags) -> float:
        return len(self.filter(tags))

    def stddev(self, tags) -> float:
        f = self.filter(tags)
        all = []
        for i in f:
            all.append(i.time)
        if len(all) == 0:
            return 0
        return statistics.stdev(all)

    def average(self, tags) -> float:
        f = self.filter(tags)
        total = 0
        sum = 0
        for i in f:
            total += 1
            sum += i.time
        if total == 0:
            return 0
        return sum / total

    def longest(self, tags):
        #print("longest {tags}",tags)
        f = self.filter(tags)
        maxx = -1
        confTest = None
        for i in f:
            maxx = max(i.time, maxx )
            # it changed, so update the return val.
            if i.time == maxx:
                confTest = i
        return confTest

    def shortest(self, tags) -> float:
        minn = 100000
        f = self.filter(tags)
        for i in f:
            #print(i.time, i.name)
            minn = min(i.time, minn )  
        return minn

class ConformanceTest:
    # unique identified 
    def id(self):
        return self.source+self.name

    # no timestamp means, skipped...
    def valid(self) -> bool:
        return float(self.time) > 0

    # return true if this test matches all tags
    def match(self, intags):
        # sig-network
        for t in intags:
            if t not in self.tags:
                return False
        return True

    # node is a type of xml.dom.minidom.Element 
    def __init__(self, node):
        self.time = -1
        self.name = "none"
        self.tags = {}
        self.source = ""
        try:
            #print(node.toxml())
            #print(node.attributes.items())

            if node.getAttribute('name'):
                self.name = node.getAttribute("name")

            if float(node.getAttribute("time")) > 0:
                self.time = float(node.getAttribute("time"))

            for tt in TAGS:
                if tt in self.name:
                    #print(f"\tadding tag {tt} to {self.name}")
                    self.tags[tt]=1
        
        except Exception as e:
            print("fail",e)

import os

def main():
    root = "./k8s-conformance/v1.22/"
    for dirr in os.listdir(root):
        print("########",dirr,"########")
        process(root+dirr+"/junit_01.xml")


def process (f):
    run = ConformanceRun()

    Load_XML = xml.dom.minidom.parse(f)
    print (Load_XML.firstChild.getAttribute("time"))
    
    while Load_XML.firstChild.childNodes.length > 1:
        try:
            node = Load_XML.firstChild.childNodes.pop()
            if type(node) == xml.dom.minidom.Element:
                confTest = ConformanceTest(node)
                if confTest.valid():
                    run.add(confTest)
        except Exception as e:
            print(type(node),"<-ERROR", e)

    #print(run.summary())
    print("*****************************")
    print(f"\t\tshrt \t avg \t dev \t  long \t cnt \t totaltime(s)")
    summarys = {}
    table = []

    # CSV output of all timings... 
    for t in TAGS:
        s = format(run.shortest([t]), '.1f')
        a = format(run.average([t]), '.1f')
        d = format(run.stddev([t]), '.1f')

        l = "x"
        longestDesc = "x"

        confTest = run.longest([t])
        if confTest is not None:
            l = format(confTest.time , '.1f')
            summarys[f"longest {t}"]=f"{confTest.name} {confTest.time}"

        c = run.count([t])
        to = format(run.total([t]),'.1f')

        table.append([t,s,a,d,l,c,to])

        print(f"{t} \t {s}s \t {a} \t {d} \t {l}s\t {c}\t {to}")

    for t in summarys:
        print()
        print(t,":",summarys[t])
        
if __name__ == "__main__":
    main()

#Otherwise, its a "text" field
#<class 'xml.dom.minidom.Text'>
#if type(node) == xml.dom.minidom.Text: 
#    print(node.data)
