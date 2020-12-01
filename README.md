# etwstacks
Allows dumping stacks of running goroutines via ETW's capture state mechanism.

To use, simply include this package so that `init` is called:

```golang
import (
    _ "github.com/kevpar/etwstacks"
)
```

To capture stacks, there are many tools that can be used to start ETW traces. You will need one that also supports the ETW capture state mechanism. 

## WPR

WPR (Windows Performance Recorder) is one such tool that can be used. To use it to collect a trace, use etwstacks.wprp from this repo. Then, start and stop a trace as follows, this will capture the stacks due to the `CaptureStateOnSave` element in the profile:

```
> wpr -start etwstacks.wprp
> wpr -stop trace.etl
```

Then the trace can be opened using any ETW analysis tool, such as WPA (Windows Performance Analyzer).