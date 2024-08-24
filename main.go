package main

import (
	"goroutines-pipeline/services"
	"image"
	"strings"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Name      string
}

func load(imagesPath []string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, imgPath := range imagesPath {
			job := Job{InputPath: imgPath,
				OutPath: strings.Replace(imgPath, "img/", "img/output/", 1),
				Name:    strings.Split(imgPath, "/")[1],
			}
			job.Image = services.ReadImage(imgPath)
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(jobData <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range jobData {
			job.Image = services.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func changeColor(jobData <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range jobData {
			job.Image = services.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func save(jobData <-chan Job) <-chan string {
	out := make(chan string)
	go func() {
		for job := range jobData {
			services.WriteImage(job.OutPath, job.Image)
			out <- job.Name
		}
		close(out)
	}()
	return out
}

func loadImage(imagesPath []string) []Job {
	var jobs []Job
	for _, imgPath := range imagesPath {
		job := Job{InputPath: imgPath,
			OutPath: strings.Replace(imgPath, "img/", "img/output/", 1),
			Name:    strings.Split(imgPath, "/")[1],
		}
		job.Image = services.ReadImage(imgPath)
		jobs = append(jobs, job)
	}
	return jobs
}

// Collect / Fan out
func resizeImage(job *[]Job) <-chan Job {
	out := make(chan Job)
	for _, job := range *job {
		go func(job Job) {
			job.Image = services.Resize(job.Image)
			out <- job
		}(job)
	}
	return out
}

// Collect / Fan in
func collectJobs(input <-chan Job, imageCnt int) []Job {
	var resizedJobs []Job
	for i := 0; i < imageCnt; i++ {
		job := <-input
		resizedJobs = append(resizedJobs, job)
	}
	return resizedJobs
}

func saveImages(jobs *[]Job) {
	for _, job := range *jobs {
		services.WriteImage(job.OutPath, job.Image)
	}
}

func main() {
	imagePaths := []string{"img/image4.jpeg",
		"img/image5.jpeg",
		"img/image6.jpeg",
	}

	//chan1 := load(imagePaths)
	//chan2 := resize(chan1)
	//chan3 := changeColor(chan2)
	//writeResults := save(chan3)
	//for val := range writeResults {
	//	fmt.Printf("%s Process Successfully\n", val)
	//}

	jobs := loadImage(imagePaths)
	// fan out
	out := resizeImage(&jobs)
	// fan in
	resizedJobs := collectJobs(out, len(jobs))
	saveImages(&resizedJobs)

}
