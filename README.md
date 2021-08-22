## Node Traversal Algorithm

- Task:
    You are given a file which represents a graph that contains (at least) 6 nodes. Each node has a connection to each other node, and each connection distance is
    goes one way (i.e A->B has a distance m while B->A might have a distance n). The objective is to find the shortest path between a given source and destination
    node.

- Use:
    The program has 3 distinct modes for case handling:
    - Case file:
        If the user provides the path to a properly structured `case-file` (an example of such a file can be seen in `sampleCaseFile.yaml`), the algorithm will
        be run directly against the file.
    - Gen File:
        If the user provides a path to a properly structure `gen-file` (an example of such a file can be seen in `sampleGenFile.yaml`), a case yaml will be generated
        automatically and the algorithm will be run against it.
    - Manual:
        If neither a `case-file` nor a `gen-file` is provided when running the app from the CLI, the program will prompt the end user to provide n discrete nodes,
        to choose the distance range between them and the source node, as well as the location for the file. The program will then generate a case yaml and run the
        algorithm against the generated yaml.

- Notes
    Many design questions were left to me. Notably, the structure of the input file. I elected to use yaml, and to represent the graph as a flat list where the
    A->B connections would serve as the keys and the distances would be the values; this made extracting data from the file into the a flat data store very simple.

    I elected to use a parallelized brute force approach, as I find this easier to reason about. The algorithm scales to the 10x use case (i.e 60 nodes), but fails
    at the 100x case; further iteration would involve replacing the map and mutex combo with a the `sync.Map`, which was purpose built for such a case. Further
    augmentations could be made in the form of batch processing and memoization, should the problem remain intractable after the purely technical upgrades.

    This would benefit from further testing and hardening, especially around the three use methodologies; it will not be hard to produce unusual behavior.

