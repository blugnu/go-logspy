<div align="center" style="margin-bottom:20px">
  <!-- <img src=".assets/banner.png" alt="go-logspy" /> -->
  <div align="center">
    <a href="https://github.com/blugnu/go-logspy/actions/workflows/qa.yml"><img alt="build-status" src="https://github.com/blugnu/go-logspy/actions/workflows/qa.yml/badge.svg?branch=master&style=flat-square"/></a>
    <a href="https://goreportcard.com/report/github.com/blugnu/go-logspy" ><img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/go-logspy"/></a>
    <a><img alt="go version >= 1.14" src="https://img.shields.io/github/go-mod/go-version/blugnu/go-logspy?style=flat-square"/></a>
    <a href="https://github.com/blugnu/go-logspy/blob/master/LICENSE"><img alt="MIT License" src="https://img.shields.io/github/license/blugnu/go-logspy?color=%234275f5&style=flat-square"/></a>
    <a href="https://coveralls.io/github/blugnu/go-logspy?branch=master"><img alt="coverage" src="https://img.shields.io/coveralls/github/blugnu/go-logspy?style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/blugnu/go-logspy"><img alt="docs" src="https://pkg.go.dev/badge/github.com/blugnu/go-logspy"/></a>
  </div>
</div>

<br>

# go-logspy

A truly trivial package for assisting in incorporating log output verification in `go` unit tests.

LogSpy works with any log package that provides a mechanism for redirecting log output to an `io.Writer`, such as the built-in `log` package or `logrus`.

## How LogSpy Works
Log output is redirected into a `bytes.Buffer` (the "sink").

After code under test has been executed, the contents of the log are tested using either normal string testing techniques, or any of the various helper methods provided (intended to simplify common log tests).

That's it.

<br>
<hr>
<br>

## How to Use LogSpy
LogSpy is as trivial to use as it is to understand:

1. Configure and sink the log
2. `Reset()` LogSpy before each test
3. Execute code to be tested
4. Test log content with helpers (or your preferred string testing techniques)

<br>

### **1. Configure and Sink the Log**
In your tests, redirect log output to the `logspy.Sink()`.

When using a package such as `logrus`, this is typically combined with configuring your log package for test runs in the same way that it is configured for normal execution of the code under test.

If using `go test`, a good place for this could be a `TestMain`:

```golang
func TestMain(m *testing.M) {
    // Confgure the log formatter
    // (TIP: Use a common func to ensure identical formatting
    //        in both test and application logging)
    logrus.SetFormatter(&logrus.JSONFormatter{})

    // Redirect log output to the logspy sink 
	logrus.SetOutput(logspy.Sink())

    // Run the tests and set the exit value to the result
	os.Exit(m.Run())
}
```

**NOTE:** *Since `go test` runs tests at the package level, it is necessary to configure and sink the log in each package; a common func called from `TestMain()` in each package might be one way to achieve this.*

<br>

### **2. `Reset()` Logs Before Each Test (ARRANGE)**
With log formatting and sink in place, you then need to ensure you `Reset()` LogSpy at the beginning of each test:

```golang
func TestThatSomeFunctionEmitsExpectedLogs(t *testing.T) {
    // ARRANGE
    logspy.Reset()

    // ACT
    SomeFunction()

    // ASSERT
    ...
}
```

This is important as it ensures that log output captured in one test does not "pollute" the log in others.

<br>

### **3. Execute Code Under Test (ACT)**
This doesn't need any explanation, right?

<br>

### **4. Test Log Content (ASSERT)**
The raw log output captured by LogSpy is available using the `logspy.String()` function.  Using this, the log content can be tested using whatever techniques you prefer when testing string values.

However, LogSpy also provides a number of helper functions to simplify common tests of log content.

For example, to test that an expected number of log entries have been emitted, the `NumEntries()` function can be used, which returns the number of non-empty log entries (since the most recent `Reset()`).

An example showing this in use:

```golang
func TestThatNumEntriesReturnsTheNumberOfLogEntries(t *testing.T) {
	// ARRANGE
	logspy.Reset()

	// ACT
	log.Println("output 1")
	log.Println("output 2")
	log.Println("output 3")

	// ASSERT
	wanted := 3
	got := logspy.NumEntries()
	if wanted != got {
		t.Errorf("Wanted %d log entries, got %d", wanted, got)
	}
}
```
