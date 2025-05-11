package shell

import (
	"os"
	"os/exec"
)

func (s *Shell) editScript(scriptName string, filetype string) error {
	if os.Getenv("EDITOR") == "" {
		s.shell.Println(messagePrompt + "Error: EDITOR environment variable is not set.")
		return nil
	}
	cmd := exec.Command(os.Getenv("EDITOR"), "-filetype", filetype, scriptName)
	cmd.Stdin = s.stdin
	cmd.Stdout = s.stdout
	cmd.Stderr = s.stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
