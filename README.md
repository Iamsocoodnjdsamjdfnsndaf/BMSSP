# BMSSP - A Go Implementation of a Near-Optimal Shortest Path Algorithm

This repository contains a Go implementation of the **Bounded Multi-Source Shortest Path (BMSSP)** algorithm. This algorithm is the core of the groundbreaking research from Tsinghua University that achieved the first major theoretical speedup for the single-source shortest path problem in over 40 years, effectively breaking the long-standing "sorting barrier" of Dijkstra's algorithm.

The implementation is based on the paper: **"A New Algorithm for the Single-Source Shortest Path Problem"** by Ruobing Chen, et al.

## Overview

The single-source shortest path (SSSP) problem is a fundamental challenge in computer science: find the shortest path from a starting "source" node to all other nodes in a weighted graph. For decades, Dijkstra's algorithm has been the standard solution. However, its performance is limited by a time complexity of roughly $O(m + n \log n)$, which becomes a bottleneck on massive graphs due to the cost of sorting nodes.

This new algorithm introduces a recursive, divide-and-conquer strategy that improves the complexity to $O(m \log^{2/3} n)$, making it significantly more scalable. This implementation serves as a practical, working demonstration of this cutting-edge theoretical concept.

## Features

* **Faithful Algorithm Implementation**: A correct and understandable implementation of the recursive BMSSP algorithm described in the paper.
* **Recursive Divide-and-Conquer Strategy**: Solves the SSSP problem by breaking it into smaller, bounded sub-problems.
* **Custom Data Structures**: Includes the specialized `DataStructureD` and `PriorityQueue` required to manage the algorithm's unique search frontiers.
* **Comprehensive Test Suite**: A robust set of unit and integration tests to verify the correctness of the implementation across various graph types and edge cases.
* **Self-Contained Core Logic**: The core `bmssp` package is self-contained and relies only on the Go standard library.

## Algorithm Pseudocode

The implementation is based on the following pseudocode from the research paper:

![BMSSP Algorithm Pseudocode](https://waifuvault.moe/f/e0086426-3adb-4dba-bc12-88e661673e10/GyAKU3rXQAAB_JI.jpg)

## Architecture & How It Works

The algorithm abandons the traditional, linear approach of Dijkstra and instead recursively partitions the problem.

### Core Components

* `bmssp/bmssp.go`: The main implementation of the recursive `BMSSP` function and its helper functions like `FindPivots` and `BaseCase`.
* `bmssp/data_structure_d.go`: The implementation of the specialized `DataStructureD` priority queue that drives the selection of sub-problems.
* `bmssp/priority_queue.go`: A standard min-priority queue implementation used by the data structures.
* `main.go`: An example executable showing how to initialize a graph and run the algorithm.

### Processing Flow

1.  **Pivot Selection**: The algorithm begins by selecting a set of "pivot" nodes (`P`) from the graph.
2.  **Initialization**: These pivots are inserted into the main controlling data structure, `D`.
3.  **Recursive Execution**: The algorithm enters a loop, pulling the highest-priority sub-problem (a set of source nodes `Si` with a boundary `Bi`) from `D`.
4.  **Sub-Problem Solving**: It then calls `BMSSP` recursively on this sub-problem with a reduced recursion depth (`l-1`). The base case (`l=0`) runs a bounded Dijkstra-like search.
5.  **Edge Relaxation & Partitioning**: After a recursive call returns, the algorithm "relaxes" the edges of the newly settled nodes. Critically, it partitions the neighboring nodes:
    * Nodes outside the child's boundary but inside the parent's are re-inserted into `D`.
    * Nodes inside the child's boundary are collected for a more efficient batch update.
6.  **Batch Update**: The collected nodes are added to `D` in a single batch operation, which is key to the algorithm's efficiency.
7.  **Termination**: The process repeats until the search frontiers are exhausted, resulting in the correct shortest-path distances for all reachable nodes.

## Usage

This repository is primarily a library demonstrating the algorithm. The `main.go` file provides a simple example of how to use it.

To run the example:

```bash
go run .
```

To use `bmssp` in your own project, import it and call the `BMSSP` function with your graph, source nodes, and parameters.

```go
import "your-module-path/bmssp"

// 1. Create your graph
g := &bmssp.Graph{...}

// 2. Define sources, boundary, and recursion depth
S := []int{0}
B := 1000
l := 3

// 3. Initialize distance map
dist := make(map[int]float64)
// ... initialize distances

// 4. Run the algorithm
bmssp.BMSSP(l, B, S, g, dist)

// 'dist' now contains the shortest path distances
```

## Building and Testing

You can build the example executable and run the test suite using standard Go commands.

```bash
# Build the example executable
go build -o bmssp-example .

# Run the comprehensive test suite
go test ./...

# Run the benchmarks to see performance on small/medium graphs
go test -bench=BenchmarkBMSSP ./bmssp
```

### Benchmark Results

Example benchmarks on an AMD Ryzen 9 9950X3D:

```
BenchmarkBMSSP_SmallGraph-32         977748          1036 ns/op
BenchmarkBMSSP_MediumGraph-32        117789         10139 ns/op
BenchmarkBMSSP_MultiSource-32         95085         12132 ns/op
```

**Note**: The true performance advantage of the BMSSP algorithm's improved time complexity is realized on extremely large-scale graphs. These benchmarks serve to validate the efficiency and correctness of this Go implementation on a smaller scale.
