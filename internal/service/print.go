package service

type Printer interface {
	// Print is a convenience method to Print to the defined output, fallback to Stderr if not set.
	Print(i ...any)

	// Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
	Println(i ...any)

	// Printf is a convenience method to Printf to the defined output, fallback to Stderr if not set.
	Printf(format string, i ...any)

	// PrintErr is a convenience method to Print to the defined Err output, fallback to Stderr if not set.
	PrintErr(i ...any)

	// PrintErrln is a convenience method to Println to the defined Err output, fallback to Stderr if not set.
	PrintErrln(i ...any)

	// PrintErrf is a convenience method to Printf to the defined Err output, fallback to Stderr if not set.
	PrintErrf(format string, i ...any)
}
