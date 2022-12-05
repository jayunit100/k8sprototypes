# This script will take a log file and "parse" it into
# hours, minutes and seconds. It then will
# Make buckets for every hour minute, and print out all
# counts for all strings within those buckets.


# This parses a log line of format
# 2022-11-29 00:00:51 ...
# You can reimplement the "parse" function to extract minutes, hours, seconds in any way you want.
# Just make sure you set self.hours and self.minutes and self.seconds...
class LogLine:
    def __init__(self,f):
        x = self.parse(f)
        self.exit = x

    def parse(self, f):
        self.string=f
        if "2022" in self.string and self.string.find("-") > -1 and self.string.find(":") > -1:
            self.date = f.split(" ")[0]
            if len(f.split(" ")) < 3:
                self.error = f
                return -1


            self.year = self.date.split("-")[0]
            self.month = self.date.split("-")[1]
            if len(self.date.split("-")) < 3:
                self.error = "date -"
                return -4
            self.day = self.date.split("-")[2]

            self.time = f.split(" ")[1]

            self.time = self.time.rstrip()

            if self.time.find(":") == -1:
                self.error = f
                return -2

            self.hours = self.time.split(":")[0]
            self.minutes = self.time.split(":")[1]
            self.seconds = self.time.split(":")[2]
            self.error=None
            return 0

        else:
            self.error = f
            return -3


filename="gobuilds.log"
filename="gobuild.log"
# These are filters.  For example, anything with "sql" or "database" will be tagged with the "SQL" key.
# That way we can print out all log lines about things doing "SQL oriented" stuff during a given time bucket.
filters = {
    "SQL":["sql", "database"],
    "TF" :["tanzu-framework"],
    "GOBUILD":["Calling SetOptions", "Calling GetOptions", "No method Get", "buildtree", "tzinfo", "GetCluster","Source Provenance", "gobuilds.Sync", "gobuilds.M", "gobuilds.P"],
    "GOBUILD2":["Copying"],
    "GOBUILD2_SHA":["md6 sha1"],
    "SETUP":["No Match for argument", "Running transaction check", "yum", "Transaction Summary", "Dependency Installed"],
    "DOCKER":["Pulling fs layer", "docker tag", "Pulling from"],
    "OSX":["darwin"],
    "WINDOWS":["windows-amd"],
    "TF2":["IS_OFFICIAL_BUILD","harbor-repo.vmware.com/dockerhub","provider-bundle"],
    "TF2_PROVIDER":["provider"],
    "GOMOD":["go: downloading"],
    "TF2_CORE":["cli/core"],
    "GOLANG_BUILD":["/usr/local/go/bin/go","build -o"],
    #"GIT":["git"],
    "PLUGIN":["building plugin"],
    "PLUGIN_DONE":["succesfully built local"],
    "STANDALONG_PLUG":["standalone-plugins"],
    "PINNIPIED":["pinnipied"],
    "GOLANG_NOT_FOUND":["go: Command not found"],
    "MGMT":["management-cluster"],
    "SHA":["sha256"],
    "z_":[""],
    "zip":["gzip -c"],
    "X_COMPILE":["for os in windows linux darwin"],
    "POLL_HOST":["connecting to host"]
}

# converts a map to an array of strings with sorted keys for printing
# input: { a: 1, b: 2}
# output: ["a:1", "b:2"]
def format_map(m):
    t = []
    for k in sorted(m):
        if k != "z_":
            t.append(str(k)+":"+str(m[k]))
    return t

# read a string and return a map of all the tags which the string
# contains.
# input : "The brown cow is nice"
# output: { animals:1 , colors:1} 
# reads the global "filters" variable as the source. 
def tag(f):
    vals = {}
    for k in filters:
        print("\t", "checking", k)
        for v in filters[k]:
            if v in f:
                if k not in vals:
                    vals[k] = 1
                else:
                    vals[k]+=1
    return vals




with open(filename,'r') as f:
    lines = f.readlines()
    hourminutes={

    }
    labels_for_timebucket={

    }

    for s in lines:
        if len(s) < 2:
            print("skip")
        else:
            ll = LogLine(s)
            if ll.error == None:
                key = ll.hours+":"+ll.minutes
                if key in hourminutes:
                    hourminutes[key] += 1
                else:
                    hourminutes[key] = 1

                labels = tag(ll.string)
                for l in labels:
                    if key not in labels_for_timebucket:
                        labels_for_timebucket[key] = {}
                    if l not in labels_for_timebucket[key]:
                        labels_for_timebucket[key][l] = 0
                    # prints out a 1 for each label, but we should add it eventually....
                    labels_for_timebucket[key][l] += labels[l]
                    #if  labels_for_timebucket[key][l] > 4:
                    #    print("too many labels")
                    #    print(labels_for_timebucket)
                    #    exit(-1)
            else:
                print(ll.exit, "ERROR ------>", ll.error)

    print(hourminutes)
    print(labels_for_timebucket)

    print("********************************************")
    global_counts_of_all_labels = {
    }

    for key in sorted(hourminutes):
        print(key, "===> ", format_map(labels_for_timebucket[key]))

        for task in labels_for_timebucket[key]:
            if task not in global_counts_of_all_labels:
                global_counts_of_all_labels[task] = 0
            global_counts_of_all_labels[task] += labels_for_timebucket[key][task]


    print("***************** global *********************")
    for key in sorted(global_counts_of_all_labels):
        if key != "null":
            print(key, global_counts_of_all_labels[key])
