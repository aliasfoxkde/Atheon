<img src="docs/logo.svg" alt="Leakr" width="72" /><br>

# Leakr

![Scanners](https://img.shields.io/badge/scanners-6-blue)
![Java](https://img.shields.io/badge/Java-17%2B-orange)
![Maven Central](https://img.shields.io/maven-central/v/io.github.horadomu/leakr)
![License](https://img.shields.io/badge/license-MIT-green)

**Secret detection for Java applications. A library you embed in your code, and a CLI you run anywhere.**

---

## What is Leakr?

Leakr detects leaked secrets — API keys, tokens, and credentials — in your code, config files, and environment variables. It does one thing and does it well.

- **As a library**: call `runner.scanString(...)`, `runner.scanFile(...)`, or `runner.scanDir(...)` directly from your Java application. No subprocess. No native binary. Just a dependency.
- **As a CLI tool**: run `leakr scan <path>` from any terminal on any platform with a JRE.

---

## Install

**Maven**

```xml
<dependency>
    <groupId>io.github.horadomu</groupId>
    <artifactId>leakr</artifactId>
    <version>1.0.0</version>
</dependency>
```

**Gradle**

```gradle
implementation 'io.github.horadomu:leakr:1.0.0'
```

---

## Library Usage

```java
import leakr.core.*;
import java.nio.file.*;
import java.util.List;

Runner runner = new Runner(new Registry());

List<Finding> findings = runner.scanString("AKIAIOSFODNN7EXAMPLE");
List<Finding> findings = runner.scanFile(Path.of("config.yaml"));
List<Finding> findings = runner.scanDir(Path.of("/path/to/project"));
List<Finding> findings = runner.scanEnv();
```

Each `Finding` exposes: `scanner`, `severity`, `file`, `line`, `description`, `match`. No configuration required.

---

## CLI Quickstart

```bash
git clone https://github.com/HoraDomu/Leakr.git
cd Leakr
mvn package -q
```

This produces two JARs in `target/`:
- `leakr-1.0.0-cli.jar` — self-contained JAR for CLI use
- `leakr-1.0.0.jar` — thin library JAR for embedding as a dependency

---

## Run the CLI

```bash
# Scan a directory
java -jar target/leakr-1.0.0-cli.jar scan /path/to/project

# Scan a single file
java -jar target/leakr-1.0.0-cli.jar scan config.env

# Scan environment variables
java -jar target/leakr-1.0.0-cli.jar scan --env

# Pipe content from stdin
git diff | java -jar target/leakr-1.0.0-cli.jar scan --stdin

# JSON output
java -jar target/leakr-1.0.0-cli.jar scan /path/to/project --json

# Exclude directories
java -jar target/leakr-1.0.0-cli.jar scan . --exclude target,dist,node_modules

# Filter by file extension
java -jar target/leakr-1.0.0-cli.jar scan . --ext .env,.yaml,.json,.tf

# List registered scanners
java -jar target/leakr-1.0.0-cli.jar list
```

Exit code `0` means clean. Exit code `1` means findings were detected.

---

## Sample Output

```
[CRITICAL] aws-access-key
  file:  src/config/dev.properties:14
  desc:  Detects AWS access key IDs (AKIA...)
  match: AKIA**********************MPLE

─────────────────────────────
found 1 potential secret(s)

files: 42  size: 318.7 KB  time: 84ms
```

---

## Adding a Scanner

Leakr auto-discovers every class in the `leakr.scanners` package that implements `Scanner`. No registration, no config file — drop a class in, rebuild, and it is live.

### 1. Create the scanner

```java
package leakr.scanners;

import leakr.core.*;
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

Place it in `src/leakr/scanners/MyServiceScanner.java`.

### 2. Rebuild and verify

```bash
mvn package -q
java -jar target/leakr-1.0.0-cli.jar list
```

### 3. Add a test case

Open `src/leakr/test/ScannerTest.java` and add one entry to `CASES`:

```java
new Case("myservice-api-key",
    "myservice_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",  // must match
    "myservice_tooshort")                            // must not match
```

---

## Testing

```bash
java -cp target/leakr-1.0.0-cli.jar leakr.test.ScannerTest
```

```
PASS     aws-access-key
PASS     github-pat
PASS     openai-api-key
PASS     stripe-secret-key
PASS     slack-bot-token
PASS     twilio-account-sid

6 passed, 0 failed
```

Exit code `0` on full pass, `1` on any failure.

---

> [!WARNING]
> Leakr is feature complete. Future releases will be security patches and new scanners only. The library API, CLI, output formats, and exit codes are stable and will not change.

---

## License

MIT — free to use, modify, and distribute. The copyright notice must be preserved in all copies. See [LICENSE](LICENSE).
