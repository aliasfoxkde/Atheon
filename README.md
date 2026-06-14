<pre>
        /\     /\
       /  \___/  \
      / .-------. \
     / /         \ \
    / /   .---.   \ \
   | |   /     \   | |
   | |  |  [*]  |  | |
   | |   \     /   | |
    \ \   '---'   / /
     \ \         / /
      \ '-------' /
       \_________/
            |
           |||
</pre>

# Atheon

![Java](https://img.shields.io/badge/Java-17%2B-orange)
![Scanners](https://img.shields.io/badge/scanners-6-blue)
![License](https://img.shields.io/badge/license-MIT-green)

**A secret detection GUI for Java. Download one JAR. Run it anywhere a JRE exists.**

---

## What is Atheon?

Atheon scans your code, config files, and environment variables for leaked API keys, tokens, and credentials. It runs as a self-contained desktop GUI — no install, no configuration, no native binary. Double-click the JAR or run it from a terminal.

---

## Why it matters

Secrets end up in code all the time — an API key hardcoded for a quick test, a `.env` file committed by accident, a config template that shipped with real credentials still in it. Once a secret is in git history it is effectively public. Atheon catches them before that happens.

---

## Why a GUI?

Every other secret scanner is CLI-only. Atheon's terminal-style GUI runs on any platform with a JRE — Windows, macOS, Linux, ARM, Docker, CI — without downloading a platform-specific binary. Open it, scan, done.

---

## Install

Download `atheon.jar` from [GitHub Releases](https://github.com/HoraDomu/Atheon/releases/latest).

**Run:**
```bash
java -jar atheon.jar
```

That's it. The window opens.

---

## Commands

Type any of these at the `atheon>` prompt:

| Command | What it does |
|---|---|
| `scan` | Open a folder picker and scan the directory |
| `scan file` | Open a file picker and scan one file |
| `scan env` | Scan all environment variables |
| `list` | Show every registered scanner |
| `help` | Show all commands |
| `clear` | Clear the terminal |
| `//new` | Open a second session tab |
| `//exit` | Close Atheon |

---

## Adding a scanner

Drop a single `.java` file into `src/atheon/scanners/`. Implement `Scanner`, rebuild, and it appears automatically — no registration, no config.

```java
package atheon.scanners;

import atheon.core.*;
import java.util.*;
import java.util.regex.*;

public class MyServiceScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("myservice_[a-zA-Z0-9]{32}");

    public String name()        { return "myservice-api-key"; }
    public String description() { return "Detects MyService API keys"; }
    public Severity severity()  { return Severity.HIGH; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
```

---

## Testing a scanner

Open `src/atheon/test/ScannerTest.java` and add one entry to `CASES`:

```java
new Case("myservice-api-key",
    "myservice_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",  // must match
    "myservice_short")                               // must not match
```

Build and run:

```bash
mvn package -q
java -cp target/atheon.jar atheon.test.ScannerTest
```

---

## License

MIT License — Copyright © 2026 Dominick Yanez

Free to use, modify, and distribute. The copyright notice must be preserved in all copies.
