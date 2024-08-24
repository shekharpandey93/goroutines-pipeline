# Go Concurrency Patterns: Pipelines

1. Fan Out - Multiple functions can read from the same channel until that channel is closed; this is called fan-out
2. Fan In - A function can read from multiple inputs and proceed until all are chanels are closed onto a single channel that’s closed when all the inputs are closed. This is called fan-in.
