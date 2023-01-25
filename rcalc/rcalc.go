package rcalc

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

func Run(stackDataFolder string, createFolder bool, debugMode bool) {

	defer func() {
		logger := GetLogger()
		if logger == nil {
			_, _ = fmt.Fprintln(os.Stderr, "Exiting rcalc, logger is nil")
		} else {
			logger.Info("Exiting rcalc")
		}
	}()

	if createFolder {
		if _, err := os.Stat(stackDataFolder); os.IsNotExist(err) {
			err := os.Mkdir(stackDataFolder, 0755)
			if err != nil {
				fmt.Printf("Error creating %s : %s\n", stackDataFolder, err.Error())
				return
			}
		}
	}
	logFilePath := path.Join(stackDataFolder, "rcalc-debug.log")
	fmt.Println(logFilePath)
	if debugMode {
		InitDevLogger(logFilePath)
	} else {
		InitProdLogger(logFilePath)
	}

	GetLogger().Info("Start rcalc")
	sugaredLogger := GetLogger()
	stackDataFilePath := path.Join(stackDataFolder, "stack.protobuf")

	var stack = CreateSaveOnDiskStack(stackDataFilePath)

	var message = ""
	var system = CreateSystemInstance()
	for {
		// print stack
		DisplayStack(stack, message, 3, true)

		// print prompt
		input := bufio.NewScanner(os.Stdin)

		// wait for cmd
		input.Scan()

		// interpret cmd
		var cmds = input.Text()

		// Message to display above (temp way of doing this)
		message = ""

		actions, parseErr := ParseToActions(cmds, "InteractiveShell", Registry)
		if parseErr != nil {
			sugaredLogger.Errorf("Parsing error(s): %v", parseErr)
			message = parseErr.Error()
		} else {
			runtimeContext := CreateRuntimeContext(system, stack)
			err := stack.StartSession()
			if err != nil {
				return
			}
			for _, action := range actions {
				err := runtimeContext.RunAction(action)
				if err != nil {
					message = err.Error()
					// in case of error, stop evaluation
					break
				} else {
					message = ""
				}
				if system.shouldStop() {
					err := stack.CloseSession()
					if err != nil {
						sugaredLogger.Errorf("Error while closing session")
						return
					}
					return
				}
			}
			err = stack.CloseSession()
			if err != nil {
				return
			}
		}
	}
}

func checkTypesForAction(s *Stack, a Action) (bool, error) {
	elts, _ := s.PeekN(a.NbArgs())
	ok, err := a.CheckTypes(elts...)
	if err != nil {
		return false, err
	} else {
		return ok, nil
	}
}
