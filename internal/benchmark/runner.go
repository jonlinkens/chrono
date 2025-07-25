package benchmark

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	Phrase          string
	Warmups         int
	Runs            int
	Timeout         time.Duration
	CalibrationRuns int
	SkipCalibration bool
	Command         []string
	UseCli          bool
}

type Result struct {
	Duration time.Duration
	Found    bool
}

func Run(config Config, shellOverhead time.Duration) Result {
	cmd := exec.Command(config.Command[0], config.Command[1:]...)

	if config.Phrase == "" {
		return runCommandCompletion(cmd, config.Timeout, shellOverhead)
	}

	return runPhraseDetection(cmd, config.Phrase, config.Timeout, shellOverhead)
}

func runCommandCompletion(cmd *exec.Cmd, timeout time.Duration, shellOverhead time.Duration) Result {
	startTime := time.Now()

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting command:", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	if timeout > 0 {
		select {
		case <-done:
			duration := time.Since(startTime)
			adjustedDuration := max(duration-shellOverhead, 0)
			return Result{Duration: adjustedDuration, Found: true}
		case <-time.After(timeout):
			killProcess(cmd)
			return Result{Found: false}
		}
	} else {
		<-done
		duration := time.Since(startTime)
		adjustedDuration := max(duration-shellOverhead, 0)
		return Result{Duration: adjustedDuration, Found: true}
	}
}

func runPhraseDetection(cmd *exec.Cmd, phrase string, timeout time.Duration, shellOverhead time.Duration) Result {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdout pipe:", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderr pipe:", err)
	}

	startTime := time.Now()

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting command:", err)
	}

	done := make(chan time.Duration, 1)
	cancel := make(chan struct{})
	cmdFinished := make(chan struct{})

	go scanOutput(stdout, phrase, startTime, done, cancel)
	go scanOutput(stderr, phrase, startTime, done, cancel)

	go func() {
		cmd.Wait()
		close(cmdFinished)
	}()

	if timeout > 0 {
		select {
		case duration := <-done:
			close(cancel)
			killProcess(cmd)
			adjustedDuration := max(duration-shellOverhead, 0)
			return Result{Duration: adjustedDuration, Found: true}
		case <-time.After(timeout):
			close(cancel)
			killProcess(cmd)
			return Result{Found: false}
		case <-cmdFinished:
			close(cancel)
			return Result{Found: false}
		}
	} else {
		select {
		case duration := <-done:
			close(cancel)
			killProcess(cmd)
			adjustedDuration := max(duration-shellOverhead, 0)
			return Result{Duration: adjustedDuration, Found: true}
		case <-cmdFinished:
			close(cancel)
			return Result{Found: false}
		}
	}
}

func scanOutput(reader io.Reader, phrase string, startTime time.Time, done chan time.Duration, cancel chan struct{}) {
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		select {
		case <-cancel:
			return
		default:
		}
		line := scanner.Text()
		if strings.Contains(line, phrase) {
			select {
			case done <- time.Since(startTime):
			default:
			}
			return
		}
	}
}

func killProcess(cmd *exec.Cmd) {
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}
