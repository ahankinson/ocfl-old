package libocfl

import (
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"sync"
)

type OCFLRule struct {
	Code         string
	Message      string
	CheckError   func(obj *OCFLObject) bool
	CheckWarning func(obj *OCFLObject) bool
}

var e001Check = &OCFLRule{
	Code:    "E001",
	Message: "Missing OCFL NamAsTe file",
	CheckError: func(obj *OCFLObject) bool {
		fpth := path.Join(obj.path, "0=ocfl_object_1.0")
		_, err := os.Stat(fpth)

		if err != nil {
			return false
		}
		return true
	},
}

var e002Check = &OCFLRule{
	Code:    "E002",
	Message: "Namaste file name does not match Inventory OCFL Version",
	CheckError: func(obj *OCFLObject) bool {
		log.Debug().Msg("Checking E002")
		return true
	},
}

var e003Check = &OCFLRule{
	Code:    "E003",
	Message: "Missing inventory.json in object root",
	CheckError: func(obj *OCFLObject) bool {
		log.Debug().Msg("Checking E003")
		return true
	},
}

var w001Check = &OCFLRule{
	Code:    "W001",
	Message: "Some warning",
	CheckWarning: func(obj *OCFLObject) bool {
		return false
	},
}

// Gather all the error checks into a slice
var errorChecks = []*OCFLRule{
	e001Check,
	e002Check,
	e003Check,
}

// Gather all the warning checks into a slice
var warnChecks = []*OCFLRule{
	w001Check,
}

func runValidation(obj *OCFLObject) {
	ruleChannel := make(chan *OCFLRule)
	var wg sync.WaitGroup

	// Create up to NUM workers, each responsible for
	// reading from a channel and doing work. This will continue
	// until the channel is exhausted.
	for worker := 0; worker < NumWorkers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for rule := range ruleChannel {
				log.Debug().Str("code", rule.Code).Msg("Queued up")
				if rule.CheckError != nil {
					valid := rule.CheckError(obj)
					if !valid {
						log.Error().
							Str("code", rule.Code).
							Str("path", obj.path).
							Bool("valid", valid).
							Msg(rule.Message)
					}
				} else if rule.CheckWarning != nil {
					warn := rule.CheckWarning(obj)

					// If we return a 'warning' status, set the
					// 'warn' flag to true by inverting the returned
					// status. So the warning checker returns false if
					// we need to flag the result as a warning, but the
					// error message will return 'warn': true.
					if !warn {
						log.Warn().
							Str("code", rule.Code).
							Str("path", obj.path).
							Bool("warn", !warn).
							Msg(rule.Message)
					}
				}
			}
		}()
	}

	log.Debug().Msg("Finished allocating workers")

	// Queue up all the errors and warnings by sending them
	// into the rule channel.

	// Add error-check tasks to the channel
	for _, errRule := range errorChecks {
		ruleChannel <- errRule
	}

	// Add warning-check tasks to the channel
	for _, warnRule := range warnChecks {
		ruleChannel <- warnRule
	}

	// Close the rule channel. No other rules can be added, but
	// the rules in the channel will be processed until there are
	// no more to be done.
	close(ruleChannel)

	// Wait for the worker pool to finish.
	wg.Wait()
}
