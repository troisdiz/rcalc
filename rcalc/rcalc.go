package rcalc

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

func RunFile(programPath string, progArgs []string, outputAsJson bool, stackDataFolder string, createFolder bool, debugMode bool) ([]string, error) {
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
				return nil, fmt.Errorf("cannot create stack data folder %s", stackDataFolder)
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

	var stack = CreateStack()

	var system = CreateSystemInstance()

	// interpret cmd
	progFileBytes, err := os.ReadFile(programPath)
	if err != nil {
		return nil, fmt.Errorf("errors reading program file %s (%w)", programPath, err)
	}
	progFileText := string(progFileBytes)

	actions, parseErr := ParseToActions(progFileText, "ParseFile", Registry)
	if parseErr != nil {
		sugaredLogger.Errorf("Parsing error(s): %v", parseErr)
		return nil, fmt.Errorf("czannot parse file %s (%w)", programPath, parseErr)
	} else {
		runtimeContext := CreateRuntimeContext(system, stack)
		err := stack.StartSession()
		if err != nil {
			return nil, fmt.Errorf("cannot create stack: %w", err)
		}

		// write args to stack
		for argIdx, arg := range progArgs {
			fmt.Println(arg)
			actions, err := ParseToActions(arg, "ArgLexer", Registry)
			if err != nil {
				return nil, fmt.Errorf("error parsing argument nb %d, %s", argIdx+1, arg)
			}
			if len(actions) != 1 {
				return nil, fmt.Errorf("error parsing argument nb %d, %s (too many actions)", argIdx, arg)
			}
			//TODO check type
			runtimeContext.RunAction(actions[0])
		}

		// execute program
		for _, action := range actions {
			err := runtimeContext.RunAction(action)
			if err != nil {
				return nil, fmt.Errorf("error executing file: %w", err)
			}
			//TODO handle quit gracefully in programs
			/*
				if system.shouldStop() {
					err := stack.CloseSession()
					if err != nil {
						sugaredLogger.Errorf("Error while closing session")
						return
					}
					return
				}*/
		}
		// DisplayStack(stack, "", 4, false)

		// write stack to stdout
		if outputAsJson {
			fmt.Printf("[\n")
			for _, stackElt := range stack.elts {
				//TODO need to jsonsify!
				fmt.Printf("     \"%s\",\n", stackElt.display())
			}
			fmt.Printf("]\n")
		} else {
			for _, stackElt := range stack.elts {
				fmt.Println(stackElt.display())
			}
		}
		err = stack.CloseSession()
		if err != nil {
			return nil, fmt.Errorf("cannot close stack: %w", err)
		}
	}

	return nil, fmt.Errorf("TODO finish implementation")
}

func RunRepl(stackDataFolder string, createFolder bool, debugMode bool) {

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
