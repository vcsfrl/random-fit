package service

type Printer interface {
	// Print is a convenience method to Print to the defined output, fallback to Stderr if not set.
	Print(i ...interface{})

	// Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
	Println(i ...interface{})

	// Printf is a convenience method to Printf to the defined output, fallback to Stderr if not set.
	Printf(format string, i ...interface{})

	// PrintErr is a convenience method to Print to the defined Err output, fallback to Stderr if not set.
	PrintErr(i ...interface{})

	// PrintErrln is a convenience method to Println to the defined Err output, fallback to Stderr if not set.
	PrintErrln(i ...interface{})

	// PrintErrf is a convenience method to Printf to the defined Err output, fallback to Stderr if not set.
	PrintErrf(format string, i ...interface{})
}
