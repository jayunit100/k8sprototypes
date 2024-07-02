# A hacky simulation of a k8s scheduler w/ different types of nodes

class Pod:
    # Same as node
    def __init__(self, cores: int, mem: int):
        self.cores=cores
        self.mem=mem

# Find all nodes which match this pod
def fit(pod: Pod, nn: []) -> ([],set):
    print("fitting over" , len(nn))
    reasons = []
    possible_nodes = set()
    
    for i,n in enumerate(nn):
        if n.freecores < pod.cores:
            reasons.append({"free cores too low ":n.freecores})
        elif n.freemem < pod.mem:
            reasons.append({"mem too low":n.freemem })
        else:
            reasons.append(None)
            ### Add index of the node
            possible_nodes.add(i)
    return reasons, possible_nodes

class Cluster:
    def __init__(self, nodes: []):
        self.nodes=nodes

class Node:
    # cores (int) and mem is in G
    def __init__(self, cores: int, mem: int):
        self.cores=cores
        self.mem=mem
        self.freecores=cores
        self.freemem=mem
    def __str__(self):
        return str("Node:",cores," core,",mem,"gig")

    ### Every pod you run, reduces available resources.
    ### TODO simulate burstable.  not sure how.  maybe ask RJ.
    def run(self, p: Pod):
        self.freecores=self.freecores-p.cores
        self.freemem=self.freemem-p.mem


d_node_types = {
        "small":  {"cores": 4, "memory": 16}, 
        "memory": {"cores": 4, "memory": 64},
        "large":  {"cores": 16, "memory": 256},
        "mega":   {"cores": 66111, "memory": 11512}
} 

d_topology = { 
     "small": 1110,             
     "memory": 120,            
     "large":100,
     "mega":111              
}                            
d_workload_types = {
        "e-worker": {"cores":4,"memory":16},
        "search": {"cores":8,"memory":128},
        "web": {"cores":6,"memory":32}
}

d_workloads = [
    {"type":"e-worker", "count":10},
    {"type":"web",      "count":20} 
]

class Scheduler:
    def __init(self):
        pass

    # Find a single node which is recommended for running this pod
    def schedule(self, pod: Pod, nodes) -> int:
        reasons, fits = fit(pod, nodes)
        if len(fits) == 0:
            print("node failure /schedule reasons", reasons, " fits ==", fits)
        if len(fits) > 0:
            node_idx = fits.pop()
            return node_idx
        else:
            return -1

class Simulation:
    ### node_types has info on the "Types" of nodes.  Wont change often, but confi
    ### Topology has nodes info, i.e. how many of each type.
    def __init__(self, node_types=d_node_types, topo=d_topology, workload_types=d_workload_types, workloads=d_workloads):
        self.nodes = []
        ### Make nodes
        # small,10 ; medium,8 ; ...
        for k,v in topo.items():
            cores  = node_types[k]["cores"]
            memory = node_types[k]["memory"]
            for i in range(v):
                n = Node(cores, memory)
                self.nodes.append(n)
        print("nodes",len(self.nodes))
        s = Scheduler()
        scheduled = 0
        unscheduled = 0
        total = 0
        for p in workloads:
            print(p)
            ### Lookup what type
            wt = workload_types[p["type"]]
            ### Schedule this new pod against all existing nodes
            ### Schedule has no side effects, it just finds a home for a pod
            ### TODO make schedule actually mutate the node object.
            ### TODO Descide how scheudle prints cores = 0???
            
            podd = Pod(wt["cores"], wt["memory"])
            for i in range(p["count"]):
                target_node = s.schedule(podd, self.nodes)

                print("Scheduled",wt)
            
                if target_node > -1:
                    ### This should be done via schedule - part of the TODO above
                    self.nodes[target_node].run(podd)
                    scheduled += 1
                else:
                    unscheduled += 1
                    print("pod is unschedulable:", podd)
        print("done, scheduled = ",scheduled, " out of ", scheduled+unscheduled, "%", 100*scheduled/(scheduled+unscheduled))

_nt = { "small" : 
       { 
            "cores": 4, 
            "memory": 16
       }
}
_t =  { 
        "small" : 2 
}
_wt = { 
        "e" : {
            "cores":4, 
            "memory":16
        } 
}
_w =  [ 
        { 
            "type": "e", 
            "count": 2  
        }
]
### Two nodes  (4x16G) = 8 by 32 g of capacity.  each W is 4x16.  2 of them is 8/32G.  So should schedule 100%...
Simulation(_nt, _t, _wt, _w )
