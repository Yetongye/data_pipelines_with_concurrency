# Go Image Processing Pipeline with Concurrency & Sequential Mode

## Overview

This project implements an image processing pipeline in Go to demonstrate the performance advantage of concurrency. The pipeline reads input images, resizes them while preserving their aspect ratios, converts them to grayscale, and saves the results. It supports both concurrent and sequential modes of execution, includes unit and benchmark testing, helping determine whether the startup will utilize Go concurrency.


## Image Input and Output

- **Input Location**: `images/`  
  Images to be processed should be placed in this directory (e.g., `images/image1.jpeg`, `images/image2.jpeg`).

<img src="images/image5.jpeg" alt="image5" width="400"/>


- **Output Location**: `images/output/`  
  All processed images (resized and grayscaled) will be saved here with the same filenames as the originals.
<img src="images/output/image5.jpeg" alt="image5output" width="400"/>


## How to Build and Run (Application)

### Build the executable:
This repository already includes pre-built executables:

**On macOS/Linux**
```bash
go build -o imagepipeline main.go
```

**On Windows**
```bash
go build -o imagepipeline.exe main.go
```
### Run the application:

**Concurrent mode (default):**
```bash
./imagepipeline -mode=concurrent
```

**Sequential mode:**
```bash
./imagepipeline -mode=sequential
```

**With profiling (optional):**
```bash
./imagepipeline -mode=concurrent -profile=true
```

Processed images will be saved in the `images/output/` directory.


### Error Handling and User Feedback

The application includes validation for image inputs and will not crash due to common issues. Below are possible error messages and their meanings:

| Scenario                                | Console Output                                                  |
|----------------------------------------|-----------------------------------------------------------------|
| Input file does not exist              | File does not exist: images/image7.jpeg                         |
| Image cannot be decoded                | Error loading image: images/image8.jpeg                         |
| Image is nil after reading             | Failed to read image: images/image9.jpeg                        |
| Output file cannot be created          | panic: open images/output/image7.jpeg: permission denied        |
| Invalid mode flag passed               | No specific error, the program defaults to `concurrent` silently |

To ensure a smooth run:
- Make sure all image paths in the input list are correct and point to real `.jpeg` files.
- Ensure `images/output/` exists and is writable.
- Use valid flags (`-mode=concurrent` or `-mode=sequential`).

## How to Test & What Has Been Tested

### Unit Tests (imageprocessing_test.go)
```bash
go test ./image_processing
```
Tests implemented:
- **Resize**: Checks output dimensions and aspect ratio.
- **Grayscale**: Verifies that pixel color is changed to grayscale.

### Benchmark Tests (main_test.go)
```bash
go test -bench=. -benchmem
```
Benchmarks implemented:
- `BenchmarkMainConcurrent`: Concurrent pipeline performance
- `BenchmarkMainSequential`: Sequential pipeline performance

All tests pass and validate correctness and performance.

## Aspect Ratio Preservation

Originally, all images were resized to 500x500, causing distortion. I updated the Resize() function to detect the original width and height of each image and calculate new dimensions that maintain the aspect ratio. This ensures that images are not stretched or squashed.

Sample logic:
```go
if width > height {
    newWidth = 500
    newHeight = height * 500 / width
} else {
    newHeight = 500
    newWidth = width * 500 / height
}
```

This adjustment improves the visual fidelity of the output images.

## Mode Comparison: Concurrent vs Sequential

| Mode         | Execution Time |
|--------------|----------------|
| Concurrent   | 47.21 ms       |
| Sequential   | 84.97 ms       |
| Speedup      | ≈ 44.5% faster |

### Technical Analysis

- **Concurrent Mode**  
  Utilizes Go's goroutines and channels to process images in parallel across different pipeline stages (load, resize, grayscale, save). Each stage works independently, allowing for overlapping computation and IO. This significantly reduces pipeline latency and improves throughput, especially when processing many images or large files.

- **Sequential Mode**  
  Performs image processing in a step-by-step fashion. Each image completes its full processing before the next one begins. While easier to debug and suitable for very small datasets or single-threaded systems, it fails to leverage modern multi-core CPUs efficiently.

### Recommendation for the company

In the context of aiming to reduce data engineering processing times, concurrent mode is strongly recommended. The concurrency-based pipeline aligns perfectly with the business goal:

- It reduces total execution time by over 40 percent in this implementation.
- It scales better with larger volumes of data or higher system loads.
- It demonstrates the feasibility and impact of adopting Go concurrency in real-world pipelines.

This replication confirms that concurrent design patterns in Go are not only theoretically appealing but also practically beneficial for reducing pipeline throughput time, which is critical in fast-paced, data-intensive environments.

## GenAI Tools

Generative AI tools were used throughout this project to assist with coding, testing, and performance analysis. For coding assistance, GPT-4 was used to help write modular Go functions for the concurrent pipeline, resizing. GPT assisted in designing and debugging unit tests. It is also used to analyze the performance, helping interpret benchmark results, calculate speedup ratios. And changed the readme document to a markdown file.
