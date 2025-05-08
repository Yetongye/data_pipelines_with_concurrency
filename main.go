package main

import (
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"log"
	"os"
	"strings"

	"flag"
	"runtime/pprof"
	"time"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			// Check if the file exists
			if _, err := os.Stat(p); os.IsNotExist(err) {
				log.Printf("File does not exist: %s\n", p)
				continue
			}

			img := imageprocessing.ReadImage(p)
			// Check if the image was loaded successfully
			if img == nil {
				log.Printf("Error loading image: %s", p)
				continue
			}

			log.Printf("Loading image: %s", p)

			job := Job{InputPath: p,
				Image:   img,
				OutPath: strings.Replace(p, "images/", "images/output/", 1)}
			// Check if the image was read successfully
			if job.Image == nil {
				log.Printf("Failed to read image: %s", p)
				continue
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input { // Read from the channel
			imageprocessing.WriteImage(job.OutPath, job.Image)
			log.Printf("Saving image: %s", job.OutPath)
			out <- true
		}
		close(out)
	}()
	return out
}

// runConcurrent processes images concurrently
// using goroutines and channels
// This is the recommended way to process images as it takes advantage of concurrency
// and is faster than the sequential version
// It is also more complex and harder to debug but it is worth it for the performance gain
func runConcurrent(paths []string) {
	ch1 := loadImage(paths)
	ch2 := resize(ch1)
	ch3 := convertToGrayscale(ch2)
	result := saveImage(ch3)

	for success := range result {
		if success {
			log.Println("Success!")
		} else {
			log.Println("Failed to save image")
		}
	}
}

// runSequential processes images one by one
// without using goroutines
// This is useful for testing and debuggingthe image processing functions
// It is not recommended for production use as it does not take advantage of concurrency
// and may be slower than the concurrent version
func runSequential(paths []string) {
	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			log.Printf("File not found: %s", p)
			continue
		}
		img := imageprocessing.ReadImage(p)
		if img == nil {
			log.Printf("Failed to read image: %s", p)
			continue
		}
		img = imageprocessing.Resize(img)
		img = imageprocessing.Grayscale(img)

		outPath := strings.Replace(p, "images/", "images/output/", 1)
		imageprocessing.WriteImage(outPath, img)
		log.Printf("Wrote: %s", outPath)
	}
}

func main() {

	mode := flag.String("mode", "concurrent", "Set mode: 'concurrent' or 'sequential'")
	profile := flag.Bool("profile", false, "Enable CPU and memory profiling")
	flag.Parse()

	if *profile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		defer func() {
			mf, err := os.Create("mem.prof")
			if err == nil {
				pprof.WriteHeapProfile(mf)
				mf.Close()
			}
		}()
	}

	imagePaths := []string{"images/image7.jpeg",
		"images/image8.jpeg",
		"images/image5.jpeg",
		"images/image6.jpeg",
	}

	start := time.Now()

	if *mode == "sequential" {
		runSequential(imagePaths)
	} else {
		runConcurrent(imagePaths)
	}

	elapsed := time.Since(start)
	log.Printf("Execution time for mode '%s': %s\n", *mode, elapsed)
}
